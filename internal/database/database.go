package database

import (
	"context"
	"fmt"
	"github.com/jmoiron/sqlx"
	log "github.com/sirupsen/logrus"
	"os"
	"profile/internal/profile"
	"time"
)

type Database struct {
	Client *sqlx.DB
}

const (
	maxRetries    = 5
	retryInterval = time.Second * 5
)

func NewDatabase() (*Database, error) {
	log.Info("Setting up new database connection")

	connectionString := fmt.Sprintf(
		"host=%s port=%s user=%s dbname=%s password=%s sslmode=%s",
		getOrDefault("DB_HOST", "localhost"),
		getOrDefault("DB_PORT", "5433"),
		getOrDefault("DB_USERNAME", "postgres"),
		getOrDefault("DB_TABLE", "postgres"),
		getOrDefault("DB_PASSWORD", "postgres"),
		getOrDefault("SSL_MODE", "disable"),
	)

	log.Infof("Connection string: %v", connectionString)
	var db *sqlx.DB
	var err error

	for i := 0; i < maxRetries; i++ {
		db, err = sqlx.Connect("postgres", connectionString)
		if err == nil {
			log.Info("Connected to database")
			return &Database{
				Client: db,
			}, nil
		}

		log.Errorf("Could not connect to database: %v", err)

		if i < maxRetries-1 {
			log.Infof("Retrying database connection in %s...", retryInterval)
			time.Sleep(retryInterval)
		}
	}

	return &Database{}, fmt.Errorf("failed to connect to database after %d retries", maxRetries)
}

func (d *Database) Ping() error {
	return d.Client.Ping()
}

func (d *Database) GetProfiles(ctx context.Context) ([]profile.Profile, error) {
	var profileRows []ProfileRow
	query := `SELECT s.id, s.name, a.street AS address_street, a.city AS address_city, 
					 a.state AS address_state, a.postal_code AS address_zip, 
					 a.country AS address_country, s.user_id 
			  FROM profiles s 
			  INNER JOIN addresses a ON s.id = a.profile_id`
	err := d.Client.SelectContext(ctx, &profileRows, query)
	if err != nil {
		return []profile.Profile{}, profile.ErrFetchingProfile
	}

	var profiles []profile.Profile
	for _, profileRow := range profileRows {
		s := convertProfileRowToProfile(profileRow)
		profiles = append(profiles, s)
	}

	return profiles, nil
}
func (d *Database) GetProfile(ctx context.Context, id string) (profile.Profile, error) {
	var profileRow ProfileRow
	query := `SELECT s.id, s.name, a.street AS address_street, a.city AS address_city, 
                     a.state AS address_state, a.postal_code AS address_zip, 
                     a.country AS address_country, s.user_id 
              FROM profiles s 
              INNER JOIN addresses a ON s.id = a.profile_id 
              WHERE s.id = $1`
	err := d.Client.GetContext(ctx, &profileRow, query, id)
	if err != nil {
		log.Errorf("Error getting profile: %v", err)
		return profile.Profile{}, profile.ErrFetchingProfile
	}

	s := convertProfileRowToProfile(profileRow)

	return s, nil
}

func (d *Database) PostProfile(ctx context.Context, p profile.Profile) (profile.Profile, error) {
	tx, err := d.Client.BeginTx(ctx, nil)
	if err != nil {
		return profile.Profile{}, profile.ErrFetchingProfile
	}
	query := "INSERT INTO profiles (name, user_id) VALUES ($1, $2) RETURNING id"
	err = tx.QueryRowContext(ctx, query, p.Name, p.UserId).Scan(&p.ID)
	if err != nil {
		tx.Rollback()
		return profile.Profile{}, profile.ErrFetchingProfile
	}
	query = "INSERT INTO addresses (profile_id, street, city, state, postal_code, country) VALUES ($1, $2, $3, $4, $5, $6)"
	_, err = tx.ExecContext(ctx, query, p.ID, p.Address.Street, p.Address.City, p.Address.State, p.Address.Zip, p.Address.Country)
	if err != nil {
		tx.Rollback()
		return profile.Profile{}, profile.ErrFetchingProfile
	}
	err = tx.Commit()
	return p, err
}

func (d *Database) UpdateProfile(ctx context.Context, p profile.Profile) (profile.Profile, error) {
	tx, err := d.Client.BeginTx(ctx, nil)
	query := "UPDATE profiles SET name = $1 WHERE id = $2"
	_, err = d.Client.ExecContext(ctx, query, p.Name, p.ID)
	if err != nil {
		tx.Rollback()
		return profile.Profile{}, profile.ErrUpdatingProfile
	}
	query = "UPDATE addresses SET street = $1, city = $2, state = $3, postal_code = $4, country = $5 WHERE profile_id = $6"
	_, err = d.Client.ExecContext(ctx, query, p.Address.Street, p.Address.City, p.Address.State, p.Address.Zip, p.Address.Country, p.ID)
	if err != nil {
		tx.Rollback()
		return profile.Profile{}, profile.ErrUpdatingProfile
	}
	err = tx.Commit()
	return p, err
}

func (d *Database) DeleteProfile(ctx context.Context, s string) error {
	tx, err := d.Client.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	query := "DELETE FROM addresses WHERE profile_id = $1"
	_, err = tx.ExecContext(ctx, query, s)
	if err != nil {
		tx.Rollback()
		return profile.ErrDeletingProfile
	}
	query = "DELETE FROM profiles WHERE id = $1"
	_, err = tx.ExecContext(ctx, query, s)
	if err != nil {
		tx.Rollback()
		return profile.ErrDeletingProfile
	}
	err = tx.Commit()
	return err
}
func getOrDefault(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}

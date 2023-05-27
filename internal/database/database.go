package database

import (
	"context"
	"fmt"
	"github.com/jmoiron/sqlx"
	log "github.com/sirupsen/logrus"
	"os"
	"profile/internal/seller"
	"time"
)

type Database struct {
	Client *sqlx.DB
}

const maxRetries = 5
const retryInterval = time.Second * 5

func NewDatabase() (*Database, error) {
	log.Info("Setting up new database connection")

	connectionString := fmt.Sprintf(
		"host=%s port=%s user=%s dbname=%s password=%s sslmode=%s",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_USERNAME"),
		os.Getenv("DB_TABLE"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("SSL_MODE"),
	)

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

func (d *Database) GetSeller(ctx context.Context, id string) (seller.Seller, error) {
	var sellerRow SellerRow
	query := `SELECT s.id, s.name, a.street AS address_street, a.city AS address_city, 
                     a.state AS address_state, a.postal_code AS address_zip, 
                     a.country AS address_country, s.user_id 
              FROM sellers s 
              INNER JOIN addresses a ON s.id = a.seller_id 
              WHERE s.id = $1`
	err := d.Client.GetContext(ctx, &sellerRow, query, id)
	s := convertSellerRowToSeller(sellerRow)
	if err != nil {
		return s, err
	}
	return s, nil
}

func (d *Database) PostSeller(ctx context.Context, seller seller.Seller) (seller.Seller, error) {
	tx, err := d.Client.BeginTx(ctx, nil)
	if err != nil {
		return seller, err
	}
	query := "INSERT INTO sellers (name, user_id) VALUES ($1, $2) RETURNING id"
	err = tx.QueryRowContext(ctx, query, seller.Name, seller.UserId).Scan(&seller.ID)
	if err != nil {
		tx.Rollback()
		return seller, err
	}
	query = "INSERT INTO addresses (seller_id, street, city, state, postal_code, country) VALUES ($1, $2, $3, $4, $5, $6)"
	_, err = tx.ExecContext(ctx, query, seller.ID, seller.Address.Street, seller.Address.City, seller.Address.State, seller.Address.Zip, seller.Address.Country)
	if err != nil {
		tx.Rollback()
		return seller, err
	}
	err = tx.Commit()
	return seller, err
}

func (d *Database) UpdateSeller(ctx context.Context, seller seller.Seller) (seller.Seller, error) {
	query := "UPDATE sellers SET name = $1 WHERE id = $2"
	_, err := d.Client.ExecContext(ctx, query, seller.Name, seller.ID)
	if err != nil {
		return seller, err
	}
	query = "UPDATE addresses SET street = $1, city = $2, state = $3, postal_code = $4, country = $5 WHERE seller_id = $6"
	_, err = d.Client.ExecContext(ctx, query, seller.Address.Street, seller.Address.City, seller.Address.State, seller.Address.Zip, seller.Address.Country, seller.ID)
	return seller, err
}

func (d *Database) DeleteSeller(ctx context.Context, s string) error {
	tx, err := d.Client.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	query := "DELETE FROM addresses WHERE seller_id = $1"
	_, err = tx.ExecContext(ctx, query, s)
	if err != nil {
		tx.Rollback()
		return err
	}
	query = "DELETE FROM sellers WHERE id = $1"
	_, err = tx.ExecContext(ctx, query, s)
	if err != nil {
		tx.Rollback()
		return err
	}
	err = tx.Commit()
	return err
}

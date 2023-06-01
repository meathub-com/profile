package profile

import (
	"context"
	"errors"
)

var (
	ErrFetchingProfile = errors.New("Error fetching profile")
	ErrPostingProfile  = errors.New("Error posting profile")
	ErrUpdatingProfile = errors.New("Error updating profile")
	ErrDeletingProfile = errors.New("Error deleting profile")
)

type Profile struct {
	ID      string
	Name    string
	Address Address
	UserId  string
}
type Address struct {
	Street  string
	City    string
	State   string
	Zip     string
	Country string
}

type Store interface {
	GetProfile(context.Context, string) (Profile, error)
	PostProfile(context.Context, Profile) (Profile, error)
	UpdateProfile(context.Context, Profile) (Profile, error)
	DeleteProfile(context.Context, string) error
	GetProfiles(context.Context) ([]Profile, error)
}
type Service struct {
	Store Store
}

func NewService(store Store) *Service {
	return &Service{
		Store: store,
	}
}
func (s *Service) GetProfile(ctx context.Context, id string) (Profile, error) {
	return s.Store.GetProfile(ctx, id)
}
func (s *Service) PostProfile(ctx context.Context, profile Profile) (Profile, error) {
	return s.Store.PostProfile(ctx, profile)
}
func (s *Service) UpdateProfile(ctx context.Context, profile Profile) (Profile, error) {
	return s.Store.UpdateProfile(ctx, profile)
}
func (s *Service) DeleteProfile(ctx context.Context, id string) error {
	return s.Store.DeleteProfile(ctx, id)
}
func (s *Service) GetProfiles(ctx context.Context) ([]Profile, error) {
	return s.Store.GetProfiles(ctx)
}

package seller

import "context"

type Seller struct {
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

type SellerStore interface {
	GetSeller(context.Context, string) (Seller, error)
	PostSeller(context.Context, Seller) (Seller, error)
	UpdateSeller(context.Context, Seller) (Seller, error)
	DeleteSeller(context.Context, string) error
}
type Service struct {
	Store SellerStore
}

func NewService(store SellerStore) *Service {
	return &Service{
		Store: store,
	}
}
func (s *Service) GetSeller(ctx context.Context, id string) (Seller, error) {
	return s.Store.GetSeller(ctx, id)
}
func (s *Service) PostSeller(ctx context.Context, seller Seller) (Seller, error) {
	return s.Store.PostSeller(ctx, seller)
}
func (s *Service) UpdateSeller(ctx context.Context, seller Seller) (Seller, error) {
	return s.Store.UpdateSeller(ctx, seller)
}
func (s *Service) DeleteSeller(ctx context.Context, id string) error {
	return s.Store.DeleteSeller(ctx, id)
}

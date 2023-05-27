package transport

import (
	"context"
	"profile/internal/seller"
)

type SellerService interface {
	GetSeller(ctx context.Context, id string) (seller.Seller, error)
}

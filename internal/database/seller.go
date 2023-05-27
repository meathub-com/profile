package database

import "profile/internal/seller"

type SellerRow struct {
	ID     string `db:"id"`
	Name   string `db:"name"`
	UserID string `db:"user_id"`
	AddressRow
}

type AddressRow struct {
	Street  string `db:"address_street"`
	City    string `db:"address_city"`
	State   string `db:"address_state"`
	Zip     string `db:"address_zip"`
	Country string `db:"address_country"`
}

func convertSellerRowToSeller(sellerRow SellerRow) seller.Seller {
	return seller.Seller{
		ID:      sellerRow.ID,
		Name:    sellerRow.Name,
		Address: convertAddressRowToAddress(sellerRow.AddressRow),
		UserId:  sellerRow.UserID,
	}
}
func convertAddressRowToAddress(addressRow AddressRow) seller.Address {
	return seller.Address{
		Street:  addressRow.Street,
		City:    addressRow.City,
		State:   addressRow.State,
		Zip:     addressRow.Zip,
		Country: addressRow.Country,
	}
}

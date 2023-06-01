package database

import "profile/internal/profile"

type ProfileRow struct {
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

func convertProfileRowToProfile(profileRow ProfileRow) profile.Profile {
	return profile.Profile{
		ID:      profileRow.ID,
		Name:    profileRow.Name,
		Address: convertAddressRowToAddress(profileRow.AddressRow),
		UserId:  profileRow.UserID,
	}
}
func convertAddressRowToAddress(addressRow AddressRow) profile.Address {
	return profile.Address{
		Street:  addressRow.Street,
		City:    addressRow.City,
		State:   addressRow.State,
		Zip:     addressRow.Zip,
		Country: addressRow.Country,
	}
}

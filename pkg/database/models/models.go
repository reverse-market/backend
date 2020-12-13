package models

import "errors"

var (
	ErrDuplicateEmail = errors.New("duplicate email")
	ErrNoRecord       = errors.New("no such record")
)

type User struct {
	ID               int    `json:"id"`
	Name             string `json:"name"`
	Email            string `json:"email"`
	Photo            string `json:"photo"`
	DefaultAddressID *int   `json:"default_address_id"`
}

type Address struct {
	ID     int
	UserID int
	Info   AddressInfo
}

type AddressInfo struct {
	Name       string `json:"name"`
	Region     string `json:"region"`
	City       string `json:"city"`
	Street     string `json:"street"`
	Number     string `json:"number"`
	Building   string `json:"building"`
	Appartment string `json:"appartment"`
	Zip        int    `json:"zip"`
	LastName   string `json:"last_name"`
	FirstName  string `json:"first_name"`
	FatherName string `json:"father_name"`
}

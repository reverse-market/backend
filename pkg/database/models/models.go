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

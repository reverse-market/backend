package models

import "errors"

var (
	ErrDuplicateEmail = errors.New("duplicate email")
	ErrNoRecord       = errors.New("no such record")
)

type User struct {
	ID               int     `json:"id"`
	Name             string  `json:"name"`
	Email            string  `json:"email"`
	Avatar           *string `json:"avatar"`
	DefaultAddressID *int    `json:"default_address_id"`
}

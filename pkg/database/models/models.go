package models

import (
	"errors"
	"github.com/reverse-market/backend/pkg/simpletime"
)

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

type Category struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Photo string `json:"photo"`
}

type Tag struct {
	ID         int    `json:"id"`
	CategoryID *int   `json:"-"`
	Name       string `json:"name"`
}

type TagFilters struct {
	CategoryID *int
	Search     string
}

type Request struct {
	ID          int                   `json:"id"`
	UserID      int                   `json:"-"`
	Username    string                `json:"username"`
	CategoryID  int                   `json:"category_id"`
	Name        string                `json:"name"`
	ItemName    string                `json:"item_name"`
	Description string                `json:"description"`
	Photos      []string              `json:"photos"`
	Price       int                   `json:"price"`
	Quantity    int                   `json:"quantity"`
	Date        simpletime.SimpleTime `json:"date"`
	Finished    bool                  `json:"-"`
	Tags        []*struct {
		ID   int    `json:"id"`
		Name string `json:"name"`
	} `json:"tags"`
}

type RequestFilters struct {
	Page          int
	Size          int
	CategoryID    *int
	Tags          []int
	PriceFrom     *int
	PriceTo       *int
	SortColumn    string
	SortDirection string
	Search        string
}

type Proposal struct {
	ID          int                   `json:"id"`
	UserID      int                   `json:"-"`
	Username    string                `json:"username"`
	RequestID   int                   `json:"request_id"`
	Name        string                `json:"name"`
	ItemName    string                `json:"item_name"`
	Description string                `json:"description"`
	Photos      []string              `json:"photos"`
	Price       int                   `json:"price"`
	Quantity    int                   `json:"quantity"`
	Date        simpletime.SimpleTime `json:"date"`
	BoughtById  *int                  `json:"-"`
}

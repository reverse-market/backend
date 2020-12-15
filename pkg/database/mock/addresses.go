package mock

import (
	"context"
	"github.com/reverse-market/backend/pkg/database/models"
)

type AddressRepository struct{}

var mockAddresses = []*models.Address{
	{
		ID:     1,
		UserID: 1,
		Info: models.AddressInfo{
			Name:       "Mock address 1",
			Region:     "Some region",
			City:       "Some city",
			Street:     "Mock street",
			Number:     "1",
			Building:   "1",
			Appartment: "1",
			Zip:        123,
			LastName:   "Ivanov",
			FirstName:  "Ivan",
			FatherName: "Ivanovich",
		},
	},
	{
		ID:     2,
		UserID: 2,
		Info: models.AddressInfo{
			Name:       "Mock address",
			Region:     "Some region",
			City:       "Some city",
			Street:     "Mock street",
			Number:     "1",
			Building:   "1",
			Appartment: "1",
			Zip:        123,
			LastName:   "Petrov",
			FirstName:  "Petr",
			FatherName: "Petrovich",
		},
	},
}

func (a AddressRepository) Add(ctx context.Context, address *models.Address) (int, error) {
	return 3, nil
}

func (a AddressRepository) GetByID(ctx context.Context, id int) (*models.Address, error) {
	switch id {
	case 1:
		return mockAddresses[0], nil
	case 2:
		return mockAddresses[1], nil
	default:
		return nil, models.ErrNoRecord
	}
}

func (a AddressRepository) GetByUserID(ctx context.Context, userID int) ([]*models.Address, error) {
	switch userID {
	case 1:
		return []*models.Address{mockAddresses[0]}, nil
	default:
		return []*models.Address{}, nil
	}
}

func (a AddressRepository) Update(ctx context.Context, address *models.Address) error {
	return nil
}

func (a AddressRepository) Delete(ctx context.Context, id int) error {
	return nil
}

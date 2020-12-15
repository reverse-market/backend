package mock

import (
	"context"
	"github.com/reverse-market/backend/pkg/database/models"
)

type UserRepository struct {
}

var (
	id1       = 1
	id2       = 2
	mockUsers = []*models.User{
		{
			ID:               1,
			Name:             "Ivanov Ivan",
			Email:            "ivanov@mail.ru",
			Photo:            "/image/smth1.jpg",
			DefaultAddressID: &id1,
		},
		{
			ID:               2,
			Name:             "Petrov Petr",
			Email:            "petrox@mail.ru",
			Photo:            "/image/smth.jpg",
			DefaultAddressID: &id2,
		},
	}
)

func (ur *UserRepository) Add(ctx context.Context, user *models.User) (int, error) {
	if user.Email == "ivanov@mail.ru" || user.Email == "petrox@mail.ru" {
		return 0, models.ErrDuplicateEmail
	} else {
		return 3, nil
	}
}

func (ur *UserRepository) GetByID(ctx context.Context, id int) (*models.User, error) {
	if id > 0 && id <= len(mockUsers) {
		return mockUsers[id-1], nil
	} else {
		return nil, models.ErrNoRecord
	}
}

func (ur *UserRepository) GetByEmail(ctx context.Context, email string) (*models.User, error) {
	switch email {
	case "ivanov@mail.ru":
		return mockUsers[0], nil
	case "petrox@mail.ru":
		return mockUsers[1], nil
	default:
		return nil, models.ErrNoRecord
	}
}

func (ur *UserRepository) GetAll(ctx context.Context) ([]*models.User, error) {
	return mockUsers, nil
}

func (ur *UserRepository) Update(ctx context.Context, user *models.User) error {
	return nil
}

func (ur *UserRepository) Delete(ctx context.Context, id int) error {
	return nil
}

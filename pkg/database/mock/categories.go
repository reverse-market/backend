package mock

import (
	"context"
	"github.com/reverse-market/backend/pkg/database/models"
)

type CategoriesRepository struct {
}

var mockCategory = &models.Category{
	ID:    1,
	Name:  "Mock category",
	Photo: "/image/smth.jpg",
}

func (c CategoriesRepository) GetByID(ctx context.Context, id int) (*models.Category, error) {
	switch id {
	case 1:
		return mockCategory, nil
	default:
		return nil, models.ErrNoRecord
	}
}

func (c CategoriesRepository) GetAll(ctx context.Context) ([]*models.Category, error) {
	return []*models.Category{mockCategory}, nil
}

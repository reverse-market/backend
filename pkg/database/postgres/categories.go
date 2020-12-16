package postgres

import (
	"context"
	"errors"
	"github.com/georgysavva/scany/pgxscan"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/reverse-market/backend/pkg/database/models"
)

type CategoriesRepository struct {
	DB *pgxpool.Pool
}

func (cr *CategoriesRepository) GetByID(ctx context.Context, id int) (*models.Category, error) {
	conn, err := cr.DB.Acquire(ctx)
	if err != nil {
		return nil, err
	}
	defer conn.Release()

	stmt := "SELECT id, name, photo FROM categories WHERE id=$1"

	category := &models.Category{}
	if err := pgxscan.Get(ctx, conn, category, stmt, id); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, models.ErrNoRecord
		}
		return nil, err
	}

	return category, nil
}

func (cr *CategoriesRepository) GetAll(ctx context.Context) ([]*models.Category, error) {
	conn, err := cr.DB.Acquire(ctx)
	if err != nil {
		return nil, err
	}
	defer conn.Release()

	stmt := "SELECT id, name, photo FROM categories ORDER BY id"

	categories := make([]*models.Category, 0)
	if err := pgxscan.Select(ctx, conn, &categories, stmt); err != nil {
		return nil, err
	}

	return categories, nil
}

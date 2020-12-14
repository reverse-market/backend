package postgres

import (
	"context"
	"errors"
	"github.com/georgysavva/scany/pgxscan"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/reverse-market/backend/pkg/database/models"
)

type TagsRepository struct {
	DB *pgxpool.Pool
}

func (tr *TagsRepository) GetByID(ctx context.Context, id int) (*models.Tag, error) {
	conn, err := tr.DB.Acquire(ctx)
	if err != nil {
		return nil, err
	}
	defer conn.Release()

	stmt := "SELECT id, category_id, name FROM tags WHERE id=$1"

	tag := &models.Tag{}
	if err := pgxscan.Get(ctx, conn, tag, stmt, id); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, models.ErrNoRecord
		}
		return nil, err
	}

	return tag, nil
}

func (tr *TagsRepository) GetAll(ctx context.Context, filters *models.TagFilters) ([]*models.Tag, error) {
	conn, err := tr.DB.Acquire(ctx)
	if err != nil {
		return nil, err
	}
	defer conn.Release()

	stmt := "SELECT id, category_id, name FROM tags WHERE ($1::int IS NULL OR category_id=$1 OR category_id IS NULL) " +
		"AND ($2 = '' OR LOWER(name) LIKE CONCAT('%', $2, '%'))"

	tags := make([]*models.Tag, 0)
	if err := pgxscan.Select(ctx, conn, &tags, stmt, filters.CategoryID, filters.Search); err != nil {
		return nil, err
	}

	return tags, nil
}

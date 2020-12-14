package postgres

import (
	"context"
	"errors"
	"github.com/georgysavva/scany/pgxscan"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/reverse-market/backend/pkg/database/models"
)

type RequestsRepository struct {
	DB *pgxpool.Pool
}

func (rr *RequestsRepository) Add(ctx context.Context, request *models.Request) (int, error) {
	conn, err := rr.DB.Acquire(ctx)
	if err != nil {
		return 0, err
	}
	defer conn.Release()

	tx, err := conn.Begin(ctx)
	if err != nil {
		return 0, err
	}
	defer tx.Rollback(ctx)

	stmt := "INSERT INTO requests (user_id, category_id, name, item_name, description, photos, price, quantity, date) " +
		"VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9) RETURNING id"

	var id int
	if err := tx.QueryRow(ctx, stmt,
		request.UserID,
		request.CategoryID,
		request.Name,
		request.ItemName,
		request.Description,
		request.Photos,
		request.Price,
		request.Quantity,
		request.Date,
	).Scan(&id); err != nil {
		return 0, err
	}

	stmt = "INSERT INTO requests_tags (request_id, tag_id) VALUES ($1, $2)"
	for _, tag := range request.Tags {
		if _, err := tx.Exec(ctx, stmt, id, tag.ID); err != nil {
			return 0, err
		}
	}

	if err := tx.Commit(ctx); err != nil {
		return 0, err
	}

	return id, nil
}

func (rr *RequestsRepository) GetByID(ctx context.Context, id int) (*models.Request, error) {
	conn, err := rr.DB.Acquire(ctx)
	if err != nil {
		return nil, err
	}
	defer conn.Release()

	stmt := "SELECT * FROM requests_view WHERE id=$1"

	request := &models.Request{}
	if err := pgxscan.Get(ctx, conn, request, stmt, id); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, models.ErrNoRecord
		}
		return nil, err
	}

	return request, nil
}

func (rr *RequestsRepository) GetByUserID(ctx context.Context, userID int, search string) ([]*models.Request, error) {
	conn, err := rr.DB.Acquire(ctx)
	if err != nil {
		return nil, err
	}
	defer conn.Release()

	stmt := "SELECT * FROM requests_view WHERE user_id=$1 AND " +
		"($2 = '' OR " +
		"LOWER(name) LIKE CONCAT('%', $2, '%') OR " +
		"LOWER(item_name) LIKE CONCAT('%', $2, '%') OR " +
		"LOWER(description) LIKE CONCAT('%', $2, '%'))"

	requests := make([]*models.Request, 0)
	if err := pgxscan.Select(ctx, conn, &requests, stmt, userID, search); err != nil {
		return nil, err
	}

	return requests, nil
}

func (rr *RequestsRepository) Update(ctx context.Context, request *models.Request) error {
	conn, err := rr.DB.Acquire(ctx)
	if err != nil {
		return err
	}
	defer conn.Release()

	tx, err := conn.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	stmt := "UPDATE requests SET user_id=$1, category_id=$2, name=$3, item_name=$4," +
		" description=$5, photos=$6, price=$7, quantity=$8, date=$9 WHERE id=$10"

	if _, err := tx.Exec(ctx, stmt,
		request.UserID,
		request.CategoryID,
		request.Name,
		request.ItemName,
		request.Description,
		request.Photos,
		request.Price,
		request.Quantity,
		request.Date,
		request.ID,
	); err != nil {
		return err
	}

	stmt = "DELETE FROM requests_tags WHERE request_id=$1"
	if _, err := tx.Exec(ctx, stmt, request.ID); err != nil {
		return err
	}

	stmt = "INSERT INTO requests_tags (request_id, tag_id) VALUES ($1, $2)"
	for _, tag := range request.Tags {
		if _, err := tx.Exec(ctx, stmt, request.ID, tag.ID); err != nil {
			return err
		}
	}

	return tx.Commit(ctx)
}

func (rr *RequestsRepository) Delete(ctx context.Context, id int) error {
	conn, err := rr.DB.Acquire(ctx)
	if err != nil {
		return err
	}
	defer conn.Release()

	stmt := "DELETE FROM requests WHERE id=$1"

	if _, err := conn.Exec(ctx, stmt, id); err != nil {
		return err
	}

	return nil
}

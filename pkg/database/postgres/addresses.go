package postgres

import (
	"context"
	"errors"
	"github.com/georgysavva/scany/pgxscan"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/reverse-market/backend/pkg/database/models"
)

type AddressRepository struct {
	DB *pgxpool.Pool
}

func (ar *AddressRepository) Add(ctx context.Context, address *models.Address) (int, error) {
	conn, err := ar.DB.Acquire(ctx)
	if err != nil {
		return 0, err
	}
	defer conn.Release()

	stmt := "INSERT INTO addresses (user_id, info) VALUES ($1, $2) RETURNING id"
	var id int
	if err := conn.QueryRow(ctx, stmt,
		address.UserID,
		address.Info,
	).Scan(&id); err != nil {
		return 0, err
	}

	return id, nil
}

func (ar *AddressRepository) GetById(ctx context.Context, id int) (*models.Address, error) {
	conn, err := ar.DB.Acquire(ctx)
	if err != nil {
		return nil, err
	}
	defer conn.Release()

	address := &models.Address{}
	stmt := "SELECT id, user_id, info FROM addresses WHERE id=$1"
	if err := pgxscan.Get(ctx, conn, address, stmt, id); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, models.ErrNoRecord
		}
		return nil, err
	}

	return address, nil
}

func (ar *AddressRepository) GetByUserID(ctx context.Context, userID int) ([]*models.Address, error) {
	conn, err := ar.DB.Acquire(ctx)
	if err != nil {
		return nil, err
	}
	defer conn.Release()

	addresses := make([]*models.Address, 0)
	stmt := "SELECT id, user_id, info FROM addresses WHERE user_id=$1"
	if err := pgxscan.Select(ctx, conn, &addresses, stmt, userID); err != nil {
		return nil, err
	}

	return addresses, nil
}

func (ar *AddressRepository) Update(ctx context.Context, address *models.Address) error {
	conn, err := ar.DB.Acquire(ctx)
	if err != nil {
		return err
	}
	defer conn.Release()

	stmt := "UPDATE addresses SET user_id=$1, info=$2 WHERE id=$3"
	if _, err := conn.Exec(ctx, stmt, address.UserID, address.Info, address.ID); err != nil {
		return err
	}

	return nil
}

func (ar *AddressRepository) Delete(ctx context.Context, id int) error {
	conn, err := ar.DB.Acquire(ctx)
	if err != nil {
		return err
	}
	defer conn.Release()

	stmt := "DELETE FROM addresses WHERE id=$1"
	if _, err := conn.Exec(ctx, stmt, id); err != nil {
		return err
	}

	return nil
}

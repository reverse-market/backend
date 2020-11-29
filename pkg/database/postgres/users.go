package postgres

import (
	"context"
	"errors"

	"github.com/georgysavva/scany/pgxscan"
	"github.com/jackc/pgconn"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/reverse-market/backend/pkg/database/models"
)

type UserRepository struct {
	DB *pgxpool.Pool
}

func (ur *UserRepository) Add(ctx context.Context, user *models.User) (int, error) {
	conn, err := ur.DB.Acquire(ctx)
	if err != nil {
		return 0, err
	}
	defer conn.Release()

	stmt := "INSERT INTO users (name, email, avatar, default_address_id) " +
		"VALUES ($1, $2, $3, null) RETURNING id"

	var id int
	if err := conn.QueryRow(ctx, stmt, user.Name, user.Email, user.Avatar).Scan(&id); err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr); pgErr.Code == pgerrcode.UniqueViolation {
			return 0, models.ErrDuplicateEmail
		}
		return 0, err
	}

	return id, nil
}

func (ur *UserRepository) GetById(ctx context.Context, id int) (*models.User, error) {
	conn, err := ur.DB.Acquire(ctx)
	if err != nil {
		return nil, err
	}
	defer conn.Release()

	user := &models.User{}
	stmt := "SELECT id, name, email, avatar, default_address_id FROM users WHERE id=$1"
	if err := pgxscan.Get(ctx, conn, user, stmt, id); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, models.ErrNoRecord
		}
		return nil, err
	}

	return user, nil
}

func (ur *UserRepository) GetByEmail(ctx context.Context, email string) (*models.User, error) {
	conn, err := ur.DB.Acquire(ctx)
	if err != nil {
		return nil, err
	}
	defer conn.Release()

	user := &models.User{}
	stmt := "SELECT id, name, email, avatar, default_address_id FROM users WHERE email=$1"
	if err := pgxscan.Get(ctx, conn, user, stmt, email); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, models.ErrNoRecord
		}
		return nil, err
	}

	return user, nil
}

func (ur *UserRepository) GetAll(ctx context.Context) ([]*models.User, error) {
	conn, err := ur.DB.Acquire(ctx)
	if err != nil {
		return nil, err
	}
	defer conn.Release()

	users := make([]*models.User, 0)
	stmt := "SELECT id, name, email, avatar, default_address_id FROM users WHERE id=$1"
	if err := pgxscan.Select(ctx, conn, &users, stmt); err != nil {
		return nil, err
	}

	return users, nil
}

func (ur *UserRepository) Update(ctx context.Context, user *models.User) error {
	conn, err := ur.DB.Acquire(ctx)
	if err != nil {
		return err
	}
	defer conn.Release()

	stmt := "UPDATE users SET name=$1, email=$2, avatar=$3, default_address_id=$4 WHERE id=$5"

	if _, err := conn.Exec(ctx, stmt, user.Name, user.Email, user.Avatar, user.DefaultAddressID, user.ID); err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr); pgErr.Code == pgerrcode.UniqueViolation {
			return models.ErrDuplicateEmail
		}
		return err
	}

	return nil
}

func (ur *UserRepository) Delete(ctx context.Context, id int) error {
	conn, err := ur.DB.Acquire(ctx)
	if err != nil {
		return err
	}
	defer conn.Release()

	stmt := "DELETE FROM users WHERE id = $1"
	if _, err := conn.Exec(ctx, stmt, id); err != nil {
		return err
	}

	return nil
}

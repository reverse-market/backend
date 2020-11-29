package postgres

import (
	"context"

	"github.com/jackc/pgx/v4/pgxpool"
)

type Config struct {
	URL string
}

func NewPsqlPool(c *Config) (*pgxpool.Pool, func(), error) {
	pool, err := pgxpool.Connect(context.Background(), c.URL)
	if err != nil {
		return nil, nil, err
	}

	return pool, pool.Close, nil
}

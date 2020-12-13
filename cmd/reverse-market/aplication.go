package main

import (
	"context"
	"log"

	"github.com/reverse-market/backend/pkg/database/models"
)

type Application struct {
	config    *config
	loggers   *loggers
	tokens    tokensManager
	users     usersRepository
	addresses addressesRepository
}

type loggers struct {
	info  *log.Logger
	error *log.Logger
}

type tokensManager interface {
	CreateToken(int) (string, error)
	GetIdFromToken(string) (int, error)
}

type usersRepository interface {
	Add(context.Context, *models.User) (int, error)
	GetById(context.Context, int) (*models.User, error)
	GetByEmail(context.Context, string) (*models.User, error)
	GetAll(context.Context) ([]*models.User, error)
	Update(context.Context, *models.User) error
	Delete(context.Context, int) error
}

type addressesRepository interface {
	Add(context.Context, *models.Address) (int, error)
	GetById(context.Context, int) (*models.Address, error)
	GetByUserID(context.Context, int) ([]*models.Address, error)
	Update(context.Context, *models.Address) error
	Delete(ctx context.Context, id int) error
}

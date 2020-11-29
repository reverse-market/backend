package main

import (
	"context"
	"log"

	"github.com/reverse-market/backend/pkg/database/models"
)

type Application struct {
	config  *config
	loggers *loggers
	tokens  tokensManager
	users   usersRepository
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

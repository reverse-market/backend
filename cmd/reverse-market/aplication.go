package main

import (
	"context"
	"log"

	"github.com/reverse-market/backend/pkg/database/models"
)

type Application struct {
	config     *config
	loggers    *loggers
	tokens     tokensManager
	users      usersRepository
	addresses  addressesRepository
	categories categoriesRepository
	tags       tagsRepository
	requests   requestsRepository
	proposals  proposalsRepository
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
	GetByID(context.Context, int) (*models.User, error)
	GetByEmail(context.Context, string) (*models.User, error)
	GetAll(context.Context) ([]*models.User, error)
	Update(context.Context, *models.User) error
	Delete(context.Context, int) error
}

type addressesRepository interface {
	Add(context.Context, *models.Address) (int, error)
	GetByID(context.Context, int) (*models.Address, error)
	GetByUserID(context.Context, int) ([]*models.Address, error)
	Update(context.Context, *models.Address) error
	Delete(ctx context.Context, id int) error
}

type categoriesRepository interface {
	GetByID(context.Context, int) (*models.Category, error)
	GetAll(context.Context) ([]*models.Category, error)
}

type tagsRepository interface {
	GetByID(context.Context, int) (*models.Tag, error)
	GetAll(context.Context, *models.TagFilters) ([]*models.Tag, error)
}

type requestsRepository interface {
	Add(context.Context, *models.Request) (int, error)
	GetByID(context.Context, int) (*models.Request, error)
	GetByUserID(context.Context, int, string) ([]*models.Request, error)
	Update(context.Context, *models.Request) error
	Delete(context.Context, int) error
}

type proposalsRepository interface {
	Add(context.Context, *models.Proposal) (int, error)
	GetByID(context.Context, int) (*models.Proposal, error)
	GetByUserIDSold(context.Context, int) ([]*models.Proposal, error)
	GetByUserIDBought(context.Context, int) ([]*models.Proposal, error)
	Update(context.Context, *models.Proposal) error
	Delete(context.Context, int) error
	Buy(context.Context, int, int) error
	GetBestForRequest(context.Context, int) (int, error)
}

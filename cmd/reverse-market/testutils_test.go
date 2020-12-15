package main

import (
	"github.com/reverse-market/backend/pkg/database/mock"
	"github.com/reverse-market/backend/pkg/jwt"
	"io/ioutil"
	"log"
)

func newTestApplication() *Application {
	return &Application{
		loggers: &loggers{
			info:  log.New(ioutil.Discard, "", 0),
			error: log.New(ioutil.Discard, "", 0),
		},
		tokens:     &jwt.MockManager{},
		users:      &mock.UserRepository{},
		addresses:  &mock.AddressRepository{},
		categories: &mock.CategoriesRepository{},
		tags:       &mock.TagsRepository{},
		requests:   &mock.RequestsRepository{},
		proposals:  &mock.ProposalsRepository{},
	}
}

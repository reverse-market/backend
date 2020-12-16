//+build wireinject

package main

import (
	"github.com/google/wire"
	"github.com/reverse-market/backend/pkg/database/postgres"
	"github.com/reverse-market/backend/pkg/idtoken"
	"github.com/reverse-market/backend/pkg/jwt"
)

func initApp() (*Application, func(), error) {
	wire.Build(
		getConfig,
		NewLoggers,
		NewJwtManager,
		wire.Bind(new(tokensManager), new(*jwt.Manager)),
		NewPostgresConfig,
		postgres.NewPsqlPool,
		wire.Struct(new(idtoken.Parser)),
		wire.Bind(new(tokenParser), new(*idtoken.Parser)),
		wire.Struct(new(postgres.UserRepository), "*"),
		wire.Bind(new(usersRepository), new(*postgres.UserRepository)),
		wire.Struct(new(postgres.AddressRepository), "*"),
		wire.Bind(new(addressesRepository), new(*postgres.AddressRepository)),
		wire.Struct(new(postgres.CategoriesRepository), "*"),
		wire.Bind(new(categoriesRepository), new(*postgres.CategoriesRepository)),
		wire.Struct(new(postgres.TagsRepository), "*"),
		wire.Bind(new(tagsRepository), new(*postgres.TagsRepository)),
		wire.Struct(new(postgres.RequestsRepository), "*"),
		wire.Bind(new(requestsRepository), new(*postgres.RequestsRepository)),
		wire.Struct(new(postgres.ProposalsRepository), "*"),
		wire.Bind(new(proposalsRepository), new(*postgres.ProposalsRepository)),
		wire.Struct(new(Application), "*"),
	)

	return nil, nil, nil
}

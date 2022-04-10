package main

import (
	"crypto/rand"
	"io"
	"log"
	"net/http"

	"github.com/reverse-market/backend/pkg/database/postgres"
	"github.com/reverse-market/backend/pkg/database/redis"
	"github.com/reverse-market/backend/pkg/jwt"
	"github.com/xlab/closer"
)

func main() {
	app, cleanup, err := initApp()
	if err != nil {
		log.Fatal("error initiating application", err)
	}
	closer.Bind(func() {
		log.Print("Stopping server")
		cleanup()
	})

	server := &http.Server{
		Addr:     ":" + app.config.Port,
		Handler:  app.route(true),
		ErrorLog: app.loggers.Error(),
	}

	app.loggers.Info().Printf("Started server on port: %v", app.config.Port)
	app.loggers.Error().Fatalf("server error: %e", server.ListenAndServe())
}

func NewJwtManager(c *config) *jwt.Manager {
	return &jwt.Manager{
		Secret:     c.Secret,
		Expiration: c.TokenTTL,
	}
}

func NewPostgresConfig(c *config) *postgres.Config {
	return &postgres.Config{URL: c.DBurl}
}

func NewRedisConfig(c *config) *redis.Config {
	return &redis.Config{
		RedisUrl:             c.RedisURL,
		SessionTTl:           c.SessionTTl,
		SessionCleanupPeriod: c.SessionCleanupPeriod,
		SessionWindowPeriod:  c.SessionWindowPeriod,
	}
}

func NewRandSource(_ *config) io.Reader {
	return rand.Reader
}

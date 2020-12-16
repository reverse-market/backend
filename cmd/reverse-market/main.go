package main

import (
	"log"
	"net/http"
	"os"

	"github.com/reverse-market/backend/pkg/database/postgres"
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
		ErrorLog: app.loggers.error,
	}

	app.loggers.info.Printf("Started server on port: %v", app.config.Port)
	app.loggers.error.Fatalf("server error: %e", server.ListenAndServe())
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

func NewLoggers() *loggers {
	return &loggers{
		info:  log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime),
		error: log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile),
	}
}

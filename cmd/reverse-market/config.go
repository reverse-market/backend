package main

import (
	"time"

	"github.com/caarlos0/env"
)

type config struct {
	Production           bool          `env:"PRODUCTION" envDefault:"false"`
	Port                 string        `env:"PORT" envDefault:"8080"`
	DBurl                string        `env:"DB_URL" envDefault:"postgresql://postgres:marketpassword@db:5432/reverse-market"`
	RedisURL             string        `env:"REDIS_URL" envDefault:"redis:6379"`
	TokenTTL             time.Duration `env:"TOKEN_TTL" envDefault:"20m"`
	SessionTTl           time.Duration `env:"SESSION_TTL" envDefault:"168h"`
	SessionCleanupPeriod time.Duration `env:"SESSION_CLEANUP_PERIOD" envDefault:"60s"`
	SessionWindowPeriod  time.Duration `env:"SESSION_WINDOW_PERIOD" envDefault:"60s"`
	SessionTokenLength   int           `env:"SESSION_TOKEN_LENGTH" envDefault:"32"`
	Secret               string        `env:"SECRET" envDefault:"CGI3fpuoncYv1QfaNsrqiZeU5FOYpDPi"`
}

func getConfig() (*config, error) {
	c := &config{}
	if err := env.Parse(c); err != nil {
		return nil, err
	}
	return c, nil
}

package factory

import (
	"time"

	"github.com/minguu42/harmattan/internal/auth"
	"github.com/minguu42/harmattan/internal/database"
)

type Config struct {
	API  API
	Auth auth.Config
	DB   database.Config
}

type API struct {
	Host         string        `env:"API_HOST" default:"0.0.0.0"`
	Port         int           `env:"API_PORT" default:"8080"`
	ReadTimeout  time.Duration `env:"API_READ_TIMEOUT" default:"2s"`
	WriteTimeout time.Duration `env:"API_WRITE_TIMEOUT" default:"2s"`
	StopTimeout  time.Duration `env:"API_STOP_TIMEOUT" default:"25s"`
}

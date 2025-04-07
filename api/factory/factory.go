package factory

import (
	"fmt"

	"github.com/minguu42/harmattan/internal/auth"
	"github.com/minguu42/harmattan/internal/database"
)

type Factory struct {
	Auth *auth.Authenticator
	DB   *database.Client
}

func New(conf Config) (*Factory, error) {
	authn, err := auth.NewAuthenticator(conf.Auth)
	if err != nil {
		return nil, fmt.Errorf("failed to create authenticator: %w", err)
	}

	db, err := database.NewClient(conf.DB)
	if err != nil {
		return nil, fmt.Errorf("failed to create database client: %w", err)
	}

	return &Factory{
		Auth: authn,
		DB:   db,
	}, nil
}

func (f *Factory) Close() {
	_ = f.DB.Close()
}

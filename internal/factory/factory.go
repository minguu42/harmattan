package factory

import (
	"context"
	"fmt"

	"github.com/minguu42/harmattan/internal/alog"
	"github.com/minguu42/harmattan/internal/auth"
	"github.com/minguu42/harmattan/internal/database"
)

type Factory struct {
	Auth *auth.Authenticator
	DB   *database.Client
}

func New(ctx context.Context, conf Config, logger *alog.Logger) (*Factory, error) {
	authn, err := auth.NewAuthenticator(conf.Auth)
	if err != nil {
		return nil, fmt.Errorf("failed to create authenticator: %w", err)
	}

	db, err := database.NewClient(ctx, conf.DB, logger)
	if err != nil {
		return nil, fmt.Errorf("failed to create database client: %w", err)
	}

	return &Factory{
		Auth: authn,
		DB:   db,
	}, nil
}

func (f *Factory) Close() error {
	if err := f.DB.Close(); err != nil {
		return fmt.Errorf("failed to close database: %w", err)
	}
	return nil
}

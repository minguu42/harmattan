package factory

import (
	"context"

	"github.com/minguu42/harmattan/internal/alog"
	"github.com/minguu42/harmattan/internal/auth"
	"github.com/minguu42/harmattan/internal/database"
	"github.com/minguu42/harmattan/internal/lib/errors"
)

type Factory struct {
	Auth *auth.Authenticator
	DB   *database.Client
}

func New(ctx context.Context, conf Config, logger *alog.Logger) (*Factory, error) {
	authn, err := auth.NewAuthenticator(conf.Auth)
	if err != nil {
		return nil, errors.Wrap(err)
	}

	db, err := database.NewClient(ctx, conf.DB, logger)
	if err != nil {
		return nil, errors.Wrap(err)
	}

	return &Factory{
		Auth: authn,
		DB:   db,
	}, nil
}

func (f *Factory) Close() error {
	if err := f.DB.Close(); err != nil {
		return errors.Wrap(err)
	}
	return nil
}

package factory

import (
	"context"

	"github.com/minguu42/harmattan/internal/auth"
	"github.com/minguu42/harmattan/internal/database"
	"github.com/minguu42/harmattan/internal/lib/errtrace"
)

type Factory struct {
	Auth *auth.Authenticator
	DB   *database.Client
}

func New(ctx context.Context, conf Config) (*Factory, error) {
	authn, err := auth.NewAuthenticator(conf.Auth)
	if err != nil {
		return nil, errtrace.Wrap(err)
	}

	db, err := database.NewClient(ctx, conf.DB)
	if err != nil {
		return nil, errtrace.Wrap(err)
	}

	return &Factory{
		Auth: authn,
		DB:   db,
	}, nil
}

func (f *Factory) Close() error {
	return errtrace.Wrap(f.DB.Close())
}

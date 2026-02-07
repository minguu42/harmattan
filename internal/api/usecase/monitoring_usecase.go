package usecase

import (
	"context"

	"github.com/minguu42/harmattan/internal/database"
	"github.com/minguu42/harmattan/internal/lib/errtrace"
)

type Monitoring struct {
	Revision string
	DB       *database.Client
}

type CheckHealthOutput struct {
	Revision string
}

func (uc *Monitoring) CheckHealth(ctx context.Context) (*CheckHealthOutput, error) {
	if err := uc.DB.Ping(ctx); err != nil {
		return nil, errtrace.Wrap(err)
	}
	return &CheckHealthOutput{Revision: uc.Revision}, nil
}

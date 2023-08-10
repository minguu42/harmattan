package database

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/minguu42/opepe/gen/sqlc"
	"github.com/minguu42/opepe/pkg/domain/model"
	"github.com/minguu42/opepe/pkg/domain/repository"
)

func (db *DB) GetUserByAPIKey(ctx context.Context, apiKey string) (*model.User, error) {
	u, err := sqlc.New(db.sqlDB).GetUserByAPIKey(ctx, apiKey)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, repository.ErrModelNotFound
		}
		return nil, fmt.Errorf("q.GetUserByAPIKey failed: %w", err)
	}
	return &model.User{
		ID:        u.ID,
		Name:      u.Name,
		CreatedAt: u.CreatedAt,
		UpdatedAt: u.UpdatedAt,
	}, nil
}

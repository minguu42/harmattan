package database

import (
	"context"
	"fmt"

	"github.com/minguu42/opepe/gen/sqlc"
	"github.com/minguu42/opepe/pkg/entity"
)

func (db *DB) GetUserByAPIKey(ctx context.Context, apiKey string) (*entity.User, error) {
	u, err := sqlc.New(db._db).GetUserByAPIKey(ctx, apiKey)
	if err != nil {
		return nil, fmt.Errorf("q.GetUserByAPIKey failed: %w", err)
	}
	return &entity.User{
		ID:        u.ID,
		Name:      u.Name,
		CreatedAt: u.CreatedAt,
		UpdatedAt: u.UpdatedAt,
	}, nil
}

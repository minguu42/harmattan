package database

import (
	"context"
	"fmt"

	"github.com/minguu42/mtasks/pkg/entity"

	"github.com/minguu42/mtasks/pkg/logging"
)

func (db *DB) GetUserByAPIKey(ctx context.Context, apiKey string) (*entity.User, error) {
	q := `SELECT id, name, created_at, updated_at FROM users WHERE api_key = ?`
	logging.Debugf(q)

	var u entity.User
	if err := db.QueryRowContext(ctx, q, apiKey).Scan(&u.ID, &u.Name, &u.CreatedAt, &u.UpdatedAt); err != nil {
		return nil, fmt.Errorf("db.QueryRowContext failed: %w", err)
	}
	return &u, nil
}

package database

import (
	"context"
	"fmt"

	"github.com/minguu42/mtasks/app"
	"github.com/minguu42/mtasks/app/logging"
)

func (db *DB) GetUserByAPIKey(ctx context.Context, apiKey string) (*app.User, error) {
	q := `SELECT id, name, created_at, updated_at FROM users WHERE api_key = ?`
	logging.Debugf(q)

	var u app.User
	if err := db.QueryRowContext(ctx, q, apiKey).Scan(&u.ID, &u.Name, &u.CreatedAt, &u.UpdatedAt); err != nil {
		return nil, fmt.Errorf("db.QueryRowContext failed: %w", err)
	}
	return &u, nil
}

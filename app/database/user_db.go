package database

import (
	"context"
	"fmt"

	"github.com/minguu42/mtasks/app"
	"github.com/minguu42/mtasks/app/logging"
)

func (db *DB) GetUserByToken(ctx context.Context, token string) (*app.User, error) {
	q := `SELECT ID, Name, created_at, updated_at FROM users WHERE token = ?`
	logging.Debugf(q)

	var u app.User
	if err := db.QueryRowContext(ctx, q, token).Scan(&u.ID, &u.Name, &u.CreatedAt, &u.UpdatedAt); err != nil {
		return nil, fmt.Errorf("db.QueryRowContext failed: %w", err)
	}
	return &u, nil
}

package app

import (
	"context"
	"fmt"

	"github.com/minguu42/mtasks/pkg/logging"
)

func (db *database) getUserByToken(ctx context.Context, token string) (*user, error) {
	q := `SELECT id, name, created_at, updated_at FROM users WHERE token = ?`
	logging.Debugf(q)

	var u user
	if err := db.QueryRowContext(ctx, q, token).Scan(&u.id, &u.name, &u.createdAt, &u.updatedAt); err != nil {
		return nil, fmt.Errorf("db.QueryRowContext failed: %w", err)
	}
	return &u, nil
}

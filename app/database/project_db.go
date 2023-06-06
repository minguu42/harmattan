package database

import (
	"context"
	"fmt"
	"time"

	"github.com/minguu42/mtasks/app"
	"github.com/minguu42/mtasks/app/logging"
)

func (db *DB) CreateProject(ctx context.Context, userID int64, name string) (*app.Project, error) {
	q := `INSERT INTO projects (user_id, name, created_at, updated_at) VALUES (?, ?, ?, ?)`
	logging.Debugf(q)

	createdAt := time.Now()
	result, err := db.ExecContext(ctx, q, userID, name, createdAt, createdAt)
	if err != nil {
		return nil, fmt.Errorf("db.ExecContext failed: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return nil, fmt.Errorf("result.LastInsertId failed: %w", err)
	}

	return &app.Project{
		ID:        id,
		UserID:    userID,
		Name:      name,
		CreatedAt: createdAt,
		UpdatedAt: createdAt,
	}, nil
}

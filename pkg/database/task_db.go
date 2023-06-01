package database

import (
	"context"
	"fmt"
	"time"

	"github.com/minguu42/mtasks/pkg/app"
	"github.com/minguu42/mtasks/pkg/logging"
)

func (db *DB) CreateTask(ctx context.Context, userID int64, title string) (*app.Task, error) {
	q := `INSERT INTO tasks (user_id, title, created_at, updated_at) VALUES (?, ?, ?, ?)`
	logging.Debugf(q)

	createdAt := time.Now()
	result, err := db.ExecContext(ctx, q, userID, title, createdAt, createdAt)
	if err != nil {
		return nil, fmt.Errorf("db.ExecContext failed: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return nil, fmt.Errorf("result.LastInsertId failed: %w", err)
	}

	return &app.Task{
		ID:        id,
		UserID:    userID,
		Title:     title,
		CreatedAt: createdAt,
		UpdatedAt: createdAt,
	}, nil
}

func (db *DB) GetTasksByUserID(ctx context.Context, userID int64) ([]*app.Task, error) {
	q := `SELECT ID, Title, completed_at, created_at, updated_at FROM tasks WHERE user_id = ?`
	logging.Debugf(q)

	rows, err := db.QueryContext(ctx, q, userID)
	if err != nil {
		return nil, fmt.Errorf("db.QueryContext failed: %w", err)
	}
	defer rows.Close()

	ts := make([]*app.Task, 0, 10)
	for rows.Next() {
		var t app.Task
		if err := rows.Scan(&t.ID, &t.Title, &t.CompletedAt, &t.CreatedAt, &t.UpdatedAt); err != nil {
			return nil, fmt.Errorf("rows.Scan failed: %w", err)
		}
		ts = append(ts, &t)
	}
	return ts, nil
}

func (db *DB) GetTaskByID(ctx context.Context, id int64) (*app.Task, error) {
	q := `SELECT user_id, Title,completed_at, created_at, updated_at FROM tasks WHERE ID = ?`
	logging.Debugf(q)

	t := app.Task{ID: id}
	if err := db.QueryRowContext(ctx, q, id).Scan(&t.UserID, &t.Title, &t.CompletedAt, &t.CreatedAt, &t.UpdatedAt); err != nil {
		return nil, fmt.Errorf("db.QueryRowContext failed: %w", err)
	}
	return &t, nil
}

func (db *DB) UpdateTask(ctx context.Context, id int64, completedAt *time.Time) error {
	q := `UPDATE tasks SET completed_at = ? WHERE ID = ?`
	logging.Debugf(q)

	if _, err := db.ExecContext(ctx, q, completedAt, id); err != nil {
		return fmt.Errorf("db.ExecContext failed: %w", err)
	}
	return nil
}

func (db *DB) DeleteTask(ctx context.Context, id int64) error {
	q := `DELETE FROM tasks WHERE ID = ?`
	logging.Debugf(q)

	if _, err := db.ExecContext(ctx, q, id); err != nil {
		return fmt.Errorf("db.ExecContext failed: %w", err)
	}
	return nil
}

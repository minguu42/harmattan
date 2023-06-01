package app

import (
	"context"
	"fmt"
	"time"

	"github.com/minguu42/mtasks/pkg/logging"
)

func (db *database) createTask(ctx context.Context, userID int64, title string) (*task, error) {
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

	return &task{
		id:        id,
		userID:    userID,
		title:     title,
		createdAt: createdAt,
		updatedAt: createdAt,
	}, nil
}

func (db *database) getTasksByUserID(ctx context.Context, userID int64) ([]*task, error) {
	q := `SELECT id, title, completed_at, created_at, updated_at FROM tasks WHERE user_id = ?`
	logging.Debugf(q)

	rows, err := db.QueryContext(ctx, q, userID)
	if err != nil {
		return nil, fmt.Errorf("db.QueryContext failed: %w", err)
	}
	defer rows.Close()

	ts := make([]*task, 0, 10)
	for rows.Next() {
		var t task
		if err := rows.Scan(&t.id, &t.title, &t.completedAt, &t.createdAt, &t.updatedAt); err != nil {
			return nil, fmt.Errorf("rows.Scan failed: %w", err)
		}
		ts = append(ts, &t)
	}
	return ts, nil
}

func (db *database) getTaskByID(ctx context.Context, id int64) (*task, error) {
	q := `SELECT user_id, title,completed_at, created_at, updated_at FROM tasks WHERE id = ?`
	logging.Debugf(q)

	t := task{id: id}
	if err := db.QueryRowContext(ctx, q, id).Scan(&t.userID, &t.title, &t.completedAt, &t.createdAt, &t.updatedAt); err != nil {
		return nil, fmt.Errorf("db.QueryRowContext failed: %w", err)
	}
	return &t, nil
}

func (db *database) updateTask(ctx context.Context, id int64, completedAt *time.Time) error {
	q := `UPDATE tasks SET completed_at = ? WHERE id = ?`
	logging.Debugf(q)

	if _, err := db.ExecContext(ctx, q, completedAt, id); err != nil {
		return fmt.Errorf("db.ExecContext failed: %w", err)
	}
	return nil
}

func (db *database) deleteTask(ctx context.Context, id int64) error {
	q := `DELETE FROM tasks WHERE id = ?`
	logging.Debugf(q)

	if _, err := db.ExecContext(ctx, q, id); err != nil {
		return fmt.Errorf("db.ExecContext failed: %w", err)
	}
	return nil
}

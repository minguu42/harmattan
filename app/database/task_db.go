package database

import (
	"context"
	"fmt"
	"time"

	"github.com/minguu42/mtasks/app"
	"github.com/minguu42/mtasks/app/logging"
)

func (db *DB) CreateTask(ctx context.Context, projectID int64, title string) (*app.Task, error) {
	q := `INSERT INTO tasks (project_id, title, created_at, updated_at) VALUES (?, ?, ?, ?)`
	logging.Debugf(q)

	now := time.Now()
	result, err := db.ExecContext(ctx, q, projectID, title, now, now)
	if err != nil {
		return nil, fmt.Errorf("db.ExecContext failed: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return nil, fmt.Errorf("result.LastInsertId failed: %w", err)
	}

	return &app.Task{
		ID:        id,
		ProjectID: projectID,
		Title:     title,
		CreatedAt: now,
		UpdatedAt: now,
	}, nil
}

func (db *DB) GetTasksByProjectID(ctx context.Context, projectID int64, sort string, limit, offset int) ([]*app.Task, error) {
	q := `SELECT id, title, completed_at, created_at, updated_at FROM tasks WHERE project_id = ? ORDER BY ? LIMIT ? OFFSET ?`
	logging.Debugf(q)

	rows, err := db.QueryContext(ctx, q, projectID, generateOrderByClause(sort), limit, offset)
	if err != nil {
		return nil, fmt.Errorf("db.QueryContext failed: %w", err)
	}
	defer rows.Close()

	ts := make([]*app.Task, 0, 20)
	for rows.Next() {
		t := app.Task{ProjectID: projectID}
		if err := rows.Scan(&t.ID, &t.Title, &t.CompletedAt, &t.CreatedAt, &t.UpdatedAt); err != nil {
			return nil, fmt.Errorf("rows.Scan failed: %w", err)
		}
		ts = append(ts, &t)
	}
	return ts, nil
}

func (db *DB) GetTaskByID(ctx context.Context, id int64) (*app.Task, error) {
	q := `SELECT project_id, title, completed_at, created_at, updated_at FROM tasks WHERE id = ?`
	logging.Debugf(q)

	t := app.Task{ID: id}
	if err := db.QueryRowContext(ctx, q, id).Scan(&t.ProjectID, &t.Title, &t.CompletedAt, &t.CreatedAt, &t.UpdatedAt); err != nil {
		return nil, fmt.Errorf("db.QueryRowContext failed: %w", err)
	}
	return &t, nil
}

func (db *DB) UpdateTask(ctx context.Context, id int64, completedAt *time.Time, updatedAt time.Time) error {
	q := `UPDATE tasks SET completed_at = ?, updated_at = ? WHERE id = ?`
	logging.Debugf(q)

	if _, err := db.ExecContext(ctx, q, completedAt, updatedAt, id); err != nil {
		return fmt.Errorf("db.ExecContext failed: %w", err)
	}
	return nil
}

func (db *DB) DeleteTask(ctx context.Context, id int64) error {
	q := `DELETE FROM tasks WHERE id = ?`
	logging.Debugf(q)

	if _, err := db.ExecContext(ctx, q, id); err != nil {
		return fmt.Errorf("db.ExecContext failed: %w", err)
	}
	return nil
}

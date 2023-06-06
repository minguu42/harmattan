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

func (db *DB) GetProjectsByUserID(ctx context.Context, userID int64) ([]*app.Project, error) {
	q := `SELECT id, name, created_at, updated_at FROM projects WHERE user_id = ?`
	logging.Debugf(q)

	rows, err := db.QueryContext(ctx, q, userID)
	if err != nil {
		return nil, fmt.Errorf("db.QueryContext failed: %w", err)
	}
	defer rows.Close()

	ps := make([]*app.Project, 0, 20)
	for rows.Next() {
		p := app.Project{UserID: userID}
		if err := rows.Scan(&p.ID, &p.Name, &p.CreatedAt, &p.UpdatedAt); err != nil {
			return nil, fmt.Errorf("rows.Scan failed: %w", err)
		}
		ps = append(ps, &p)
	}
	return ps, nil
}

func (db *DB) GetProjectByID(ctx context.Context, id int64) (*app.Project, error) {
	q := `SELECT user_id, name, created_at, updated_at FROM projects WHERE id = ?`
	logging.Debugf(q)

	p := app.Project{ID: id}
	if err := db.QueryRowContext(ctx, q, id).Scan(&p.UserID, &p.Name, &p.CreatedAt, &p.UpdatedAt); err != nil {
		return nil, fmt.Errorf("row.Scan failed: %w", err)
	}
	return &p, nil
}

func (db *DB) UpdateProject(ctx context.Context, id int64, name string, updatedAt time.Time) error {
	q := `UPDATE projects SET name = ?, updated_at = ? WHERE id = ?`
	logging.Debugf(q)

	if _, err := db.ExecContext(ctx, q, name, updatedAt, id); err != nil {
		return fmt.Errorf("db.ExecContext failed: %w", err)
	}
	return nil
}

func (db *DB) DeleteProject(ctx context.Context, id int64) error {
	q := `DELETE FROM projects WHERE id = ?`
	logging.Debugf(q)

	if _, err := db.ExecContext(ctx, q, id); err != nil {
		return fmt.Errorf("db.ExecContext failed: %w", err)
	}
	return nil
}

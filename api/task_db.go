package api

import (
	"fmt"
	"time"

	"github.com/minguu42/mtasks/pkg/logging"
)

func createTask(userID uint64, title string) (*task, error) {
	createdAt := time.Now()
	q := `INSERT INTO tasks (user_id, title, created_at, updated_at) VALUES (?, ?, ?, ?)`
	logging.Debugf(q)
	result, err := db.Exec(q, userID, title, createdAt, createdAt)
	if err != nil {
		return nil, fmt.Errorf("db.Exec failed: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return nil, fmt.Errorf("result.LastInsertId failed: %w", err)
	}

	return &task{
		id:          uint64(id),
		userID:      userID,
		title:       title,
		completedAt: nil,
		createdAt:   createdAt,
		updatedAt:   createdAt,
	}, nil
}

func getTasksByUserID(userID uint64) ([]*task, error) {
	q := `SELECT id, title, completed_at, created_at, updated_at FROM tasks WHERE user_id = ?`
	logging.Debugf(q)

	rows, err := db.Query(q, userID)
	if err != nil {
		return nil, fmt.Errorf("db.Query failed: %w", err)
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

func getTaskByID(id uint64) (*task, error) {
	q := `SELECT user_id, title,completed_at, created_at, updated_at FROM tasks WHERE id = ?`
	logging.Debugf(q)

	t := task{id: id}
	if err := db.QueryRow(q, id).Scan(&t.userID, &t.title, &t.completedAt, &t.createdAt, &t.updatedAt); err != nil {
		return nil, fmt.Errorf("row.Scan failed: %w", err)
	}

	return &t, nil
}

func updateTask(id uint64, completedAt *time.Time) error {
	q := `UPDATE tasks SET completed_at = ? WHERE id = ?`
	logging.Debugf(q)

	if _, err := db.Exec(q, completedAt, id); err != nil {
		return fmt.Errorf("db.Exec failed: %w", err)
	}
	return nil
}

func destroyTask(id uint64) error {
	q := `DELETE FROM tasks WHERE id = ?`
	logging.Debugf(q)

	if _, err := db.Exec(q, id); err != nil {
		return fmt.Errorf("db.Exec failed: %w", err)
	}
	return nil
}

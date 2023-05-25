package api

import (
	"fmt"
	"time"
)

func createTask(userID uint64, title string) (*task, error) {
	createdAt := time.Now()
	q := `INSERT INTO tasks (user_id, title, created_at, updated_at) VALUES (?, ?, ?, ?)`
	Debugf(q)
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

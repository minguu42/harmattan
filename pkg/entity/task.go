package entity

import "time"

type Task struct {
	ID          int64
	ProjectID   int64
	Title       string
	CompletedAt *time.Time
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

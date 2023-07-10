package entity

import "time"

type Task struct {
	ID          string
	ProjectID   string
	Title       string
	CompletedAt *time.Time
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

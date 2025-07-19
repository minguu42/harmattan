package domain

import "time"

type TaskID string

type Task struct {
	ID          TaskID
	UserID      UserID
	ProjectID   ProjectID
	Name        string
	Content     string
	Priority    int
	DueOn       *time.Time
	CompletedAt *time.Time
	CreatedAt   time.Time
	UpdatedAt   time.Time
	Steps       Steps
	Tags        Tags
}

type Tasks []Task

package domain

import "time"

type StepID string

type Step struct {
	ID          StepID
	UserID      UserID
	TaskID      TaskID
	Name        string
	CompletedAt *time.Time
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type Steps []Step

package domain

import "time"

// MaxStepsPerTask は1タスクに作成できるステップ数の上限
const MaxStepsPerTask = 20

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

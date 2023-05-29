package app

import "time"

type task struct {
	id          int64
	userID      int64
	title       string
	completedAt *time.Time
	createdAt   time.Time
	updatedAt   time.Time
}

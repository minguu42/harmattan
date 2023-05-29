package app

import "time"

type task struct {
	id          uint64
	userID      uint64
	title       string
	completedAt *time.Time
	createdAt   time.Time
	updatedAt   time.Time
}

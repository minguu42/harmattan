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

type taskResponse struct {
	ID          uint64     `json:"id"`
	Title       string     `json:"title"`
	CompletedAt *time.Time `json:"completedAt"`
	CreatedAt   time.Time  `json:"createdAt"`
	UpdatedAt   time.Time  `json:"updatedAt"`
}

type tasksResponse struct {
	Tasks []*taskResponse `json:"tasks"`
}

type postTasksRequest struct {
	Title string `json:"title"`
}

type patchTaskRequest struct {
	IsCompleted bool `json:"isCompleted"`
}

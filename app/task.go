package app

import (
	"time"

	"github.com/minguu42/mtasks/app/ogen"
)

type Task struct {
	ID          int64
	UserID      int64
	Title       string
	CompletedAt *time.Time
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

// newTaskResponse はデータベースモデルの Task からレスポンスモデルの ogen.Task を生成する
func newTaskResponse(t *Task) ogen.Task {
	completedAt := ogen.OptDateTime{}
	if t.CompletedAt != nil {
		completedAt = ogen.NewOptDateTime(*t.CompletedAt)
	}
	return ogen.Task{
		ID:          t.ID,
		Title:       t.Title,
		CompletedAt: completedAt,
		CreatedAt:   t.CreatedAt,
		UpdatedAt:   t.UpdatedAt,
	}
}

// newTasksResponse はデータベースモデルの Task のスライスからレスポンスモデルの ogen.Task のスライスを生成する
func newTasksResponse(ts []*Task) []ogen.Task {
	tasks := make([]ogen.Task, 0, len(ts))
	for _, t := range ts {
		tasks = append(tasks, newTaskResponse(t))
	}
	return tasks
}

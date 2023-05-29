package app

import (
	"time"

	"github.com/minguu42/mtasks/pkg/ogen"
)

type task struct {
	id          int64
	userID      int64
	title       string
	completedAt *time.Time
	createdAt   time.Time
	updatedAt   time.Time
}

// newTaskResponse はデータベースモデルの task からレスポンスモデルの ogen.Task を生成する
func newTaskResponse(t *task) ogen.Task {
	completedAt := ogen.OptDateTime{}
	if t.completedAt != nil {
		completedAt = ogen.NewOptDateTime(*t.completedAt)
	}
	return ogen.Task{
		ID:          t.id,
		Title:       t.title,
		CompletedAt: completedAt,
		CreatedAt:   t.createdAt,
		UpdatedAt:   t.updatedAt,
	}
}

// newTasksResponse はデータベースモデルの task のスライスからレスポンスモデルの ogen.Task のスライスを生成する
func newTasksResponse(ts []*task) []ogen.Task {
	tasks := make([]ogen.Task, 0, len(ts))
	for _, t := range ts {
		tasks = append(tasks, newTaskResponse(t))
	}
	return tasks
}

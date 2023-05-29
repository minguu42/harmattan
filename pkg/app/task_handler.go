package app

import (
	"context"
	"fmt"
	"time"

	"github.com/minguu42/mtasks/pkg/logging"
	"github.com/minguu42/mtasks/pkg/ogen"
)

var (
	token = "rAM9Fm9huuWEKLdCwHBcju9Ty_-TL2tDsAicmMrXmUnaCGp3RtywzYpMDPdEtYtR"
)

// PostTasks は POST /tasks に対応するハンドラ
func (h *Handler) PostTasks(_ context.Context, req *ogen.PostTasksReq) (ogen.PostTasksRes, error) {
	u, err := getUserByToken(token)
	if err != nil {
		logging.Errorf("getUserByToken failed: %v", err)
		return &ogen.PostTasksUnauthorized{}, nil
	}

	t, err := createTask(u.id, req.Title)
	if err != nil {
		logging.Errorf("createTask failed: %v", err)
		return &ogen.PostTasksBadRequest{}, nil
	}

	completedAt := ogen.OptDateTime{}
	if t.completedAt != nil {
		completedAt = ogen.NewOptDateTime(*t.completedAt)
	}
	resp := ogen.Task{
		ID:          t.id,
		Title:       t.title,
		CompletedAt: completedAt,
		CreatedAt:   t.createdAt,
		UpdatedAt:   t.updatedAt,
	}
	return &ogen.TaskHeaders{
		Location: fmt.Sprintf("http://localhost:8080/tasks/%d", t.id),
		Response: resp,
	}, nil
}

// GetTasks は GET /tasks に対応するハンドラ
func (h *Handler) GetTasks(_ context.Context) (ogen.GetTasksRes, error) {
	u, err := getUserByToken(token)
	if err != nil {
		logging.Errorf("getUserByToken failed: %v", err)
		return &ogen.GetTasksUnauthorized{}, nil
	}

	ts, err := getTasksByUserID(u.id)
	if err != nil {
		logging.Errorf("getTasksByUserID failed: %v", err)
		return &ogen.GetTasksBadRequest{}, nil
	}

	tasks := make([]ogen.Task, 0, len(ts))
	for _, t := range ts {
		completedAt := ogen.OptDateTime{}
		task := ogen.Task{
			ID:          t.id,
			Title:       t.title,
			CompletedAt: completedAt,
			CreatedAt:   t.createdAt,
			UpdatedAt:   t.updatedAt,
		}
		tasks = append(tasks, task)
	}
	return &ogen.Tasks{Tasks: tasks}, nil
}

// PatchTask は PATCH /tasks/{taskID} に対応するハンドラ
func (h *Handler) PatchTask(_ context.Context, req *ogen.PatchTaskReq, params ogen.PatchTaskParams) (ogen.PatchTaskRes, error) {
	u, err := getUserByToken(token)
	if err != nil {
		logging.Errorf("getUserByToken failed: %v", err)
		return &ogen.PatchTaskUnauthorized{}, nil
	}

	t, err := getTaskByID(params.TaskID)
	if err != nil {
		logging.Errorf("getTaskByID failed: %v", err)
		return &ogen.PatchTaskBadRequest{}, nil
	}

	if t.userID != u.id {
		logging.Errorf("t.userID != user.id")
		return &ogen.PatchTaskNotFound{}, nil
	}

	if req.IsCompleted.IsSet() {
		if req.IsCompleted.Value {
			now := time.Now()
			if err := updateTask(params.TaskID, &now); err != nil {
				logging.Errorf("updateTask failed: %v", err)
				// TODO: InternalServerError の方が望ましい
				return &ogen.PatchTaskBadRequest{}, nil
			}
		} else {
			if err := updateTask(params.TaskID, nil); err != nil {
				logging.Errorf("updateTask failed: %v", err)
				// TODO: InternalServerError の方が望ましい
				return &ogen.PatchTaskBadRequest{}, nil
			}
		}
	}

	completedAt := ogen.OptDateTime{}
	if t.completedAt != nil {
		completedAt = ogen.NewOptDateTime(*t.completedAt)
	}
	return &ogen.Task{
		ID:          t.id,
		Title:       t.title,
		CompletedAt: completedAt,
		CreatedAt:   t.createdAt,
		UpdatedAt:   t.updatedAt,
	}, nil
}

// DeleteTask は DELETE /tasks/{taskID} に対応するハンドラ
func (h *Handler) DeleteTask(_ context.Context, params ogen.DeleteTaskParams) (ogen.DeleteTaskRes, error) {
	u, err := getUserByToken(token)
	if err != nil {
		logging.Errorf("getUserByToken failed: %v", err)
		return &ogen.DeleteTaskUnauthorized{}, nil
	}

	t, err := getTaskByID(params.TaskID)
	if err != nil {
		logging.Errorf("getTaskByID failed: %v", err)
		return &ogen.DeleteTaskBadRequest{}, nil
	}

	if t.userID != u.id {
		logging.Errorf("t.userID != u.id")
		return &ogen.DeleteTaskNotFound{}, nil
	}

	if err := destroyTask(t.id); err != nil {
		logging.Errorf("destroyTask failed: %v", err)
		return &ogen.DeleteTaskBadRequest{}, nil
	}

	return &ogen.DeleteTaskNoContent{}, nil
}

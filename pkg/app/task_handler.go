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
func (h *Handler) PostTasks(ctx context.Context, req *ogen.PostTasksReq) (ogen.PostTasksRes, error) {
	u, err := h.Repository.getUserByToken(ctx, token)
	if err != nil {
		logging.Errorf("getUserByToken failed: %v", err)
		return &ogen.PostTasksUnauthorized{}, nil
	}

	t, err := h.Repository.createTask(ctx, u.id, req.Title)
	if err != nil {
		logging.Errorf("createTask failed: %v", err)
		return &ogen.PostTasksBadRequest{}, nil
	}

	return &ogen.TaskHeaders{
		Location: fmt.Sprintf("http://localhost:8080/tasks/%d", t.id),
		Response: newTaskResponse(t),
	}, nil
}

// GetTasks は GET /tasks に対応するハンドラ
func (h *Handler) GetTasks(ctx context.Context) (ogen.GetTasksRes, error) {
	u, err := h.Repository.getUserByToken(ctx, token)
	if err != nil {
		logging.Errorf("getUserByToken failed: %v", err)
		return &ogen.GetTasksUnauthorized{}, nil
	}

	ts, err := h.Repository.getTasksByUserID(ctx, u.id)
	if err != nil {
		logging.Errorf("getTasksByUserID failed: %v", err)
		return &ogen.GetTasksBadRequest{}, nil
	}

	return &ogen.Tasks{Tasks: newTasksResponse(ts)}, nil
}

// PatchTask は PATCH /tasks/{taskID} に対応するハンドラ
func (h *Handler) PatchTask(ctx context.Context, req *ogen.PatchTaskReq, params ogen.PatchTaskParams) (ogen.PatchTaskRes, error) {
	u, err := h.Repository.getUserByToken(ctx, token)
	if err != nil {
		logging.Errorf("getUserByToken failed: %v", err)
		return &ogen.PatchTaskUnauthorized{}, nil
	}

	t, err := h.Repository.getTaskByID(ctx, params.TaskID)
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
			if err := h.Repository.updateTask(ctx, params.TaskID, &now); err != nil {
				logging.Errorf("updateTask failed: %v", err)
				// TODO: InternalServerError の方が望ましい
				return &ogen.PatchTaskBadRequest{}, nil
			}
		} else {
			if err := h.Repository.updateTask(ctx, params.TaskID, nil); err != nil {
				logging.Errorf("updateTask failed: %v", err)
				// TODO: InternalServerError の方が望ましい
				return &ogen.PatchTaskBadRequest{}, nil
			}
		}
	}

	resp := newTaskResponse(t)
	return &resp, nil
}

// DeleteTask は DELETE /tasks/{taskID} に対応するハンドラ
func (h *Handler) DeleteTask(ctx context.Context, params ogen.DeleteTaskParams) (ogen.DeleteTaskRes, error) {
	u, err := h.Repository.getUserByToken(ctx, token)
	if err != nil {
		logging.Errorf("getUserByToken failed: %v", err)
		return &ogen.DeleteTaskUnauthorized{}, nil
	}

	t, err := h.Repository.getTaskByID(ctx, params.TaskID)
	if err != nil {
		logging.Errorf("getTaskByID failed: %v", err)
		return &ogen.DeleteTaskBadRequest{}, nil
	}

	if t.userID != u.id {
		logging.Errorf("t.userID != u.id")
		return &ogen.DeleteTaskNotFound{}, nil
	}

	if err := h.Repository.deleteTask(ctx, t.id); err != nil {
		logging.Errorf("destroyTask failed: %v", err)
		return &ogen.DeleteTaskBadRequest{}, nil
	}

	return &ogen.DeleteTaskNoContent{}, nil
}

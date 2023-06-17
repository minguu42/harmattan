package handler

import (
	"context"
	"time"

	"github.com/minguu42/mtasks/pkg/entity"

	"github.com/minguu42/mtasks/pkg/logging"
	"github.com/minguu42/mtasks/pkg/ogen"
)

// CreateTask は POST /projects/{projectID}/tasks に対応するハンドラ
func (h *Handler) CreateTask(ctx context.Context, req *ogen.CreateTaskReq, params ogen.CreateTaskParams) (*ogen.Task, error) {
	u, ok := ctx.Value(userKey{}).(*entity.User)
	if !ok {
		logging.Errorf("ctx.Value(userKey{}).(*User) failed")
		return nil, errUnauthorized
	}

	p, err := h.Repository.GetProjectByID(ctx, params.ProjectID)
	if err != nil {
		logging.Errorf("repository.GetProjectByID failed: %v", err)
		return nil, errInternalServerError
	}
	if u.ID != p.UserID {
		logging.Errorf("u.ID != p.UserID")
		return nil, errTaskNotFound
	}

	t, err := h.Repository.CreateTask(ctx, params.ProjectID, req.Title)
	if err != nil {
		logging.Errorf("repository.Create failed: %v", err)
		return nil, errInternalServerError
	}

	return newTaskResponse(t), nil
}

// ListTasks は GET /projects/{projectID}/tasks に対応するハンドラ
func (h *Handler) ListTasks(ctx context.Context, params ogen.ListTasksParams) (*ogen.Tasks, error) {
	u, ok := ctx.Value(userKey{}).(*entity.User)
	if !ok {
		logging.Errorf("ctx.Value(userKey{}).(*User) failed")
		return nil, errUnauthorized
	}

	p, err := h.Repository.GetProjectByID(ctx, params.ProjectID)
	if err != nil {
		logging.Errorf("repository.GetProjectByID failed: %v", err)
		return nil, errInternalServerError
	}
	if u.ID != p.UserID {
		logging.Errorf("u.ID != p.UserID")
		return nil, errTaskNotFound
	}

	ts, err := h.Repository.GetTasksByProjectID(ctx, p.ID, string(params.Sort.Or(ogen.ListTasksSortMinusCreatedAt)), params.Limit.Or(10)+1, params.Offset.Or(0))
	if err != nil {
		logging.Errorf("repository.GetTasksByProjectID failed: %v", err)
		return nil, errInternalServerError
	}

	hasNext := false
	if len(ts) == params.Limit.Or(10)+1 {
		hasNext = true
		ts = ts[:params.Limit.Or(10)]
	}

	return &ogen.Tasks{
		Tasks:   newTasksResponse(ts),
		HasNext: hasNext,
	}, nil
}

// UpdateTask は PATCH /projects/{projectID}/tasks/{taskID} に対応するハンドラ
func (h *Handler) UpdateTask(ctx context.Context, req *ogen.UpdateTaskReq, params ogen.UpdateTaskParams) (*ogen.Task, error) {
	u, ok := ctx.Value(userKey{}).(*entity.User)
	if !ok {
		logging.Errorf("ctx.Value(userKey{}).(*User) failed")
		return nil, errUnauthorized
	}

	p, err := h.Repository.GetProjectByID(ctx, params.ProjectID)
	if err != nil {
		logging.Errorf("repository.GetProjectByID failed: %v", err)
		return nil, errInternalServerError
	}
	if u.ID != p.UserID {
		logging.Errorf("u.ID != p.UserID")
		return nil, errProjectNotFound
	}
	t, err := h.Repository.GetTaskByID(ctx, params.TaskID)
	if err != nil {
		logging.Errorf("repository.GetTaskByID failed: %v", err)
		return nil, errInternalServerError
	}
	if p.ID != t.ProjectID {
		logging.Errorf("p.ID != t.ProjectID")
		return nil, errTaskNotFound
	}

	if !req.IsCompleted.IsSet() {
		logging.Errorf("value contains nothing")
		return nil, errBadRequest
	}
	now := time.Now()
	if req.IsCompleted.Value {
		t.CompletedAt = &now
	} else {
		t.CompletedAt = nil
	}
	t.UpdatedAt = now
	if err := h.Repository.UpdateTask(ctx, params.TaskID, t.CompletedAt, t.UpdatedAt); err != nil {
		logging.Errorf("repository.UpdateTask failed: %v", err)
		return nil, errInternalServerError
	}

	return newTaskResponse(t), nil
}

// DeleteTask は DELETE /projects/{projectID}/tasks/{taskID} に対応するハンドラ
func (h *Handler) DeleteTask(ctx context.Context, params ogen.DeleteTaskParams) error {
	u, ok := ctx.Value(userKey{}).(*entity.User)
	if !ok {
		logging.Errorf("ctx.Value(userKey{}).(*User) failed")
		return errUnauthorized
	}

	p, err := h.Repository.GetProjectByID(ctx, params.ProjectID)
	if err != nil {
		logging.Errorf("repository.GetProjectByID failed: %v", err)
		return errInternalServerError
	}
	if u.ID != p.UserID {
		logging.Errorf("u.ID != p.UserID")
		return errProjectNotFound
	}
	t, err := h.Repository.GetTaskByID(ctx, params.TaskID)
	if err != nil {
		logging.Errorf("repository.GetTaskByID failed: %v", err)
		return errInternalServerError
	}
	if p.ID != t.ProjectID {
		logging.Errorf("p.ID != t.ProjectID")
		return errTaskNotFound
	}

	if err := h.Repository.DeleteTask(ctx, t.ID); err != nil {
		logging.Errorf("repository.DeleteTask failed: %v", err)
		return errInternalServerError
	}

	return nil
}

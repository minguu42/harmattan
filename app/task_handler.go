package app

import (
	"context"
	"time"

	"github.com/minguu42/mtasks/app/logging"
	"github.com/minguu42/mtasks/app/ogen"
)

// CreateTask は POST /projects/{projectID}/tasks に対応するハンドラ
func (h *handler) CreateTask(ctx context.Context, req *ogen.CreateTaskReq, params ogen.CreateTaskParams) (ogen.CreateTaskRes, error) {
	u, ok := ctx.Value(userKey{}).(*User)
	if !ok {
		logging.Errorf("ctx.Value(userKey{}).(*User) failed")
		return &ogen.CreateTaskInternalServerError{
			Message: messageInternalServerError,
			Debug:   "ctx.Value(userKey{}).(*User) failed",
		}, nil
	}

	p, err := h.repository.GetProjectByID(ctx, params.ProjectID)
	if err != nil {
		logging.Errorf("repository.GetProjectByID failed: %v", err)
		return &ogen.CreateTaskInternalServerError{
			Message: messageInternalServerError,
			Debug:   err.Error(),
		}, nil
	}
	if u.ID != p.UserID {
		logging.Errorf("u.ID != p.UserID")
		return &ogen.CreateTaskNotFound{
			Message: messageNotFound,
			Debug:   err.Error(),
		}, nil
	}

	t, err := h.repository.CreateTask(ctx, params.ProjectID, req.Title)
	if err != nil {
		logging.Errorf("repository.Create failed: %v", err)
		return &ogen.CreateTaskInternalServerError{
			Message: messageInternalServerError,
			Debug:   err.Error(),
		}, nil
	}

	return newTaskResponse(t), nil
}

// ListTasks は GET /projects/{projectID}/tasks に対応するハンドラ
func (h *handler) ListTasks(ctx context.Context, params ogen.ListTasksParams) (ogen.ListTasksRes, error) {
	u, ok := ctx.Value(userKey{}).(*User)
	if !ok {
		logging.Errorf("ctx.Value(userKey{}).(*User) failed")
		return &ogen.ListTasksInternalServerError{
			Message: messageInternalServerError,
			Debug:   "ctx.Value(userKey{}).(*User) failed",
		}, nil
	}

	p, err := h.repository.GetProjectByID(ctx, params.ProjectID)
	if err != nil {
		logging.Errorf("repository.GetProjectByID failed: %v", err)
		return &ogen.ListTasksInternalServerError{
			Message: messageInternalServerError,
			Debug:   err.Error(),
		}, nil
	}
	if u.ID != p.UserID {
		logging.Errorf("u.ID != p.UserID")
		return &ogen.ListTasksNotFound{
			Message: messageNotFound,
			Debug:   err.Error(),
		}, nil
	}

	ts, err := h.repository.GetTasksByProjectID(ctx, p.ID, string(params.Sort.Or(ogen.ListTasksSortMinusCreatedAt)), params.Limit.Or(10)+1, params.Offset.Or(0))
	if err != nil {
		logging.Errorf("repository.GetTasksByProjectID failed: %v", err)
		return &ogen.ListTasksInternalServerError{
			Message: messageInternalServerError,
			Debug:   err.Error(),
		}, nil
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
func (h *handler) UpdateTask(ctx context.Context, req *ogen.UpdateTaskReq, params ogen.UpdateTaskParams) (ogen.UpdateTaskRes, error) {
	u, ok := ctx.Value(userKey{}).(*User)
	if !ok {
		logging.Errorf("ctx.Value(userKey{}).(*User) failed")
		return &ogen.UpdateTaskInternalServerError{
			Message: messageInternalServerError,
			Debug:   "ctx.Value(userKey{}).(*User) failed",
		}, nil
	}

	p, err := h.repository.GetProjectByID(ctx, params.ProjectID)
	if err != nil {
		logging.Errorf("repository.GetProjectByID failed: %v", err)
		return &ogen.UpdateTaskInternalServerError{
			Message: messageInternalServerError,
			Debug:   err.Error(),
		}, nil
	}
	if u.ID != p.UserID {
		logging.Errorf("u.ID != p.UserID")
		return &ogen.UpdateTaskNotFound{
			Message: messageNotFound,
			Debug:   err.Error(),
		}, nil
	}
	t, err := h.repository.GetTaskByID(ctx, params.TaskID)
	if err != nil {
		logging.Errorf("repository.GetTaskByID failed: %v", err)
		return &ogen.UpdateTaskInternalServerError{
			Message: messageInternalServerError,
			Debug:   err.Error(),
		}, nil
	}
	if p.ID != t.ProjectID {
		logging.Errorf("p.ID != t.ProjectID")
		return &ogen.UpdateTaskNotFound{
			Message: messageNotFound,
			Debug:   err.Error(),
		}, nil
	}

	if !req.IsCompleted.IsSet() {
		logging.Errorf("value contains nothing")
		return &ogen.UpdateTaskBadRequest{
			Message: messageBadRequest,
			Debug:   err.Error(),
		}, nil
	}
	now := time.Now()
	if req.IsCompleted.Value {
		t.CompletedAt = &now
	} else {
		t.CompletedAt = nil
	}
	t.UpdatedAt = now
	if err := h.repository.UpdateTask(ctx, params.TaskID, t.CompletedAt, t.UpdatedAt); err != nil {
		logging.Errorf("repository.UpdateTask failed: %v", err)
		return &ogen.UpdateTaskInternalServerError{
			Message: messageInternalServerError,
			Debug:   err.Error(),
		}, nil
	}

	return newTaskResponse(t), nil
}

// DeleteTask は DELETE /projects/{projectID}/tasks/{taskID} に対応するハンドラ
func (h *handler) DeleteTask(ctx context.Context, params ogen.DeleteTaskParams) (ogen.DeleteTaskRes, error) {
	u, ok := ctx.Value(userKey{}).(*User)
	if !ok {
		logging.Errorf("ctx.Value(userKey{}).(*User) failed")
		return &ogen.DeleteTaskInternalServerError{
			Message: messageInternalServerError,
			Debug:   "ctx.Value(userKey{}).(*User) failed",
		}, nil
	}

	p, err := h.repository.GetProjectByID(ctx, params.ProjectID)
	if err != nil {
		logging.Errorf("repository.GetProjectByID failed: %v", err)
		return &ogen.DeleteTaskInternalServerError{
			Message: messageInternalServerError,
			Debug:   err.Error(),
		}, nil
	}
	if u.ID != p.UserID {
		logging.Errorf("u.ID != p.UserID")
		return &ogen.DeleteTaskNotFound{
			Message: messageNotFound,
			Debug:   err.Error(),
		}, nil
	}
	t, err := h.repository.GetTaskByID(ctx, params.TaskID)
	if err != nil {
		logging.Errorf("repository.GetTaskByID failed: %v", err)
		return &ogen.DeleteTaskInternalServerError{
			Message: messageInternalServerError,
			Debug:   err.Error(),
		}, nil
	}
	if p.ID != t.ProjectID {
		logging.Errorf("p.ID != t.ProjectID")
		return &ogen.DeleteTaskNotFound{
			Message: messageNotFound,
			Debug:   err.Error(),
		}, nil
	}

	if err := h.repository.DeleteTask(ctx, t.ID); err != nil {
		logging.Errorf("repository.DeleteTask failed: %v", err)
		return &ogen.DeleteTaskInternalServerError{
			Message: messageInternalServerError,
			Debug:   err.Error(),
		}, nil
	}

	return &ogen.DeleteTaskNoContent{}, nil
}

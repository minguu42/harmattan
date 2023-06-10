package app

import (
	"context"
	"fmt"
	"net/url"
	"time"

	"github.com/minguu42/mtasks/app/logging"
	"github.com/minguu42/mtasks/app/ogen"
)

// PostTasks は POST /projects/{projectID}/tasks に対応するハンドラ
func (h *handler) PostTasks(ctx context.Context, req *ogen.PostTasksReq, params ogen.PostTasksParams) (ogen.PostTasksRes, error) {
	u, ok := ctx.Value(userKey{}).(*User)
	if !ok {
		logging.Errorf("ctx.Value(userKey{}).(*User) failed")
		return &ogen.PostTasksInternalServerError{}, nil
	}

	p, err := h.repository.GetProjectByID(ctx, params.ProjectID)
	if err != nil {
		logging.Errorf("repository.GetProjectByID failed: %v", err)
		return &ogen.PostTasksInternalServerError{}, nil
	}
	if u.ID != p.UserID {
		logging.Errorf("u.ID != p.UserID")
		return &ogen.PostTasksNotFound{}, nil
	}

	t, err := h.repository.CreateTask(ctx, params.ProjectID, req.Title)
	if err != nil {
		logging.Errorf("repository.Create failed: %v", err)
		return &ogen.PostTasksInternalServerError{}, nil
	}

	location, err := url.ParseRequestURI(fmt.Sprintf("http://localhost:8080/projects/%d/tasks/%d", params.ProjectID, t.ID))
	if err != nil {
		logging.Errorf("url.ParseRequestURI failed: %v", err)
		return &ogen.PostTasksInternalServerError{}, nil
	}
	return &ogen.TaskHeaders{
		Location: *location,
		Response: newTaskResponse(t),
	}, nil
}

// GetTasks は GET /projects/{projectID}/tasks に対応するハンドラ
func (h *handler) GetTasks(ctx context.Context, params ogen.GetTasksParams) (ogen.GetTasksRes, error) {
	u, ok := ctx.Value(userKey{}).(*User)
	if !ok {
		logging.Errorf("ctx.Value(userKey{}).(*User) failed")
		return &ogen.GetTasksInternalServerError{}, nil
	}

	p, err := h.repository.GetProjectByID(ctx, params.ProjectID)
	if err != nil {
		logging.Errorf("repository.GetProjectByID failed: %v", err)
		return &ogen.GetTasksInternalServerError{}, nil
	}
	if u.ID != p.UserID {
		logging.Errorf("u.ID != p.UserID")
		return &ogen.GetTasksNotFound{}, nil
	}

	ts, err := h.repository.GetTasksByProjectID(ctx, p.ID, string(params.Sort.Or(ogen.GetTasksSortMinusCreatedAt)), params.Limit.Or(10), params.Offset.Or(0))
	if err != nil {
		logging.Errorf("repository.GetTasksByProjectID failed: %v", err)
		return &ogen.GetTasksInternalServerError{}, nil
	}

	return &ogen.Tasks{Tasks: newTasksResponse(ts)}, nil
}

// PatchTask は PATCH /projects/{projectID}/tasks/{taskID} に対応するハンドラ
func (h *handler) PatchTask(ctx context.Context, req *ogen.PatchTaskReq, params ogen.PatchTaskParams) (ogen.PatchTaskRes, error) {
	u, ok := ctx.Value(userKey{}).(*User)
	if !ok {
		logging.Errorf("ctx.Value(userKey{}).(*User) failed")
		return &ogen.PatchTaskInternalServerError{}, nil
	}

	p, err := h.repository.GetProjectByID(ctx, params.ProjectID)
	if err != nil {
		logging.Errorf("repository.GetProjectByID failed: %v", err)
		return &ogen.PatchTaskInternalServerError{}, nil
	}
	if u.ID != p.UserID {
		logging.Errorf("u.ID != p.UserID")
		return &ogen.PatchTaskNotFound{}, nil
	}
	t, err := h.repository.GetTaskByID(ctx, params.TaskID)
	if err != nil {
		logging.Errorf("repository.GetTaskByID failed: %v", err)
		return &ogen.PatchTaskInternalServerError{}, nil
	}
	if p.ID != t.ProjectID {
		logging.Errorf("p.ID != t.ProjectID")
		return &ogen.PatchTaskNotFound{}, nil
	}

	if !req.IsCompleted.IsSet() {
		logging.Errorf("value contains nothing")
		return &ogen.PatchTaskBadRequest{}, nil
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
		return &ogen.PatchTaskInternalServerError{}, nil
	}

	resp := newTaskResponse(t)
	return &resp, nil
}

// DeleteTask は DELETE /projects/{projectID}/tasks/{taskID} に対応するハンドラ
func (h *handler) DeleteTask(ctx context.Context, params ogen.DeleteTaskParams) (ogen.DeleteTaskRes, error) {
	u, ok := ctx.Value(userKey{}).(*User)
	if !ok {
		logging.Errorf("ctx.Value(userKey{}).(*User) failed")
		return &ogen.DeleteTaskInternalServerError{}, nil
	}

	p, err := h.repository.GetProjectByID(ctx, params.ProjectID)
	if err != nil {
		logging.Errorf("repository.GetProjectByID failed: %v", err)
		return &ogen.DeleteTaskInternalServerError{}, nil
	}
	if u.ID != p.UserID {
		logging.Errorf("u.ID != p.UserID")
		return &ogen.DeleteTaskNotFound{}, nil
	}
	t, err := h.repository.GetTaskByID(ctx, params.TaskID)
	if err != nil {
		logging.Errorf("repository.GetTaskByID failed: %v", err)
		return &ogen.DeleteTaskInternalServerError{}, nil
	}
	if p.ID != t.ProjectID {
		logging.Errorf("p.ID != t.ProjectID")
		return &ogen.DeleteTaskNotFound{}, nil
	}

	if err := h.repository.DeleteTask(ctx, t.ID); err != nil {
		logging.Errorf("repository.DeleteTask failed: %v", err)
		return &ogen.DeleteTaskInternalServerError{}, nil
	}

	return &ogen.DeleteTaskNoContent{}, nil
}

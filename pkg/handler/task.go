package handler

import (
	"context"
	"errors"

	"github.com/minguu42/opepe/gen/ogen"
	"github.com/minguu42/opepe/pkg/entity"
	"github.com/minguu42/opepe/pkg/logging"
	"github.com/minguu42/opepe/pkg/repository"
	"github.com/minguu42/opepe/pkg/ttime"
)

// CreateTask は POST /projects/{projectID}/tasks に対応するハンドラ
func (h *Handler) CreateTask(ctx context.Context, req *ogen.CreateTaskReq, params ogen.CreateTaskParams) (*ogen.Task, error) {
	u, ok := ctx.Value(userKey{}).(*entity.User)
	if !ok {
		return nil, errUnauthorized
	}

	p, err := h.Repository.GetProjectByID(ctx, params.ProjectID)
	if err != nil {
		if errors.Is(err, repository.ErrRecordNotFound) {
			return nil, errProjectNotFound
		}
		logging.Errorf(ctx, "repository.GetProjectByID failed: %v", err)
		return nil, errInternalServerError
	}
	if !u.HasProject(p) {
		logging.Errorf(ctx, "user does not have the project")
		return nil, errProjectNotFound
	}

	t := entity.Task{
		ProjectID: params.ProjectID,
		Title:     req.Title,
		Content:   req.Content,
		Priority:  req.Priority,
		DueOn:     req.DueOn.Ptr(),
	}
	if err := h.Repository.SaveTask(ctx, &t); err != nil {
		logging.Errorf(ctx, "repository.SaveTask failed: %s", err)
		return nil, errInternalServerError
	}

	return newTaskResponse(&t), nil
}

// ListTasks は GET /projects/{projectID}/tasks に対応するハンドラ
func (h *Handler) ListTasks(ctx context.Context, params ogen.ListTasksParams) (*ogen.Tasks, error) {
	u, ok := ctx.Value(userKey{}).(*entity.User)
	if !ok {
		return nil, errUnauthorized
	}

	p, err := h.Repository.GetProjectByID(ctx, params.ProjectID)
	if err != nil {
		if errors.Is(err, repository.ErrRecordNotFound) {
			return nil, errProjectNotFound
		}
		logging.Errorf(ctx, "repository.GetProjectByID failed: %v", err)
		return nil, errInternalServerError
	}
	if !u.HasProject(p) {
		logging.Errorf(ctx, "user does not have the project")
		return nil, errProjectNotFound
	}

	limit := params.Limit.Or(defaultLimit)
	ts, err := h.Repository.GetTasksByProjectID(ctx, p.ID, limit+1, params.Offset.Or(defaultOffset))
	if err != nil {
		logging.Errorf(ctx, "repository.GetTasksByProjectID failed: %v", err)
		return nil, errInternalServerError
	}

	hasNext := false
	if len(ts) == limit+1 {
		hasNext = true
		ts = ts[:limit]
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
		return nil, errUnauthorized
	}

	p, err := h.Repository.GetProjectByID(ctx, params.ProjectID)
	if err != nil {
		if errors.Is(err, repository.ErrRecordNotFound) {
			return nil, errProjectNotFound
		}
		logging.Errorf(ctx, "repository.GetProjectByID failed: %v", err)
		return nil, errInternalServerError
	}
	if !u.HasProject(p) {
		logging.Errorf(ctx, "user does not have the project")
		return nil, errProjectNotFound
	}
	t, err := h.Repository.GetTaskByID(ctx, params.TaskID)
	if err != nil {
		if errors.Is(err, repository.ErrRecordNotFound) {
			return nil, errTaskNotFound
		}
		logging.Errorf(ctx, "repository.GetTaskByID failed: %v", err)
		return nil, errInternalServerError
	}
	if !p.ContainsTask(t) {
		logging.Errorf(ctx, "project does not contain the task")
		return nil, errTaskNotFound
	}

	dueOn := t.DueOn
	if req.DueOn.IsSet() {
		dueOn = req.DueOn.Ptr()
	}
	newTask := entity.Task{
		ID:          t.ID,
		ProjectID:   t.ProjectID,
		Title:       req.Title.Or(t.Title),
		Content:     req.Content.Or(t.Content),
		Priority:    req.Priority.Or(t.Priority),
		DueOn:       dueOn,
		CompletedAt: nil,
		CreatedAt:   t.CreatedAt,
		UpdatedAt:   ttime.Now(ctx),
	}
	if err := h.Repository.SaveTask(ctx, &newTask); err != nil {
		logging.Errorf(ctx, "repository.SaveTask failed: %s", err)
		return nil, errInternalServerError
	}

	return newTaskResponse(&newTask), nil
}

// DeleteTask は DELETE /projects/{projectID}/tasks/{taskID} に対応するハンドラ
func (h *Handler) DeleteTask(ctx context.Context, params ogen.DeleteTaskParams) error {
	u, ok := ctx.Value(userKey{}).(*entity.User)
	if !ok {
		return errUnauthorized
	}

	p, err := h.Repository.GetProjectByID(ctx, params.ProjectID)
	if err != nil {
		if errors.Is(err, repository.ErrRecordNotFound) {
			return errProjectNotFound
		}
		logging.Errorf(ctx, "repository.GetProjectByID failed: %v", err)
		return errInternalServerError
	}
	if !u.HasProject(p) {
		logging.Errorf(ctx, "user does not have the project")
		return errProjectNotFound
	}
	t, err := h.Repository.GetTaskByID(ctx, params.TaskID)
	if err != nil {
		if errors.Is(err, repository.ErrRecordNotFound) {
			return errTaskNotFound
		}
		logging.Errorf(ctx, "repository.GetTaskByID failed: %v", err)
		return errInternalServerError
	}
	if !p.ContainsTask(t) {
		logging.Errorf(ctx, "project does not contain the task")
		return errTaskNotFound
	}

	if err := h.Repository.DeleteTask(ctx, t.ID); err != nil {
		logging.Errorf(ctx, "repository.DeleteTask failed: %v", err)
		return errInternalServerError
	}

	return nil
}

package handler

import (
	"context"
	"fmt"
	"time"

	"github.com/minguu42/harmattan/api/apperr"
	"github.com/minguu42/harmattan/api/usecase"
	"github.com/minguu42/harmattan/internal/domain"
	"github.com/minguu42/harmattan/internal/lib/opt"
	"github.com/minguu42/harmattan/internal/lib/pointers"
	"github.com/minguu42/harmattan/internal/openapi"
)

func (h *handler) CreateTask(ctx context.Context, req *openapi.CreateTaskReq, params openapi.CreateTaskParams) (*openapi.Task, error) {
	var errs []error
	errs = append(errs, validateTaskName(req.Name)...)
	if len(errs) > 0 {
		return nil, apperr.DomainValidationError(errs)
	}

	out, err := h.task.CreateTask(ctx, &usecase.CreateTaskInput{
		ProjectID: domain.ProjectID(params.ProjectID),
		Name:      req.Name,
		Priority:  req.Priority,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to execute CreateTask usecase: %w", err)
	}
	return convertTask(out.Task), nil
}

func (h *handler) ListTasks(ctx context.Context, params openapi.ListTasksParams) (*openapi.ListTasksOK, error) {
	out, err := h.task.ListTasks(ctx, &usecase.ListTasksInput{
		ProjectID: domain.ProjectID(params.ProjectID),
		Limit:     params.Limit.Value,
		Offset:    params.Offset.Value,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to execute ListTasks usecase: %w", err)
	}
	return &openapi.ListTasksOK{
		Tasks:   convertTasks(out.Tasks),
		HasNext: out.HasNext,
	}, nil
}

func (h *handler) UpdateTask(ctx context.Context, req *openapi.UpdateTaskReq, params openapi.UpdateTaskParams) (*openapi.Task, error) {
	var errs []error
	if name, ok := req.Name.Get(); ok {
		errs = append(errs, validateTaskName(name)...)
	}
	if len(errs) > 0 {
		return nil, apperr.DomainValidationError(errs)
	}

	out, err := h.task.UpdateTask(ctx, &usecase.UpdateTaskInput{
		ID:          domain.TaskID(params.TaskID),
		ProjectID:   domain.ProjectID(params.ProjectID),
		Name:        opt.Option[string]{V: req.Name.Value, Valid: req.Name.Set},
		Content:     opt.Option[string]{V: req.Content.Value, Valid: req.Content.Set},
		Priority:    opt.Option[int]{V: req.Priority.Value, Valid: req.Priority.Set},
		DueOn:       opt.Option[*time.Time]{V: pointers.RefOrNil(req.DueOn.Null, req.DueOn.Value), Valid: req.DueOn.Set},
		CompletedAt: opt.Option[*time.Time]{V: pointers.RefOrNil(req.CompletedAt.Null, req.CompletedAt.Value), Valid: req.CompletedAt.Set},
	})
	if err != nil {
		return nil, fmt.Errorf("failed to execute UpdateTask usecase: %w", err)
	}
	return convertTask(out.Task), nil
}

func (h *handler) DeleteTask(ctx context.Context, params openapi.DeleteTaskParams) error {
	if err := h.task.DeleteTask(ctx, &usecase.DeleteTaskInput{
		ProjectID: domain.ProjectID(params.ProjectID),
		ID:        domain.TaskID(params.TaskID),
	}); err != nil {
		return fmt.Errorf("failed to execute DeleteTask usecase: %w", err)
	}
	return nil
}

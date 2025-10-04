package handler

import (
	"context"
	"fmt"
	"time"

	"github.com/minguu42/harmattan/internal/api/handler/openapi"
	"github.com/minguu42/harmattan/internal/api/usecase"
	"github.com/minguu42/harmattan/internal/domain"
	"github.com/minguu42/harmattan/internal/lib/pointers"
)

func (h *handler) CreateTask(ctx context.Context, req *openapi.CreateTaskReq, params openapi.CreateTaskParams) (*openapi.Task, error) {
	var errs []error
	errs = append(errs, validateTaskName(req.Name)...)
	if len(errs) > 0 {
		return nil, usecase.DomainValidationError(errs)
	}

	out, err := h.task.CreateTask(ctx, &usecase.CreateTaskInput{
		ProjectID: domain.ProjectID(params.ProjectID),
		Name:      req.Name,
		Priority:  req.Priority,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to execute CreateTask usecase: %w", err)
	}
	return convertTask(out.Task, out.Tags), nil
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
		Tasks:   convertTasks(out.Tasks, out.Tags),
		HasNext: out.HasNext,
	}, nil
}

func (h *handler) GetTask(ctx context.Context, params openapi.GetTaskParams) (*openapi.Task, error) {
	out, err := h.task.GetTask(ctx, &usecase.GetTaskInput{ID: domain.TaskID(params.TaskID)})
	if err != nil {
		return nil, fmt.Errorf("failed to execute GetTask usecase: %w", err)
	}
	return convertTask(out.Task, out.Tags), nil
}

func (h *handler) UpdateTask(ctx context.Context, req *openapi.UpdateTaskReq, params openapi.UpdateTaskParams) (*openapi.Task, error) {
	var errs []error
	if name, ok := req.Name.Get(); ok {
		errs = append(errs, validateTaskName(name)...)
	}
	if len(errs) > 0 {
		return nil, usecase.DomainValidationError(errs)
	}

	out, err := h.task.UpdateTask(ctx, &usecase.UpdateTaskInput{
		ID:          domain.TaskID(params.TaskID),
		Name:        usecase.Option[string]{V: req.Name.Value, Valid: req.Name.Set},
		TagIDs:      usecase.Option[[]domain.TagID]{V: convertSlice[domain.TagID](req.TagIds), Valid: req.TagIds != nil},
		Content:     usecase.Option[string]{V: req.Content.Value, Valid: req.Content.Set},
		Priority:    usecase.Option[int]{V: req.Priority.Value, Valid: req.Priority.Set},
		DueOn:       usecase.Option[*time.Time]{V: pointers.Ternary(req.DueOn.Null, nil, &req.DueOn.Value), Valid: req.DueOn.Set},
		CompletedAt: usecase.Option[*time.Time]{V: pointers.Ternary(req.CompletedAt.Null, nil, &req.CompletedAt.Value), Valid: req.CompletedAt.Set},
	})
	if err != nil {
		return nil, fmt.Errorf("failed to execute UpdateTask usecase: %w", err)
	}
	return convertTask(out.Task, out.Tags), nil
}

func (h *handler) DeleteTask(ctx context.Context, params openapi.DeleteTaskParams) error {
	if err := h.task.DeleteTask(ctx, &usecase.DeleteTaskInput{ID: domain.TaskID(params.TaskID)}); err != nil {
		return fmt.Errorf("failed to execute DeleteTask usecase: %w", err)
	}
	return nil
}

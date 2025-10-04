package handler

import (
	"context"
	"fmt"
	"time"

	openapi2 "github.com/minguu42/harmattan/internal/api/handler/openapi"
	usecase2 "github.com/minguu42/harmattan/internal/api/usecase"
	"github.com/minguu42/harmattan/internal/domain"
	"github.com/minguu42/harmattan/internal/lib/pointers"
)

func (h *handler) CreateTask(ctx context.Context, req *openapi2.CreateTaskReq, params openapi2.CreateTaskParams) (*openapi2.Task, error) {
	var errs []error
	errs = append(errs, validateTaskName(req.Name)...)
	if len(errs) > 0 {
		return nil, usecase2.DomainValidationError(errs)
	}

	out, err := h.task.CreateTask(ctx, &usecase2.CreateTaskInput{
		ProjectID: domain.ProjectID(params.ProjectID),
		Name:      req.Name,
		Priority:  req.Priority,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to execute CreateTask usecase: %w", err)
	}
	return convertTask(out.Task, out.Tags), nil
}

func (h *handler) ListTasks(ctx context.Context, params openapi2.ListTasksParams) (*openapi2.ListTasksOK, error) {
	out, err := h.task.ListTasks(ctx, &usecase2.ListTasksInput{
		ProjectID: domain.ProjectID(params.ProjectID),
		Limit:     params.Limit.Value,
		Offset:    params.Offset.Value,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to execute ListTasks usecase: %w", err)
	}
	return &openapi2.ListTasksOK{
		Tasks:   convertTasks(out.Tasks, out.Tags),
		HasNext: out.HasNext,
	}, nil
}

func (h *handler) GetTask(ctx context.Context, params openapi2.GetTaskParams) (*openapi2.Task, error) {
	out, err := h.task.GetTask(ctx, &usecase2.GetTaskInput{ID: domain.TaskID(params.TaskID)})
	if err != nil {
		return nil, fmt.Errorf("failed to execute GetTask usecase: %w", err)
	}
	return convertTask(out.Task, out.Tags), nil
}

func (h *handler) UpdateTask(ctx context.Context, req *openapi2.UpdateTaskReq, params openapi2.UpdateTaskParams) (*openapi2.Task, error) {
	var errs []error
	if name, ok := req.Name.Get(); ok {
		errs = append(errs, validateTaskName(name)...)
	}
	if len(errs) > 0 {
		return nil, usecase2.DomainValidationError(errs)
	}

	out, err := h.task.UpdateTask(ctx, &usecase2.UpdateTaskInput{
		ID:          domain.TaskID(params.TaskID),
		Name:        usecase2.Option[string]{V: req.Name.Value, Valid: req.Name.Set},
		TagIDs:      usecase2.Option[[]domain.TagID]{V: convertSlice[domain.TagID](req.TagIds), Valid: req.TagIds != nil},
		Content:     usecase2.Option[string]{V: req.Content.Value, Valid: req.Content.Set},
		Priority:    usecase2.Option[int]{V: req.Priority.Value, Valid: req.Priority.Set},
		DueOn:       usecase2.Option[*time.Time]{V: pointers.Ternary(req.DueOn.Null, nil, &req.DueOn.Value), Valid: req.DueOn.Set},
		CompletedAt: usecase2.Option[*time.Time]{V: pointers.Ternary(req.CompletedAt.Null, nil, &req.CompletedAt.Value), Valid: req.CompletedAt.Set},
	})
	if err != nil {
		return nil, fmt.Errorf("failed to execute UpdateTask usecase: %w", err)
	}
	return convertTask(out.Task, out.Tags), nil
}

func (h *handler) DeleteTask(ctx context.Context, params openapi2.DeleteTaskParams) error {
	if err := h.task.DeleteTask(ctx, &usecase2.DeleteTaskInput{ID: domain.TaskID(params.TaskID)}); err != nil {
		return fmt.Errorf("failed to execute DeleteTask usecase: %w", err)
	}
	return nil
}

package handler

import (
	"context"
	"fmt"

	"github.com/minguu42/harmattan/api/usecase"
	"github.com/minguu42/harmattan/internal/domain"
	"github.com/minguu42/harmattan/internal/openapi"
)

func convertTask(task *domain.Task) *openapi.Task {
	return &openapi.Task{
		ID:          string(task.ID),
		ProjectID:   string(task.ProjectID),
		Name:        task.Name,
		Content:     task.Content,
		Priority:    task.Priority,
		DueOn:       convertDatePtr(task.DueOn),
		CompletedAt: convertDateTimePtr(task.CompletedAt),
		CreatedAt:   task.CreatedAt,
		UpdatedAt:   task.UpdatedAt,
		Steps:       convertSteps(task.Steps),
		Tags:        convertTags(task.Tags),
	}
}

func convertTasks(tasks domain.Tasks) []openapi.Task {
	ts := make([]openapi.Task, 0, len(tasks))
	for _, t := range tasks {
		ts = append(ts, *convertTask(&t))
	}
	return ts
}

func (h *handler) CreateTask(ctx context.Context, req *openapi.CreateTaskReq, params openapi.CreateTaskParams) (*openapi.Task, error) {
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

func (h *handler) ListTasks(ctx context.Context, params openapi.ListTasksParams) (*openapi.Tasks, error) {
	out, err := h.task.ListTasks(ctx, &usecase.ListTasksInput{
		ProjectID: domain.ProjectID(params.ProjectID),
		Limit:     params.Limit.Value,
		Offset:    params.Offset.Value,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to execute ListTasks usecase: %w", err)
	}
	return &openapi.Tasks{
		Tasks:   convertTasks(out.Tasks),
		HasNext: out.HasNext,
	}, nil
}

func (h *handler) UpdateTask(ctx context.Context, req *openapi.UpdateTaskReq, params openapi.UpdateTaskParams) (*openapi.Task, error) {
	out, err := h.task.UpdateTask(ctx, &usecase.UpdateTaskInput{
		ID:          domain.TaskID(params.TaskID),
		ProjectID:   domain.ProjectID(params.ProjectID),
		Name:        convertOptString(req.Name),
		Content:     convertOptString(req.Content),
		Priority:    convertOptInt(req.Priority),
		DueOn:       convertOptDateTime(req.DueOn),
		CompletedAt: convertOptDateTime(req.CompletedAt),
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

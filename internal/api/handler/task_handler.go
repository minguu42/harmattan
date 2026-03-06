package handler

import (
	"context"
	"errors"
	"time"
	"unicode/utf8"

	"github.com/minguu42/harmattan/internal/api/apierror"
	"github.com/minguu42/harmattan/internal/api/openapi"
	"github.com/minguu42/harmattan/internal/api/usecase"
	"github.com/minguu42/harmattan/internal/domain"
	"github.com/minguu42/harmattan/internal/lib/errtrace"
)

func (h *Handler) CreateTask(ctx context.Context, req *openapi.CreateTaskReq, params openapi.CreateTaskParams) (*openapi.Task, error) {
	var errs []error
	errs = append(errs, validateTaskName(req.Name)...)
	if len(errs) > 0 {
		return nil, errtrace.Wrap(apierror.DomainValidationError(errs))
	}

	out, err := h.Task.CreateTask(ctx, &usecase.CreateTaskInput{
		ProjectID: domain.ProjectID(params.ProjectID),
		Name:      req.Name,
		Priority:  req.Priority.Value,
	})
	if err != nil {
		return nil, errtrace.Wrap(err)
	}
	return convertTask(out.Task, out.Tags), nil
}

func (h *Handler) ListTasks(ctx context.Context, params openapi.ListTasksParams) (*openapi.ListTasksOK, error) {
	out, err := h.Task.ListTasks(ctx, &usecase.ListTasksInput{
		ProjectID:     domain.ProjectID(params.ProjectID),
		Limit:         params.Limit.Value,
		Offset:        params.Offset.Value,
		ShowCompleted: params.ShowCompleted.Value,
	})
	if err != nil {
		return nil, errtrace.Wrap(err)
	}
	return &openapi.ListTasksOK{
		Tasks:   convertTasks(out.Tasks, out.Tags),
		HasNext: out.HasNext,
	}, nil
}

func (h *Handler) GetTask(ctx context.Context, params openapi.GetTaskParams) (*openapi.Task, error) {
	out, err := h.Task.GetTask(ctx, &usecase.GetTaskInput{ID: domain.TaskID(params.TaskID)})
	if err != nil {
		return nil, errtrace.Wrap(err)
	}
	return convertTask(out.Task, out.Tags), nil
}

func (h *Handler) UpdateTask(ctx context.Context, req *openapi.UpdateTaskReq, params openapi.UpdateTaskParams) (*openapi.Task, error) {
	var errs []error
	if name, ok := req.Name.Get(); ok {
		errs = append(errs, validateTaskName(name)...)
	}
	if len(errs) > 0 {
		return nil, errtrace.Wrap(apierror.DomainValidationError(errs))
	}

	out, err := h.Task.UpdateTask(ctx, &usecase.UpdateTaskInput{
		ID:          domain.TaskID(params.TaskID),
		Name:        usecase.Option[string]{V: req.Name.Value, Valid: req.Name.Set},
		TagIDs:      usecase.Option[[]domain.TagID]{V: convertSlice[domain.TagID](req.TagIds), Valid: req.TagIds != nil},
		Content:     usecase.Option[string]{V: req.Content.Value, Valid: req.Content.Set},
		Priority:    usecase.Option[int]{V: req.Priority.Value, Valid: req.Priority.Set},
		DueOn:       usecase.Option[*time.Time]{V: ternary(req.DueOn.Null, nil, &req.DueOn.Value), Valid: req.DueOn.Set},
		CompletedAt: usecase.Option[*time.Time]{V: ternary(req.CompletedAt.Null, nil, &req.CompletedAt.Value), Valid: req.CompletedAt.Set},
	})
	if err != nil {
		return nil, errtrace.Wrap(err)
	}
	return convertTask(out.Task, out.Tags), nil
}

func (h *Handler) DeleteTask(ctx context.Context, params openapi.DeleteTaskParams) error {
	if err := h.Task.DeleteTask(ctx, &usecase.DeleteTaskInput{ID: domain.TaskID(params.TaskID)}); err != nil {
		return errtrace.Wrap(err)
	}
	return nil
}

var ErrTaskNameLength = errors.New("タスク名は1文字以上100文字以下で指定できます")

func validateTaskName(name string) []error {
	var errs []error
	if utf8.RuneCountInString(name) < 1 || 100 < utf8.RuneCountInString(name) {
		errs = append(errs, ErrTaskNameLength)
	}
	return errs
}

func convertTask(task *domain.Task, tags domain.Tags) *openapi.Task {
	return &openapi.Task{
		ID:          string(task.ID),
		ProjectID:   string(task.ProjectID),
		Name:        task.Name,
		Content:     task.Content,
		Priority:    task.Priority,
		DueOn:       convertOptDate(task.DueOn),
		CompletedAt: convertOptDateTime(task.CompletedAt),
		CreatedAt:   task.CreatedAt,
		UpdatedAt:   task.UpdatedAt,
		Steps:       convertSteps(task.Steps),
		Tags:        convertTags(tags),
	}
}

func convertTasks(tasks domain.Tasks, tags domain.Tags) []openapi.Task {
	tagByID := tags.TagByID()

	ts := make([]openapi.Task, 0, len(tasks))
	for _, t := range tasks {
		taskTags := make(domain.Tags, 0, len(t.TagIDs))
		for _, id := range t.TagIDs {
			taskTags = append(taskTags, tagByID[id])
		}
		ts = append(ts, *convertTask(&t, taskTags))
	}
	return ts
}

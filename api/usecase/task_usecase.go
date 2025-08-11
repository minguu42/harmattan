package usecase

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/minguu42/harmattan/api/apperr"
	"github.com/minguu42/harmattan/internal/auth"
	"github.com/minguu42/harmattan/internal/database"
	"github.com/minguu42/harmattan/internal/domain"
	"github.com/minguu42/harmattan/internal/lib/clock"
	"github.com/minguu42/harmattan/internal/lib/idgen"
	"github.com/minguu42/harmattan/internal/lib/opt"
)

type Task struct {
	DB *database.Client
}

type TaskOutput struct {
	Task *domain.Task
}

type CreateTaskInput struct {
	ProjectID domain.ProjectID
	Name      string
	Priority  int
}

func (uc *Task) CreateTask(ctx context.Context, in *CreateTaskInput) (*TaskOutput, error) {
	user := auth.MustUserFromContext(ctx)

	p, err := uc.DB.GetProjectByID(ctx, in.ProjectID)
	if err != nil {
		if errors.Is(err, database.ErrModelNotFound) {
			return nil, apperr.ProjectNotFoundError(err)
		}
		return nil, fmt.Errorf("failed to get project: %w", err)
	}
	if !user.HasProject(p) {
		return nil, apperr.ProjectAccessDeniedError()
	}

	now := clock.Now(ctx)
	t := domain.Task{
		ID:        domain.TaskID(idgen.ULID(ctx)),
		UserID:    user.ID,
		ProjectID: in.ProjectID,
		Name:      in.Name,
		Priority:  in.Priority,
		CreatedAt: now,
		UpdatedAt: now,
	}
	if err := uc.DB.CreateTask(ctx, &t); err != nil {
		return nil, fmt.Errorf("failed to create task: %w", err)
	}
	return &TaskOutput{Task: &t}, nil
}

type ListTasksInput struct {
	ProjectID domain.ProjectID
	Limit     int
	Offset    int
}

type ListTasksOutput struct {
	Tasks   domain.Tasks
	HasNext bool
}

func (uc *Task) ListTasks(ctx context.Context, in *ListTasksInput) (*ListTasksOutput, error) {
	user := auth.MustUserFromContext(ctx)

	p, err := uc.DB.GetProjectByID(ctx, in.ProjectID)
	if err != nil {
		if errors.Is(err, database.ErrModelNotFound) {
			return nil, apperr.ProjectNotFoundError(err)
		}
		return nil, fmt.Errorf("failed to get project: %w", err)
	}
	if !user.HasProject(p) {
		return nil, apperr.ProjectAccessDeniedError()
	}

	ts, err := uc.DB.ListTasks(ctx, in.ProjectID, in.Limit+1, in.Offset)
	if err != nil {
		return nil, fmt.Errorf("failed to list tasks: %w", err)
	}

	hasNext := false
	if len(ts) == in.Limit+1 {
		ts = ts[:in.Limit]
		hasNext = true
	}
	return &ListTasksOutput{Tasks: ts, HasNext: hasNext}, nil
}

type GetTaskInput struct {
	ID domain.TaskID
}

func (uc *Task) GetTask(ctx context.Context, in *GetTaskInput) (*TaskOutput, error) {
	user := auth.MustUserFromContext(ctx)

	t, err := uc.DB.GetTaskByID(ctx, in.ID)
	if err != nil {
		if errors.Is(err, database.ErrModelNotFound) {
			return nil, apperr.TaskNotFoundError(err)
		}
		return nil, fmt.Errorf("failed to get task: %w", err)
	}
	if !user.HasTask(t) {
		return nil, apperr.TaskAccessDeniedError()
	}

	return &TaskOutput{Task: t}, nil
}

type UpdateTaskInput struct {
	ID          domain.TaskID
	Name        opt.Option[string]
	Content     opt.Option[string]
	Priority    opt.Option[int]
	DueOn       opt.Option[*time.Time]
	CompletedAt opt.Option[*time.Time]
}

func (uc *Task) UpdateTask(ctx context.Context, in *UpdateTaskInput) (*TaskOutput, error) {
	user := auth.MustUserFromContext(ctx)

	t, err := uc.DB.GetTaskByID(ctx, in.ID)
	if err != nil {
		if errors.Is(err, database.ErrModelNotFound) {
			return nil, apperr.TaskNotFoundError(err)
		}
		return nil, fmt.Errorf("failed to get task: %w", err)
	}
	if !user.HasTask(t) {
		return nil, apperr.TaskAccessDeniedError()
	}

	if in.Name.Valid {
		t.Name = in.Name.V
	}
	if in.Content.Valid {
		t.Content = in.Content.V
	}
	if in.Priority.Valid {
		t.Priority = in.Priority.V
	}
	if in.DueOn.Valid {
		t.DueOn = in.DueOn.V
	}
	if in.CompletedAt.Valid {
		t.CompletedAt = in.CompletedAt.V
	}
	t.UpdatedAt = clock.Now(ctx)
	if err := uc.DB.UpdateTask(ctx, t); err != nil {
		return nil, fmt.Errorf("failed to update task: %w", err)
	}
	return &TaskOutput{Task: t}, nil
}

type DeleteTaskInput struct {
	ID domain.TaskID
}

func (uc *Task) DeleteTask(ctx context.Context, in *DeleteTaskInput) error {
	user := auth.MustUserFromContext(ctx)

	t, err := uc.DB.GetTaskByID(ctx, in.ID)
	if err != nil {
		if errors.Is(err, database.ErrModelNotFound) {
			return apperr.TaskNotFoundError(err)
		}
		return fmt.Errorf("failed to get task: %w", err)
	}
	if !user.HasTask(t) {
		return apperr.TaskAccessDeniedError()
	}

	if err := uc.DB.DeleteTaskByID(ctx, t.ID); err != nil {
		return fmt.Errorf("failed to delete task: %w", err)
	}
	return nil
}

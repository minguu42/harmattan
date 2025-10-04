package usecase

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/minguu42/harmattan/internal/auth"
	"github.com/minguu42/harmattan/internal/database"
	"github.com/minguu42/harmattan/internal/domain"
	"github.com/minguu42/harmattan/internal/lib/clock"
	"github.com/minguu42/harmattan/internal/lib/idgen"
)

type Task struct {
	DB *database.Client
}

type TaskOutput struct {
	Task *domain.Task
	Tags domain.Tags
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
			return nil, ProjectNotFoundError(err)
		}
		return nil, fmt.Errorf("failed to get project: %w", err)
	}
	if !user.HasProject(p) {
		return nil, ProjectAccessDeniedError()
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
	Tags    domain.Tags
	HasNext bool
}

func (uc *Task) ListTasks(ctx context.Context, in *ListTasksInput) (*ListTasksOutput, error) {
	user := auth.MustUserFromContext(ctx)

	p, err := uc.DB.GetProjectByID(ctx, in.ProjectID)
	if err != nil {
		if errors.Is(err, database.ErrModelNotFound) {
			return nil, ProjectNotFoundError(err)
		}
		return nil, fmt.Errorf("failed to get project: %w", err)
	}
	if !user.HasProject(p) {
		return nil, ProjectAccessDeniedError()
	}

	ts, err := uc.DB.ListTasks(ctx, in.ProjectID, in.Limit+1, in.Offset)
	if err != nil {
		return nil, fmt.Errorf("failed to list tasks: %w", err)
	}

	tags, err := uc.DB.GetTagsByIDs(ctx, ts.TagIDs())
	if err != nil {
		return nil, fmt.Errorf("failed to get tags: %w", err)
	}

	hasNext := false
	if len(ts) == in.Limit+1 {
		ts = ts[:in.Limit]
		hasNext = true
	}
	return &ListTasksOutput{Tasks: ts, Tags: tags, HasNext: hasNext}, nil
}

type GetTaskInput struct {
	ID domain.TaskID
}

func (uc *Task) GetTask(ctx context.Context, in *GetTaskInput) (*TaskOutput, error) {
	user := auth.MustUserFromContext(ctx)

	task, err := uc.DB.GetTaskByID(ctx, in.ID)
	if err != nil {
		if errors.Is(err, database.ErrModelNotFound) {
			return nil, TaskNotFoundError(err)
		}
		return nil, fmt.Errorf("failed to get task: %w", err)
	}
	if !user.HasTask(task) {
		return nil, TaskAccessDeniedError()
	}

	tags, err := uc.DB.GetTagsByIDs(ctx, task.TagIDs)
	if err != nil {
		return nil, fmt.Errorf("failed to get tags: %w", err)
	}
	return &TaskOutput{Task: task, Tags: tags}, nil
}

type UpdateTaskInput struct {
	ID          domain.TaskID
	Name        Option[string]
	TagIDs      Option[[]domain.TagID]
	Content     Option[string]
	Priority    Option[int]
	DueOn       Option[*time.Time]
	CompletedAt Option[*time.Time]
}

func (uc *Task) UpdateTask(ctx context.Context, in *UpdateTaskInput) (*TaskOutput, error) {
	user := auth.MustUserFromContext(ctx)

	task, err := uc.DB.GetTaskByID(ctx, in.ID)
	if err != nil {
		if errors.Is(err, database.ErrModelNotFound) {
			return nil, TaskNotFoundError(err)
		}
		return nil, fmt.Errorf("failed to get task: %w", err)
	}
	if !user.HasTask(task) {
		return nil, TaskAccessDeniedError()
	}

	if in.Name.Valid {
		task.Name = in.Name.V
	}
	var tags domain.Tags
	if in.TagIDs.Valid {
		tags, err = uc.DB.GetTagsByIDs(ctx, in.TagIDs.V)
		if err != nil {
			return nil, fmt.Errorf("failed to get tags: %w", err)
		}
		validTags := make(domain.Tags, 0, len(tags))
		for _, t := range tags {
			if user.HasTag(&t) {
				validTags = append(validTags, t)
			}
		}
		tags = validTags
		task.TagIDs = validTags.IDs()
	} else {
		tags, err = uc.DB.GetTagsByIDs(ctx, task.TagIDs)
		if err != nil {
			return nil, fmt.Errorf("failed to get tags: %w", err)
		}
	}
	if in.Content.Valid {
		task.Content = in.Content.V
	}
	if in.Priority.Valid {
		task.Priority = in.Priority.V
	}
	if in.DueOn.Valid {
		task.DueOn = in.DueOn.V
	}
	if in.CompletedAt.Valid {
		task.CompletedAt = in.CompletedAt.V
	}
	task.UpdatedAt = clock.Now(ctx)
	if err := uc.DB.UpdateTask(ctx, task); err != nil {
		return nil, fmt.Errorf("failed to update task: %w", err)
	}
	return &TaskOutput{Task: task, Tags: tags}, nil
}

type DeleteTaskInput struct {
	ID domain.TaskID
}

func (uc *Task) DeleteTask(ctx context.Context, in *DeleteTaskInput) error {
	user := auth.MustUserFromContext(ctx)

	task, err := uc.DB.GetTaskByID(ctx, in.ID)
	if err != nil {
		if errors.Is(err, database.ErrModelNotFound) {
			return TaskNotFoundError(err)
		}
		return fmt.Errorf("failed to get task: %w", err)
	}
	if !user.HasTask(task) {
		return TaskAccessDeniedError()
	}

	if err := uc.DB.DeleteTaskByID(ctx, task.ID); err != nil {
		return fmt.Errorf("failed to delete task: %w", err)
	}
	return nil
}

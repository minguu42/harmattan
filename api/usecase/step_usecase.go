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
	"github.com/minguu42/harmattan/lib/clock"
	"github.com/minguu42/harmattan/lib/idgen"
	"github.com/minguu42/harmattan/lib/pointers"
)

type Step struct {
	DB *database.Client
}

type StepOutput struct {
	Step *domain.Step
}

type CreateStepInput struct {
	TaskID domain.TaskID
	Name   string
}

func (uc *Step) CreateStep(ctx context.Context, in *CreateStepInput) (*StepOutput, error) {
	user := auth.MustUserFromContext(ctx)

	task, err := uc.DB.GetTaskByID(ctx, in.TaskID)
	if err != nil {
		if errors.Is(err, database.ErrModelNotFound) {
			return nil, apperr.ErrTaskNotFound(err)
		}
		return nil, fmt.Errorf("failed to get task: %w", err)
	}
	if !user.HasTask(task) {
		return nil, apperr.ErrTaskNotFound(errors.New("user does not own the task"))
	}

	now := clock.Now(ctx)
	s := domain.Step{
		ID:        domain.StepID(idgen.ULID(ctx)),
		UserID:    user.ID,
		TaskID:    in.TaskID,
		Name:      in.Name,
		CreatedAt: now,
		UpdatedAt: now,
	}

	if err := uc.DB.CreateStep(ctx, &s); err != nil {
		return nil, fmt.Errorf("failed to create step: %w", err)
	}
	return &StepOutput{Step: &s}, nil
}

type UpdateStepInput struct {
	ProjectID   domain.ProjectID
	TaskID      domain.TaskID
	ID          domain.StepID
	Name        *string
	CompletedAt *time.Time
}

func (uc *Step) UpdateStep(ctx context.Context, in *UpdateStepInput) (*StepOutput, error) {
	user := auth.MustUserFromContext(ctx)

	p, err := uc.DB.GetProjectByID(ctx, in.ProjectID)
	if err != nil {
		if errors.Is(err, database.ErrModelNotFound) {
			return nil, apperr.ErrProjectNotFound(err)
		}
		return nil, fmt.Errorf("failed to get project: %w", err)
	}
	if !user.HasProject(p) {
		return nil, apperr.ErrProjectNotFound(errors.New("user does not own the project"))
	}

	t, err := uc.DB.GetTaskByID(ctx, in.TaskID)
	if err != nil {
		if errors.Is(err, database.ErrModelNotFound) {
			return nil, apperr.ErrTaskNotFound(err)
		}
		return nil, fmt.Errorf("failed to get task: %w", err)
	}
	if !user.HasTask(t) {
		return nil, apperr.ErrTaskNotFound(errors.New("user does not own the task"))
	}
	if p.ID != t.ProjectID {
		return nil, apperr.ErrTaskNotFound(errors.New("task does not belong to the project"))
	}

	s, err := uc.DB.GetStepByID(ctx, in.ID)
	if err != nil {
		if errors.Is(err, database.ErrModelNotFound) {
			return nil, apperr.ErrStepNotFound(err)
		}
		return nil, fmt.Errorf("failed to get step: %w", err)
	}
	if !user.HasStep(s) {
		return nil, apperr.ErrStepNotFound(errors.New("user does not own the step"))
	}
	if t.ID != s.TaskID {
		return nil, apperr.ErrStepNotFound(errors.New("step does not belong to the task"))
	}

	if in.Name != nil {
		s.Name = *in.Name
	}
	if in.CompletedAt != nil {
		s.CompletedAt = pointers.Ref(clock.Now(ctx))
	}
	s.UpdatedAt = clock.Now(ctx)

	if err := uc.DB.UpdateStep(ctx, s); err != nil {
		return nil, fmt.Errorf("failed to update step: %w", err)
	}
	return &StepOutput{Step: s}, nil
}

type DeleteStepInput struct {
	ProjectID domain.ProjectID
	TaskID    domain.TaskID
	ID        domain.StepID
}

func (uc *Step) DeleteStep(ctx context.Context, in *DeleteStepInput) error {
	user := auth.MustUserFromContext(ctx)

	p, err := uc.DB.GetProjectByID(ctx, in.ProjectID)
	if err != nil {
		if errors.Is(err, database.ErrModelNotFound) {
			return apperr.ErrProjectNotFound(err)
		}
		return fmt.Errorf("failed to get project: %w", err)
	}
	if !user.HasProject(p) {
		return apperr.ErrProjectNotFound(errors.New("user does not own the project"))
	}

	t, err := uc.DB.GetTaskByID(ctx, in.TaskID)
	if err != nil {
		if errors.Is(err, database.ErrModelNotFound) {
			return apperr.ErrTaskNotFound(err)
		}
		return fmt.Errorf("failed to get task: %w", err)
	}
	if !user.HasTask(t) {
		return apperr.ErrTaskNotFound(errors.New("user does not own the task"))
	}
	if p.ID != t.ProjectID {
		return apperr.ErrTaskNotFound(errors.New("task does not belong to the project"))
	}

	s, err := uc.DB.GetStepByID(ctx, in.ID)
	if err != nil {
		if errors.Is(err, database.ErrModelNotFound) {
			return apperr.ErrStepNotFound(err)
		}
		return fmt.Errorf("failed to get step: %w", err)
	}
	if !user.HasStep(s) {
		return apperr.ErrStepNotFound(errors.New("user does not own the step"))
	}
	if t.ID != s.TaskID {
		return apperr.ErrStepNotFound(errors.New("step does not belong to the task"))
	}

	if err := uc.DB.DeleteStepByID(ctx, s.ID); err != nil {
		return fmt.Errorf("failed to delete step: %w", err)
	}
	return nil
}

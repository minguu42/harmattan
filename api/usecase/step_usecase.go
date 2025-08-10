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
			return nil, apperr.TaskNotFoundError(err)
		}
		return nil, fmt.Errorf("failed to get task: %w", err)
	}
	if !user.HasTask(task) {
		return nil, apperr.TaskNotFoundError(errors.New("user does not own the task"))
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
	TaskID      domain.TaskID
	ID          domain.StepID
	Name        opt.Option[string]
	CompletedAt opt.Option[*time.Time]
}

func (uc *Step) UpdateStep(ctx context.Context, in *UpdateStepInput) (*StepOutput, error) {
	user := auth.MustUserFromContext(ctx)

	t, err := uc.DB.GetTaskByID(ctx, in.TaskID)
	if err != nil {
		if errors.Is(err, database.ErrModelNotFound) {
			return nil, apperr.TaskNotFoundError(err)
		}
		return nil, fmt.Errorf("failed to get task: %w", err)
	}
	if !user.HasTask(t) {
		return nil, apperr.TaskNotFoundError(errors.New("user does not own the task"))
	}

	s, err := uc.DB.GetStepByID(ctx, in.ID)
	if err != nil {
		if errors.Is(err, database.ErrModelNotFound) {
			return nil, apperr.StepNotFoundError(err)
		}
		return nil, fmt.Errorf("failed to get step: %w", err)
	}
	if !user.HasStep(s) {
		return nil, apperr.StepNotFoundError(errors.New("user does not own the step"))
	}
	if t.ID != s.TaskID {
		return nil, apperr.StepNotFoundError(errors.New("step does not belong to the task"))
	}

	if in.Name.Valid {
		s.Name = in.Name.V
	}
	if in.CompletedAt.Valid {
		s.CompletedAt = in.CompletedAt.V
	}
	s.UpdatedAt = clock.Now(ctx)

	if err := uc.DB.UpdateStep(ctx, s); err != nil {
		return nil, fmt.Errorf("failed to update step: %w", err)
	}
	return &StepOutput{Step: s}, nil
}

type DeleteStepInput struct {
	TaskID domain.TaskID
	ID     domain.StepID
}

func (uc *Step) DeleteStep(ctx context.Context, in *DeleteStepInput) error {
	user := auth.MustUserFromContext(ctx)

	t, err := uc.DB.GetTaskByID(ctx, in.TaskID)
	if err != nil {
		if errors.Is(err, database.ErrModelNotFound) {
			return apperr.TaskNotFoundError(err)
		}
		return fmt.Errorf("failed to get task: %w", err)
	}
	if !user.HasTask(t) {
		return apperr.TaskNotFoundError(errors.New("user does not own the task"))
	}

	s, err := uc.DB.GetStepByID(ctx, in.ID)
	if err != nil {
		if errors.Is(err, database.ErrModelNotFound) {
			return apperr.StepNotFoundError(err)
		}
		return fmt.Errorf("failed to get step: %w", err)
	}
	if !user.HasStep(s) {
		return apperr.StepNotFoundError(errors.New("user does not own the step"))
	}
	if t.ID != s.TaskID {
		return apperr.StepNotFoundError(errors.New("step does not belong to the task"))
	}

	if err := uc.DB.DeleteStepByID(ctx, s.ID); err != nil {
		return fmt.Errorf("failed to delete step: %w", err)
	}
	return nil
}

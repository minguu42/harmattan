package usecase

import (
	"context"
	"errors"
	"time"

	"github.com/minguu42/harmattan/internal/auth"
	"github.com/minguu42/harmattan/internal/database"
	"github.com/minguu42/harmattan/internal/domain"
	"github.com/minguu42/harmattan/internal/lib/clock"
	"github.com/minguu42/harmattan/internal/lib/errtrace"
	"github.com/minguu42/harmattan/internal/lib/idgen"
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
			return nil, errtrace.Wrap(TaskNotFoundError())
		}
		return nil, errtrace.Wrap(err)
	}
	if !user.HasTask(task) {
		return nil, errtrace.Wrap(TaskNotFoundError())
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
		return nil, errtrace.Wrap(err)
	}
	return &StepOutput{Step: &s}, nil
}

type UpdateStepInput struct {
	ID          domain.StepID
	Name        Option[string]
	CompletedAt Option[*time.Time]
}

func (uc *Step) UpdateStep(ctx context.Context, in *UpdateStepInput) (*StepOutput, error) {
	user := auth.MustUserFromContext(ctx)

	s, err := uc.DB.GetStepByID(ctx, in.ID)
	if err != nil {
		if errors.Is(err, database.ErrModelNotFound) {
			return nil, errtrace.Wrap(StepNotFoundError())
		}
		return nil, errtrace.Wrap(err)
	}
	if !user.HasStep(s) {
		return nil, errtrace.Wrap(StepNotFoundError())
	}

	if in.Name.Valid {
		s.Name = in.Name.V
	}
	if in.CompletedAt.Valid {
		s.CompletedAt = in.CompletedAt.V
	}
	s.UpdatedAt = clock.Now(ctx)

	if err := uc.DB.UpdateStep(ctx, s); err != nil {
		return nil, errtrace.Wrap(err)
	}
	return &StepOutput{Step: s}, nil
}

type DeleteStepInput struct {
	ID domain.StepID
}

func (uc *Step) DeleteStep(ctx context.Context, in *DeleteStepInput) error {
	user := auth.MustUserFromContext(ctx)

	s, err := uc.DB.GetStepByID(ctx, in.ID)
	if err != nil {
		if errors.Is(err, database.ErrModelNotFound) {
			return errtrace.Wrap(StepNotFoundError())
		}
		return errtrace.Wrap(err)
	}
	if !user.HasStep(s) {
		return errtrace.Wrap(StepNotFoundError())
	}

	if err := uc.DB.DeleteStepByID(ctx, s.ID); err != nil {
		return errtrace.Wrap(err)
	}
	return nil
}

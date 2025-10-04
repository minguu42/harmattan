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

func (h *handler) CreateStep(ctx context.Context, req *openapi2.CreateStepReq, params openapi2.CreateStepParams) (*openapi2.Step, error) {
	var errs []error
	errs = append(errs, validateStepName(req.Name)...)
	if len(errs) > 0 {
		return nil, usecase2.DomainValidationError(errs)
	}

	out, err := h.step.CreateStep(ctx, &usecase2.CreateStepInput{
		TaskID: domain.TaskID(params.TaskID),
		Name:   req.Name,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to execute CreateStep usecase: %w", err)
	}
	return convertStep(out.Step), nil
}

func (h *handler) UpdateStep(ctx context.Context, req *openapi2.UpdateStepReq, params openapi2.UpdateStepParams) (*openapi2.Step, error) {
	var errs []error
	if name, ok := req.Name.Get(); ok {
		errs = append(errs, validateStepName(name)...)
	}
	if len(errs) > 0 {
		return nil, usecase2.DomainValidationError(errs)
	}

	out, err := h.step.UpdateStep(ctx, &usecase2.UpdateStepInput{
		ID:          domain.StepID(params.StepID),
		Name:        usecase2.Option[string]{V: req.Name.Value, Valid: req.Name.Set},
		CompletedAt: usecase2.Option[*time.Time]{V: pointers.Ternary(req.CompletedAt.Null, nil, &req.CompletedAt.Value), Valid: req.CompletedAt.Set},
	})
	if err != nil {
		return nil, fmt.Errorf("failed to execute UpdateStep usecase: %w", err)
	}
	return convertStep(out.Step), nil
}

func (h *handler) DeleteStep(ctx context.Context, params openapi2.DeleteStepParams) error {
	if err := h.step.DeleteStep(ctx, &usecase2.DeleteStepInput{ID: domain.StepID(params.StepID)}); err != nil {
		return fmt.Errorf("failed to execute DeleteStep usecase: %w", err)
	}
	return nil
}

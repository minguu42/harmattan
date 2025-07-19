package handler

import (
	"context"
	"fmt"

	"github.com/minguu42/harmattan/api/usecase"
	"github.com/minguu42/harmattan/internal/domain"
	"github.com/minguu42/harmattan/internal/openapi"
)

func convertStep(s *domain.Step) *openapi.Step {
	return &openapi.Step{
		ID:          string(s.ID),
		TaskID:      string(s.TaskID),
		Name:        s.Name,
		CompletedAt: convertDateTimePtr(s.CompletedAt),
		CreatedAt:   s.CreatedAt,
		UpdatedAt:   s.UpdatedAt,
	}
}

func (h *handler) CreateStep(ctx context.Context, req *openapi.CreateStepReq, params openapi.CreateStepParams) (*openapi.Step, error) {
	out, err := h.step.CreateStep(ctx, &usecase.CreateStepInput{
		TaskID: domain.TaskID(params.TaskID),
		Name:   req.Name,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to execute CreateStep usecase: %w", err)
	}
	return convertStep(out.Step), nil
}

func (h *handler) UpdateStep(ctx context.Context, req *openapi.UpdateStepReq, params openapi.UpdateStepParams) (*openapi.Step, error) {
	out, err := h.step.UpdateStep(ctx, &usecase.UpdateStepInput{
		ProjectID:   domain.ProjectID(params.ProjectID),
		TaskID:      domain.TaskID(params.TaskID),
		ID:          domain.StepID(params.StepID),
		Name:        convertOptString(req.Name),
		CompletedAt: convertOptDateTime(req.CompletedAt),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to execute UpdateStep usecase: %w", err)
	}
	return convertStep(out.Step), nil
}

func (h *handler) DeleteStep(ctx context.Context, params openapi.DeleteStepParams) error {
	err := h.step.DeleteStep(ctx, &usecase.DeleteStepInput{
		ProjectID: domain.ProjectID(params.ProjectID),
		TaskID:    domain.TaskID(params.TaskID),
		ID:        domain.StepID(params.StepID),
	})
	if err != nil {
		return fmt.Errorf("failed to execute DeleteStep usecase: %w", err)
	}
	return nil
}

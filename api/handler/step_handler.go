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
	input := &usecase.UpdateStepInput{
		ID: domain.StepID(params.StepID),
	}

	if req.Name.Set {
		input.Name = &req.Name.Value
	}
	if req.CompletedAt.Set {
		if req.CompletedAt.Value.IsZero() {
			empty := ""
			input.CompletedAt = &empty
		} else {
			nonEmpty := "completed"
			input.CompletedAt = &nonEmpty
		}
	}

	out, err := h.step.UpdateStep(ctx, input)
	if err != nil {
		return nil, err
	}

	step := &openapi.Step{
		ID:     string(out.Step.ID),
		TaskID: string(out.Step.TaskID),
		Name:   out.Step.Name,
	}

	if out.Step.CompletedAt != nil {
		step.CompletedAt = openapi.OptDateTime{
			Value: *out.Step.CompletedAt,
			Set:   true,
		}
	}

	return step, nil
}

func (h *handler) DeleteStep(ctx context.Context, params openapi.DeleteStepParams) error {
	return h.step.DeleteStep(ctx, &usecase.DeleteStepInput{
		ID: domain.StepID(params.StepID),
	})
}

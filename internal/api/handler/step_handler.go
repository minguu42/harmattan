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

func (h *Handler) CreateStep(ctx context.Context, req *openapi.CreateStepReq, params openapi.CreateStepParams) (*openapi.Step, error) {
	var errs []error
	errs = append(errs, validateStepName(req.Name)...)
	if len(errs) > 0 {
		return nil, errtrace.Wrap(apierror.DomainValidationError(errs))
	}

	out, err := h.Step.CreateStep(ctx, &usecase.CreateStepInput{
		TaskID: domain.TaskID(params.TaskID),
		Name:   req.Name,
	})
	if err != nil {
		return nil, errtrace.Wrap(err)
	}
	return convertStep(out.Step), nil
}

func (h *Handler) UpdateStep(ctx context.Context, req *openapi.UpdateStepReq, params openapi.UpdateStepParams) (*openapi.Step, error) {
	var errs []error
	if name, ok := req.Name.Get(); ok {
		errs = append(errs, validateStepName(name)...)
	}
	if len(errs) > 0 {
		return nil, errtrace.Wrap(apierror.DomainValidationError(errs))
	}

	out, err := h.Step.UpdateStep(ctx, &usecase.UpdateStepInput{
		ID:          domain.StepID(params.StepID),
		Name:        usecase.Option[string]{V: req.Name.Value, Valid: req.Name.Set},
		CompletedAt: usecase.Option[*time.Time]{V: ternary(req.CompletedAt.Null, nil, &req.CompletedAt.Value), Valid: req.CompletedAt.Set},
	})
	if err != nil {
		return nil, errtrace.Wrap(err)
	}
	return convertStep(out.Step), nil
}

func (h *Handler) DeleteStep(ctx context.Context, params openapi.DeleteStepParams) error {
	if err := h.Step.DeleteStep(ctx, &usecase.DeleteStepInput{ID: domain.StepID(params.StepID)}); err != nil {
		return errtrace.Wrap(err)
	}
	return nil
}

var ErrStepNameLength = errors.New("ステップ名は1文字以上100文字以下で指定できます")

func validateStepName(name string) []error {
	var errs []error
	if utf8.RuneCountInString(name) < 1 || 100 < utf8.RuneCountInString(name) {
		errs = append(errs, ErrStepNameLength)
	}
	return errs
}

func convertStep(s *domain.Step) *openapi.Step {
	return &openapi.Step{
		ID:          string(s.ID),
		TaskID:      string(s.TaskID),
		Name:        s.Name,
		CompletedAt: convertOptDateTime(s.CompletedAt),
		CreatedAt:   s.CreatedAt,
		UpdatedAt:   s.UpdatedAt,
	}
}

func convertSteps(steps domain.Steps) []openapi.Step {
	s := make([]openapi.Step, 0, len(steps))
	for _, step := range steps {
		s = append(s, *convertStep(&step))
	}
	return s
}

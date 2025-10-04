package handler

import (
	"context"
	"fmt"

	openapi2 "github.com/minguu42/harmattan/internal/api/handler/openapi"
	usecase2 "github.com/minguu42/harmattan/internal/api/usecase"
	"github.com/minguu42/harmattan/internal/domain"
)

func (h *handler) CreateTag(ctx context.Context, req *openapi2.CreateTagReq) (*openapi2.Tag, error) {
	var errs []error
	errs = append(errs, validateTagName(req.Name)...)
	if len(errs) > 0 {
		return nil, usecase2.DomainValidationError(errs)
	}

	out, err := h.tag.CreateTag(ctx, &usecase2.CreateTagInput{Name: req.Name})
	if err != nil {
		return nil, fmt.Errorf("failed to execute CreateTag usecase: %w", err)
	}
	return convertTag(out.Tag), nil
}

func (h *handler) ListTags(ctx context.Context, params openapi2.ListTagsParams) (*openapi2.ListTagsOK, error) {
	out, err := h.tag.ListTags(ctx, &usecase2.ListTagsInput{
		Limit:  params.Limit.Value,
		Offset: params.Offset.Value,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to execute ListTags usecase: %w", err)
	}
	return &openapi2.ListTagsOK{
		Tags:    convertTags(out.Tags),
		HasNext: out.HasNext,
	}, nil
}

func (h *handler) UpdateTag(ctx context.Context, req *openapi2.UpdateTagReq, params openapi2.UpdateTagParams) (*openapi2.Tag, error) {
	var errs []error
	if name, ok := req.Name.Get(); ok {
		errs = validateTagName(name)
	}
	if len(errs) > 0 {
		return nil, usecase2.DomainValidationError(errs)
	}

	out, err := h.tag.UpdateTag(ctx, &usecase2.UpdateTagInput{
		ID:   domain.TagID(params.TagID),
		Name: usecase2.Option[string]{V: req.Name.Value, Valid: req.Name.Set},
	})
	if err != nil {
		return nil, fmt.Errorf("failed to execute UpdateTag usecase: %w", err)
	}
	return convertTag(out.Tag), nil
}

func (h *handler) GetTag(ctx context.Context, params openapi2.GetTagParams) (*openapi2.Tag, error) {
	out, err := h.tag.GetTag(ctx, &usecase2.GetTagInput{
		ID: domain.TagID(params.TagID),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to execute GetTag usecase: %w", err)
	}
	return convertTag(out.Tag), nil
}

func (h *handler) DeleteTag(ctx context.Context, params openapi2.DeleteTagParams) error {
	if err := h.tag.DeleteTag(ctx, &usecase2.DeleteTagInput{ID: domain.TagID(params.TagID)}); err != nil {
		return fmt.Errorf("failed to execute DeleteTag usecase: %w", err)
	}
	return nil
}

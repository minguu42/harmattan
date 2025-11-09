package handler

import (
	"context"

	"github.com/minguu42/harmattan/internal/api/openapi"
	"github.com/minguu42/harmattan/internal/api/usecase"
	"github.com/minguu42/harmattan/internal/domain"
	"github.com/minguu42/harmattan/internal/lib/errtrace"
)

func (h *handler) CreateTag(ctx context.Context, req *openapi.CreateTagReq) (*openapi.Tag, error) {
	var errs []error
	errs = append(errs, validateTagName(req.Name)...)
	if len(errs) > 0 {
		return nil, usecase.DomainValidationError(errs)
	}

	out, err := h.tag.CreateTag(ctx, &usecase.CreateTagInput{Name: req.Name})
	if err != nil {
		return nil, errtrace.Wrap(err)
	}
	return convertTag(out.Tag), nil
}

func (h *handler) ListTags(ctx context.Context, params openapi.ListTagsParams) (*openapi.ListTagsOK, error) {
	out, err := h.tag.ListTags(ctx, &usecase.ListTagsInput{
		Limit:  params.Limit.Value,
		Offset: params.Offset.Value,
	})
	if err != nil {
		return nil, errtrace.Wrap(err)
	}
	return &openapi.ListTagsOK{
		Tags:    convertTags(out.Tags),
		HasNext: out.HasNext,
	}, nil
}

func (h *handler) UpdateTag(ctx context.Context, req *openapi.UpdateTagReq, params openapi.UpdateTagParams) (*openapi.Tag, error) {
	var errs []error
	if name, ok := req.Name.Get(); ok {
		errs = validateTagName(name)
	}
	if len(errs) > 0 {
		return nil, usecase.DomainValidationError(errs)
	}

	out, err := h.tag.UpdateTag(ctx, &usecase.UpdateTagInput{
		ID:   domain.TagID(params.TagID),
		Name: usecase.Option[string]{V: req.Name.Value, Valid: req.Name.Set},
	})
	if err != nil {
		return nil, errtrace.Wrap(err)
	}
	return convertTag(out.Tag), nil
}

func (h *handler) GetTag(ctx context.Context, params openapi.GetTagParams) (*openapi.Tag, error) {
	out, err := h.tag.GetTag(ctx, &usecase.GetTagInput{
		ID: domain.TagID(params.TagID),
	})
	if err != nil {
		return nil, errtrace.Wrap(err)
	}
	return convertTag(out.Tag), nil
}

func (h *handler) DeleteTag(ctx context.Context, params openapi.DeleteTagParams) error {
	if err := h.tag.DeleteTag(ctx, &usecase.DeleteTagInput{ID: domain.TagID(params.TagID)}); err != nil {
		return errtrace.Wrap(err)
	}
	return nil
}

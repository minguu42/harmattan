package handler

import (
	"context"
	"fmt"

	"github.com/minguu42/harmattan/api/usecase"
	"github.com/minguu42/harmattan/internal/domain"
	"github.com/minguu42/harmattan/internal/openapi"
)

func convertTag(t *domain.Tag) *openapi.Tag {
	return &openapi.Tag{
		ID:        string(t.ID),
		Name:      t.Name,
		CreatedAt: t.CreatedAt,
		UpdatedAt: t.UpdatedAt,
	}
}

func convertTags(tags domain.Tags) []openapi.Tag {
	ts := make([]openapi.Tag, 0, len(tags))
	for _, t := range tags {
		ts = append(ts, *convertTag(&t))
	}
	return ts
}

func (h *handler) CreateTag(ctx context.Context, req *openapi.CreateTagReq) (*openapi.Tag, error) {
	out, err := h.tag.CreateTag(ctx, &usecase.CreateTagInput{Name: req.Name})
	if err != nil {
		return nil, fmt.Errorf("failed to execute CreateTag usecase: %w", err)
	}
	return convertTag(out.Tag), nil
}

func (h *handler) ListTags(ctx context.Context, params openapi.ListTagsParams) (*openapi.Tags, error) {
	out, err := h.tag.ListTags(ctx, &usecase.ListTagsInput{
		Limit:  params.Limit.Value,
		Offset: params.Offset.Value,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to execute ListTags usecase: %w", err)
	}
	return &openapi.Tags{
		Tags:    convertTags(out.Tags),
		HasNext: out.HasNext,
	}, nil
}

func (h *handler) UpdateTag(ctx context.Context, req *openapi.UpdateTagReq, params openapi.UpdateTagParams) (*openapi.Tag, error) {
	out, err := h.tag.UpdateTag(ctx, &usecase.UpdateTagInput{
		ID:   domain.TagID(params.TagID),
		Name: convertOptString(req.Name),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to execute UpdateTag usecase: %w", err)
	}
	return convertTag(out.Tag), nil
}

func (h *handler) DeleteTag(ctx context.Context, params openapi.DeleteTagParams) error {
	if err := h.tag.DeleteTag(ctx, &usecase.DeleteTagInput{ID: domain.TagID(params.TagID)}); err != nil {
		return fmt.Errorf("failed to execute DeleteTag usecase: %w", err)
	}
	return nil
}

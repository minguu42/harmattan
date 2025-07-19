package usecase

import (
	"context"
	"errors"
	"fmt"

	"github.com/minguu42/harmattan/api/apperr"
	"github.com/minguu42/harmattan/internal/auth"
	"github.com/minguu42/harmattan/internal/database"
	"github.com/minguu42/harmattan/internal/domain"
	"github.com/minguu42/harmattan/lib/clock"
	"github.com/minguu42/harmattan/lib/idgen"
)

type Tag struct {
	DB *database.Client
}

type TagOutput struct {
	Tag *domain.Tag
}

type TagsOutput struct {
	Tags    domain.Tags
	HasNext bool
}

type CreateTagInput struct {
	Name string
}

func (uc *Tag) CreateTag(ctx context.Context, in *CreateTagInput) (*TagOutput, error) {
	user := auth.MustUserFromContext(ctx)

	now := clock.Now(ctx)
	t := domain.Tag{
		ID:        domain.TagID(idgen.ULID(ctx)),
		UserID:    user.ID,
		Name:      in.Name,
		CreatedAt: now,
		UpdatedAt: now,
	}
	if err := uc.DB.CreateTag(ctx, &t); err != nil {
		return nil, fmt.Errorf("failed to create tag: %w", err)
	}
	return &TagOutput{Tag: &t}, nil
}

type ListTagsInput struct {
	Limit  int
	Offset int
}

func (uc *Tag) ListTags(ctx context.Context, in *ListTagsInput) (*TagsOutput, error) {
	user := auth.MustUserFromContext(ctx)

	ts, err := uc.DB.ListTags(ctx, user.ID, in.Limit+1, in.Offset)
	if err != nil {
		return nil, fmt.Errorf("failed to list tags: %w", err)
	}

	hasNext := false
	if len(ts) == in.Limit+1 {
		ts = ts[:in.Limit]
		hasNext = true
	}
	return &TagsOutput{Tags: ts, HasNext: hasNext}, nil
}

type UpdateTagInput struct {
	ID   domain.TagID
	Name *string
}

func (uc *Tag) UpdateTag(ctx context.Context, in *UpdateTagInput) (*TagOutput, error) {
	user := auth.MustUserFromContext(ctx)

	t, err := uc.DB.GetTagByID(ctx, in.ID)
	if err != nil {
		if errors.Is(err, database.ErrModelNotFound) {
			return nil, apperr.TagNotFoundError(err)
		}
		return nil, fmt.Errorf("failed to get tag: %w", err)
	}
	if !user.HasTag(t) {
		return nil, apperr.TagNotFoundError(errors.New("user does not own the tag"))
	}

	if in.Name != nil {
		t.Name = *in.Name
	}
	t.UpdatedAt = clock.Now(ctx)
	if err := uc.DB.UpdateTag(ctx, t); err != nil {
		return nil, fmt.Errorf("failed to update tag: %w", err)
	}
	return &TagOutput{Tag: t}, nil
}

type DeleteTagInput struct {
	ID domain.TagID
}

func (uc *Tag) DeleteTag(ctx context.Context, in *DeleteTagInput) error {
	user := auth.MustUserFromContext(ctx)

	t, err := uc.DB.GetTagByID(ctx, in.ID)
	if err != nil {
		if errors.Is(err, database.ErrModelNotFound) {
			return apperr.TagNotFoundError(err)
		}
		return fmt.Errorf("failed to get tag: %w", err)
	}
	if !user.HasTag(t) {
		return apperr.TagNotFoundError(errors.New("user does not own the tag"))
	}

	if err := uc.DB.DeleteTagByID(ctx, t.ID); err != nil {
		return fmt.Errorf("failed to delete tag: %w", err)
	}
	return nil
}

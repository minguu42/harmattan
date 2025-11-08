package usecase

import (
	"context"

	"github.com/minguu42/harmattan/internal/auth"
	"github.com/minguu42/harmattan/internal/database"
	"github.com/minguu42/harmattan/internal/domain"
	"github.com/minguu42/harmattan/internal/lib/clock"
	"github.com/minguu42/harmattan/internal/lib/errors"
	"github.com/minguu42/harmattan/internal/lib/idgen"
)

type Tag struct {
	DB *database.Client
}

type TagOutput struct {
	Tag *domain.Tag
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
		return nil, errors.Wrap(err)
	}
	return &TagOutput{Tag: &t}, nil
}

type ListTagsInput struct {
	Limit  int
	Offset int
}

type ListTagsOutput struct {
	Tags    domain.Tags
	HasNext bool
}

func (uc *Tag) ListTags(ctx context.Context, in *ListTagsInput) (*ListTagsOutput, error) {
	user := auth.MustUserFromContext(ctx)

	ts, err := uc.DB.ListTags(ctx, user.ID, in.Limit+1, in.Offset)
	if err != nil {
		return nil, errors.Wrap(err)
	}

	hasNext := false
	if len(ts) == in.Limit+1 {
		ts = ts[:in.Limit]
		hasNext = true
	}
	return &ListTagsOutput{Tags: ts, HasNext: hasNext}, nil
}

type UpdateTagInput struct {
	ID   domain.TagID
	Name Option[string]
}

func (uc *Tag) UpdateTag(ctx context.Context, in *UpdateTagInput) (*TagOutput, error) {
	user := auth.MustUserFromContext(ctx)

	t, err := uc.DB.GetTagByID(ctx, in.ID)
	if err != nil {
		if errors.Is(err, database.ErrModelNotFound) {
			return nil, TagNotFoundError(err)
		}
		return nil, errors.Wrap(err)
	}
	if !user.HasTag(t) {
		return nil, TagAccessDeniedError()
	}

	if in.Name.Valid {
		t.Name = in.Name.V
	}
	t.UpdatedAt = clock.Now(ctx)
	if err := uc.DB.UpdateTag(ctx, t); err != nil {
		return nil, errors.Wrap(err)
	}
	return &TagOutput{Tag: t}, nil
}

type GetTagInput struct {
	ID domain.TagID
}

func (uc *Tag) GetTag(ctx context.Context, in *GetTagInput) (*TagOutput, error) {
	user := auth.MustUserFromContext(ctx)

	t, err := uc.DB.GetTagByID(ctx, in.ID)
	if err != nil {
		if errors.Is(err, database.ErrModelNotFound) {
			return nil, TagNotFoundError(err)
		}
		return nil, errors.Wrap(err)
	}
	if !user.HasTag(t) {
		return nil, TagAccessDeniedError()
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
			return TagNotFoundError(err)
		}
		return errors.Wrap(err)
	}
	if !user.HasTag(t) {
		return TagAccessDeniedError()
	}

	if err := uc.DB.DeleteTagByID(ctx, t.ID); err != nil {
		return errors.Wrap(err)
	}
	return nil
}

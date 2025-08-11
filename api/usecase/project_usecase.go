package usecase

import (
	"context"
	"errors"
	"fmt"

	"github.com/minguu42/harmattan/api/apperr"
	"github.com/minguu42/harmattan/internal/auth"
	"github.com/minguu42/harmattan/internal/database"
	"github.com/minguu42/harmattan/internal/domain"
	"github.com/minguu42/harmattan/internal/lib/clock"
	"github.com/minguu42/harmattan/internal/lib/idgen"
	"github.com/minguu42/harmattan/internal/lib/opt"
)

type Project struct {
	DB *database.Client
}

type ProjectOutput struct {
	Project *domain.Project
}

type CreateProjectInput struct {
	Name  string
	Color domain.ProjectColor
}

func (uc *Project) CreateProject(ctx context.Context, in *CreateProjectInput) (*ProjectOutput, error) {
	user := auth.MustUserFromContext(ctx)

	now := clock.Now(ctx)
	p := domain.Project{
		ID:        domain.ProjectID(idgen.ULID(ctx)),
		UserID:    user.ID,
		Name:      in.Name,
		Color:     in.Color,
		CreatedAt: now,
		UpdatedAt: now,
	}
	if err := uc.DB.CreateProject(ctx, &p); err != nil {
		return nil, fmt.Errorf("failed to create project: %w", err)
	}
	return &ProjectOutput{Project: &p}, nil
}

type ListProjectsInput struct {
	Limit  int
	Offset int
}

type ListProjectsOutput struct {
	Projects domain.Projects
	HasNext  bool
}

func (uc *Project) ListProjects(ctx context.Context, in *ListProjectsInput) (*ListProjectsOutput, error) {
	user := auth.MustUserFromContext(ctx)

	ps, err := uc.DB.ListProjects(ctx, user.ID, in.Limit+1, in.Offset)
	if err != nil {
		return nil, fmt.Errorf("failed to list projects: %w", err)
	}

	hasNext := false
	if len(ps) == in.Limit+1 {
		ps = ps[:in.Limit]
		hasNext = true
	}
	return &ListProjectsOutput{Projects: ps, HasNext: hasNext}, nil
}

type GetProjectInput struct {
	ID domain.ProjectID
}

func (uc *Project) GetProject(ctx context.Context, in *GetProjectInput) (*ProjectOutput, error) {
	user := auth.MustUserFromContext(ctx)

	p, err := uc.DB.GetProjectByID(ctx, in.ID)
	if err != nil {
		if errors.Is(err, database.ErrModelNotFound) {
			return nil, apperr.ProjectNotFoundError(err)
		}
		return nil, fmt.Errorf("failed to get project: %w", err)
	}
	if !user.HasProject(p) {
		return nil, apperr.ProjectAccessDeniedError()
	}

	return &ProjectOutput{Project: p}, nil
}

type UpdateProjectInput struct {
	ID         domain.ProjectID
	Name       opt.Option[string]
	Color      opt.Option[domain.ProjectColor]
	IsArchived opt.Option[bool]
}

func (uc *Project) UpdateProject(ctx context.Context, in *UpdateProjectInput) (*ProjectOutput, error) {
	user := auth.MustUserFromContext(ctx)

	p, err := uc.DB.GetProjectByID(ctx, in.ID)
	if err != nil {
		if errors.Is(err, database.ErrModelNotFound) {
			return nil, apperr.ProjectNotFoundError(err)
		}
		return nil, fmt.Errorf("failed to get project: %w", err)
	}
	if !user.HasProject(p) {
		return nil, apperr.ProjectAccessDeniedError()
	}

	if in.Name.Valid {
		p.Name = in.Name.V
	}
	if in.Color.Valid {
		p.Color = in.Color.V
	}
	if in.IsArchived.Valid {
		p.IsArchived = in.IsArchived.V
	}
	p.UpdatedAt = clock.Now(ctx)
	if err := uc.DB.UpdateProject(ctx, p); err != nil {
		return nil, fmt.Errorf("failed to update project: %w", err)
	}
	return &ProjectOutput{Project: p}, nil
}

type DeleteProjectInput struct {
	ID domain.ProjectID
}

func (uc *Project) DeleteProject(ctx context.Context, in *DeleteProjectInput) error {
	user := auth.MustUserFromContext(ctx)

	p, err := uc.DB.GetProjectByID(ctx, in.ID)
	if err != nil {
		if errors.Is(err, database.ErrModelNotFound) {
			return apperr.ProjectNotFoundError(err)
		}
		return fmt.Errorf("failed to get project: %w", err)
	}
	if !user.HasProject(p) {
		return apperr.ProjectAccessDeniedError()
	}

	if err := uc.DB.DeleteProjectByID(ctx, p.ID); err != nil {
		return fmt.Errorf("failed to delete project: %w", err)
	}
	return nil
}

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

type Project struct {
	DB *database.Client
}

type ProjectOutput struct {
	Project *domain.Project
}

type ProjectsOutput struct {
	Projects domain.Projects
	HasNext  bool
}

type CreateProjectInput struct {
	Name  string
	Color string
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

func (uc *Project) ListProjects(ctx context.Context, in *ListProjectsInput) (*ProjectsOutput, error) {
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
	return &ProjectsOutput{Projects: ps, HasNext: hasNext}, nil
}

type UpdateProjectInput struct {
	ID         domain.ProjectID
	Name       *string
	Color      *string
	IsArchived *bool
}

func (uc *Project) UpdateProject(ctx context.Context, in *UpdateProjectInput) (*ProjectOutput, error) {
	user := auth.MustUserFromContext(ctx)

	p, err := uc.DB.GetProjectByID(ctx, in.ID)
	if err != nil {
		if errors.Is(err, database.ErrModelNotFound) {
			return nil, apperr.ErrProjectNotFound(err)
		}
		return nil, fmt.Errorf("failed to get project: %w", err)
	}
	if !user.HasProject(p) {
		return nil, apperr.ErrProjectNotFound(errors.New("user does not own the project"))
	}

	if in.Name != nil {
		p.Name = *in.Name
	}
	if in.Color != nil {
		p.Color = *in.Color
	}
	if in.IsArchived != nil {
		p.IsArchived = *in.IsArchived
	}
	p.UpdatedAt = clock.Now(ctx)
	if err := uc.DB.UpdateProject(ctx, p); err != nil {
		return nil, fmt.Errorf("failed to update project: %w", err)
	}
	return &ProjectOutput{Project: p}, nil
}

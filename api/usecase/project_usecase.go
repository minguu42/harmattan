package usecase

import (
	"context"
	"fmt"

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
	user := auth.UserFromContext(ctx)

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
	user := auth.UserFromContext(ctx)

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

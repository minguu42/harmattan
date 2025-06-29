package handler

import (
	"context"
	"fmt"

	"github.com/minguu42/harmattan/api/usecase"
	"github.com/minguu42/harmattan/internal/domain"
	"github.com/minguu42/harmattan/internal/oapi"
)

func convertProject(project *domain.Project) *oapi.Project {
	return &oapi.Project{
		ID:         string(project.ID),
		Name:       project.Name,
		Color:      project.Color,
		IsArchived: project.IsArchived,
		CreatedAt:  project.CreatedAt,
		UpdatedAt:  project.UpdatedAt,
	}
}

func convertProjects(projects domain.Projects) []oapi.Project {
	ps := make([]oapi.Project, 0, len(projects))
	for _, p := range projects {
		ps = append(ps, *convertProject(&p))
	}
	return ps
}

func (h *handler) CreateProject(ctx context.Context, req *oapi.CreateProjectReq) (*oapi.Project, error) {
	out, err := h.project.CreateProject(ctx, &usecase.CreateProjectInput{
		Name:  req.Name,
		Color: req.Color,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to execute CreateProject usecase: %w", err)
	}
	return convertProject(out.Project), nil
}

func (h *handler) ListProjects(ctx context.Context, params oapi.ListProjectsParams) (*oapi.Projects, error) {
	out, err := h.project.ListProjects(ctx, &usecase.ListProjectsInput{
		Limit:  params.Limit.Value,
		Offset: params.Offset.Value,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to execute ListProjects usecase: %w", err)
	}
	return &oapi.Projects{
		Projects: convertProjects(out.Projects),
		HasNext:  out.HasNext,
	}, nil
}

func (h *handler) UpdateProject(ctx context.Context, req *oapi.UpdateProjectReq, params oapi.UpdateProjectParams) (*oapi.Project, error) {
	out, err := h.project.UpdateProject(ctx, &usecase.UpdateProjectInput{
		ID:         domain.ProjectID(params.ProjectID),
		Name:       convertOptString(req.Name),
		Color:      convertOptString(req.Color),
		IsArchived: convertOptBool(req.IsArchived),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to execute UpdateProject usecase: %w", err)
	}
	return convertProject(out.Project), nil
}

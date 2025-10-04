package handler

import (
	"context"
	"fmt"

	openapi2 "github.com/minguu42/harmattan/internal/api/handler/openapi"
	usecase2 "github.com/minguu42/harmattan/internal/api/usecase"
	"github.com/minguu42/harmattan/internal/domain"
)

func (h *handler) CreateProject(ctx context.Context, req *openapi2.CreateProjectReq) (*openapi2.Project, error) {
	var errs []error
	errs = append(errs, validateProjectName(req.Name)...)
	if len(errs) > 0 {
		return nil, usecase2.DomainValidationError(errs)
	}

	out, err := h.project.CreateProject(ctx, &usecase2.CreateProjectInput{
		Name:  req.Name,
		Color: domain.ProjectColor(req.Color),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to execute CreateProject usecase: %w", err)
	}
	return convertProject(out.Project), nil
}

func (h *handler) ListProjects(ctx context.Context, params openapi2.ListProjectsParams) (*openapi2.ListProjectsOK, error) {
	out, err := h.project.ListProjects(ctx, &usecase2.ListProjectsInput{
		Limit:  params.Limit.Value,
		Offset: params.Offset.Value,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to execute ListProjects usecase: %w", err)
	}
	return &openapi2.ListProjectsOK{
		Projects: convertProjects(out.Projects),
		HasNext:  out.HasNext,
	}, nil
}

func (h *handler) GetProject(ctx context.Context, params openapi2.GetProjectParams) (*openapi2.Project, error) {
	out, err := h.project.GetProject(ctx, &usecase2.GetProjectInput{
		ID: domain.ProjectID(params.ProjectID),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to execute GetProject usecase: %w", err)
	}
	return convertProject(out.Project), nil
}

func (h *handler) UpdateProject(ctx context.Context, req *openapi2.UpdateProjectReq, params openapi2.UpdateProjectParams) (*openapi2.Project, error) {
	var errs []error
	if name, ok := req.Name.Get(); ok {
		errs = append(errs, validateProjectName(name)...)
	}
	if len(errs) > 0 {
		return nil, usecase2.DomainValidationError(errs)
	}

	out, err := h.project.UpdateProject(ctx, &usecase2.UpdateProjectInput{
		ID:         domain.ProjectID(params.ProjectID),
		Name:       usecase2.Option[string]{V: req.Name.Value, Valid: req.Name.Set},
		Color:      usecase2.Option[domain.ProjectColor]{V: domain.ProjectColor(req.Color.Value), Valid: req.Color.Set},
		IsArchived: usecase2.Option[bool]{V: req.IsArchived.Value, Valid: req.IsArchived.Set},
	})
	if err != nil {
		return nil, fmt.Errorf("failed to execute UpdateProject usecase: %w", err)
	}
	return convertProject(out.Project), nil
}

func (h *handler) DeleteProject(ctx context.Context, params openapi2.DeleteProjectParams) error {
	if err := h.project.DeleteProject(ctx, &usecase2.DeleteProjectInput{ID: domain.ProjectID(params.ProjectID)}); err != nil {
		return fmt.Errorf("failed to execute DeleteProject usecase: %w", err)
	}
	return nil
}

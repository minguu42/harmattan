package handler

import (
	"context"
	"fmt"

	"github.com/minguu42/harmattan/api/apperr"
	"github.com/minguu42/harmattan/api/usecase"
	"github.com/minguu42/harmattan/internal/domain"
	"github.com/minguu42/harmattan/internal/openapi"
	"github.com/minguu42/harmattan/lib/opt"
)

func (h *handler) CreateProject(ctx context.Context, req *openapi.CreateProjectReq) (*openapi.Project, error) {
	var errs []error
	errs = append(errs, validateProjectName(req.Name)...)
	if len(errs) > 0 {
		return nil, apperr.DomainValidationError(errs)
	}

	out, err := h.project.CreateProject(ctx, &usecase.CreateProjectInput{
		Name:  req.Name,
		Color: domain.ProjectColor(req.Color),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to execute CreateProject usecase: %w", err)
	}
	return convertProject(out.Project), nil
}

func (h *handler) ListProjects(ctx context.Context, params openapi.ListProjectsParams) (*openapi.ListProjectsOK, error) {
	out, err := h.project.ListProjects(ctx, &usecase.ListProjectsInput{
		Limit:  params.Limit.Value,
		Offset: params.Offset.Value,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to execute ListProjects usecase: %w", err)
	}
	return &openapi.ListProjectsOK{
		Projects: convertProjects(out.Projects),
		HasNext:  out.HasNext,
	}, nil
}

func (h *handler) UpdateProject(ctx context.Context, req *openapi.UpdateProjectReq, params openapi.UpdateProjectParams) (*openapi.Project, error) {
	var errs []error
	if name, ok := req.Name.Get(); ok {
		errs = append(errs, validateProjectName(name)...)
	}
	if len(errs) > 0 {
		return nil, apperr.DomainValidationError(errs)
	}

	out, err := h.project.UpdateProject(ctx, &usecase.UpdateProjectInput{
		ID:         domain.ProjectID(params.ProjectID),
		Name:       opt.Option[string]{V: req.Name.Value, Valid: req.Name.Set},
		Color:      opt.Option[domain.ProjectColor]{V: domain.ProjectColor(req.Color.Value), Valid: req.Color.Set},
		IsArchived: opt.Option[bool]{V: req.IsArchived.Value, Valid: req.IsArchived.Set},
	})
	if err != nil {
		return nil, fmt.Errorf("failed to execute UpdateProject usecase: %w", err)
	}
	return convertProject(out.Project), nil
}

func (h *handler) GetProject(ctx context.Context, params openapi.GetProjectParams) (*openapi.Project, error) {
	out, err := h.project.GetProject(ctx, &usecase.GetProjectInput{
		ID: domain.ProjectID(params.ProjectID),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to execute GetProject usecase: %w", err)
	}
	return convertProject(out.Project), nil
}

func (h *handler) DeleteProject(ctx context.Context, params openapi.DeleteProjectParams) error {
	if err := h.project.DeleteProject(ctx, &usecase.DeleteProjectInput{ID: domain.ProjectID(params.ProjectID)}); err != nil {
		return fmt.Errorf("failed to execute DeleteProject usecase: %w", err)
	}
	return nil
}

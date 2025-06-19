package handler

import (
	"context"
	"fmt"

	"github.com/minguu42/harmattan/api/usecase"
	"github.com/minguu42/harmattan/internal/oapi"
)

func (h *handler) CreateProject(ctx context.Context, req *oapi.CreateProjectReq) (*oapi.Project, error) {
	out, err := h.project.CreateProject(ctx, &usecase.CreateProjectInput{
		Name:  req.Name,
		Color: req.Color,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to execute CreateProject usecase: %w", err)
	}
	return &oapi.Project{
		ID:         string(out.Project.ID),
		Name:       out.Project.Name,
		Color:      out.Project.Color,
		IsArchived: out.Project.IsArchived,
		CreatedAt:  out.Project.CreatedAt,
		UpdatedAt:  out.Project.UpdatedAt,
	}, nil
}

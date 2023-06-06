package app

import (
	"context"
	"fmt"
	"net/url"

	"github.com/minguu42/mtasks/app/logging"
	"github.com/minguu42/mtasks/app/ogen"
)

// PostProjects は POST /projects に対応するハンドラ
func (h *handler) PostProjects(ctx context.Context, req *ogen.PostProjectsReq, params ogen.PostProjectsParams) (ogen.PostProjectsRes, error) {
	u, err := h.repository.GetUserByAPIKey(ctx, params.XAPIKey)
	if err != nil {
		logging.Errorf("repository.GetUserByAPIKey failed: %v", err)
		return &ogen.PostProjectsUnauthorized{}, nil
	}

	p, err := h.repository.CreateProject(ctx, u.ID, req.Name)
	if err != nil {
		logging.Errorf("repository.CreateProject failed: %v", err)
		return &ogen.PostProjectsInternalServerError{}, nil
	}

	location, err := url.ParseRequestURI(fmt.Sprintf("http://localhost:8080/projects/%d", p.ID))
	if err != nil {
		logging.Errorf("url.ParseRequestURI failed: %v", err)
		return &ogen.PostProjectsInternalServerError{}, nil
	}
	return &ogen.ProjectHeaders{
		Location: *location,
		Response: newProjectResponse(p),
	}, nil
}

// GetProjects は GET /projects に対応するハンドラ
func (h *handler) GetProjects(_ context.Context, _ ogen.GetProjectsParams) (ogen.GetProjectsRes, error) {
	return &ogen.GetProjectsNotImplemented{}, nil
}

// PatchProject は PATCH /projects/{projectID} に対応するハンドラ
func (h *handler) PatchProject(_ context.Context, _ *ogen.PatchProjectReq, _ ogen.PatchProjectParams) (ogen.PatchProjectRes, error) {
	return &ogen.PatchProjectNotImplemented{}, nil
}

// DeleteProject は DELETE /projects/{projectID} に対応するハンドラ
func (h *handler) DeleteProject(_ context.Context, _ ogen.DeleteProjectParams) (ogen.DeleteProjectRes, error) {
	return &ogen.DeleteProjectNotImplemented{}, nil
}

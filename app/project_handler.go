package app

import (
	"context"
	"fmt"
	"net/url"
	"time"

	"github.com/minguu42/mtasks/app/logging"
	"github.com/minguu42/mtasks/app/ogen"
)

// PostProjects は POST /projects に対応するハンドラ
func (h *handler) PostProjects(ctx context.Context, req *ogen.PostProjectsReq, _ ogen.PostProjectsParams) (ogen.PostProjectsRes, error) {
	u, ok := ctx.Value(userKey{}).(*User)
	if !ok {
		logging.Errorf("ctx.Value(userKey{}).(*User) failed")
		return &ogen.PostProjectsInternalServerError{}, nil
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
func (h *handler) GetProjects(ctx context.Context, params ogen.GetProjectsParams) (ogen.GetProjectsRes, error) {
	u, ok := ctx.Value(userKey{}).(*User)
	if !ok {
		logging.Errorf("ctx.Value(userKey{}).(*User) failed")
		return &ogen.GetProjectsInternalServerError{}, nil
	}

	ps, err := h.repository.GetProjectsByUserID(ctx, u.ID, string(params.Sort.Or(ogen.GetProjectsSortMinusCreatedAt)), params.Limit.Or(10), params.Offset.Or(0))
	if err != nil {
		logging.Errorf("repository.GetProjectsByUserID failed: %v", err)
		return &ogen.GetProjectsInternalServerError{}, nil
	}

	return &ogen.Projects{Projects: newProjectsResponse(ps)}, nil
}

// PatchProject は PATCH /projects/{projectID} に対応するハンドラ
func (h *handler) PatchProject(ctx context.Context, req *ogen.PatchProjectReq, params ogen.PatchProjectParams) (ogen.PatchProjectRes, error) {
	u, ok := ctx.Value(userKey{}).(*User)
	if !ok {
		logging.Errorf("ctx.Value(userKey{}).(*User) failed")
		return &ogen.PatchProjectInternalServerError{}, nil
	}

	p, err := h.repository.GetProjectByID(ctx, params.ProjectID)
	if err != nil {
		logging.Errorf("repository.GetProjectByID failed: %v", err)
		return &ogen.PatchProjectInternalServerError{}, nil
	}

	if u.ID != p.UserID {
		logging.Errorf("u.ID != p.UserID")
		return &ogen.PatchProjectNotFound{}, nil
	}

	if !req.Name.IsSet() {
		logging.Errorf("value contains nothing")
		return &ogen.PatchProjectBadRequest{}, nil
	}
	p.Name = req.Name.Value
	p.UpdatedAt = time.Now()
	if err := h.repository.UpdateProject(ctx, params.ProjectID, p.Name, p.UpdatedAt); err != nil {
		logging.Errorf("repository.UpdateProject failed: %v", err)
		return &ogen.PatchProjectInternalServerError{}, nil
	}

	resp := newProjectResponse(p)
	return &resp, nil
}

// DeleteProject は DELETE /projects/{projectID} に対応するハンドラ
func (h *handler) DeleteProject(ctx context.Context, params ogen.DeleteProjectParams) (ogen.DeleteProjectRes, error) {
	u, ok := ctx.Value(userKey{}).(*User)
	if !ok {
		logging.Errorf("ctx.Value(userKey{}).(*User) failed")
		return &ogen.DeleteProjectInternalServerError{}, nil
	}

	p, err := h.repository.GetProjectByID(ctx, params.ProjectID)
	if err != nil {
		logging.Errorf("repository.GetProjectByID failed: %v", err)
		return &ogen.DeleteProjectInternalServerError{}, nil
	}

	if u.ID != p.UserID {
		logging.Errorf("u.ID != p.UserID")
		return &ogen.DeleteProjectNotFound{}, nil
	}

	if err := h.repository.DeleteProject(ctx, p.ID); err != nil {
		logging.Errorf("repository.DeleteProject failed: %v", err)
		return &ogen.DeleteProjectInternalServerError{}, nil
	}

	return &ogen.DeleteProjectNoContent{}, nil
}

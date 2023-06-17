package handler

import (
	"context"
	"time"

	"github.com/minguu42/mtasks/pkg/entity"

	"github.com/minguu42/mtasks/pkg/logging"
	"github.com/minguu42/mtasks/pkg/ogen"
)

// CreateProject は POST /projects に対応するハンドラ
func (h *Handler) CreateProject(ctx context.Context, req *ogen.CreateProjectReq) (*ogen.Project, error) {
	u, ok := ctx.Value(userKey{}).(*entity.User)
	if !ok {
		logging.Errorf("ctx.Value(userKey{}).(*User) failed")
		return nil, errUnauthorized
	}

	p, err := h.Repository.CreateProject(ctx, u.ID, req.Name)
	if err != nil {
		logging.Errorf("repository.CreateProject failed: %v", err)
		return nil, errInternalServerError
	}

	return newProjectResponse(p), nil
}

// ListProjects は GET /projects に対応するハンドラ
func (h *Handler) ListProjects(ctx context.Context, params ogen.ListProjectsParams) (*ogen.Projects, error) {
	u, ok := ctx.Value(userKey{}).(*entity.User)
	if !ok {
		logging.Errorf("ctx.Value(userKey{}).(*User) failed")
		return nil, errUnauthorized
	}

	ps, err := h.Repository.GetProjectsByUserID(ctx, u.ID, string(params.Sort.Or(ogen.ListProjectsSortMinusCreatedAt)), params.Limit.Or(10)+1, params.Offset.Or(0))
	if err != nil {
		logging.Errorf("repository.GetProjectsByUserID failed: %v", err)
		return nil, errInternalServerError
	}

	hasNext := false
	if len(ps) == params.Limit.Or(10)+1 {
		hasNext = true
		ps = ps[:params.Limit.Or(10)]
	}

	return &ogen.Projects{
		Projects: newProjectsResponse(ps),
		HasNext:  hasNext,
	}, nil
}

// UpdateProject は PATCH /projects/{projectID} に対応するハンドラ
func (h *Handler) UpdateProject(ctx context.Context, req *ogen.UpdateProjectReq, params ogen.UpdateProjectParams) (*ogen.Project, error) {
	u, ok := ctx.Value(userKey{}).(*entity.User)
	if !ok {
		logging.Errorf("ctx.Value(userKey{}).(*User) failed")
		return nil, errUnauthorized
	}

	p, err := h.Repository.GetProjectByID(ctx, params.ProjectID)
	if err != nil {
		logging.Errorf("repository.GetProjectByID failed: %v", err)
		return nil, errInternalServerError
	}

	if !u.HasProject(p) {
		logging.Errorf("user does not have the project")
		return nil, errProjectNotFound
	}

	if !req.Name.IsSet() {
		logging.Errorf("value contains nothing")
		return nil, errBadRequest
	}

	p.Name = req.Name.Value
	p.UpdatedAt = time.Now()
	if err := h.Repository.UpdateProject(ctx, params.ProjectID, p.Name, p.UpdatedAt); err != nil {
		logging.Errorf("repository.UpdateProject failed: %v", err)
		return nil, errInternalServerError
	}

	return newProjectResponse(p), nil
}

// DeleteProject は DELETE /projects/{projectID} に対応するハンドラ
func (h *Handler) DeleteProject(ctx context.Context, params ogen.DeleteProjectParams) error {
	u, ok := ctx.Value(userKey{}).(*entity.User)
	if !ok {
		logging.Errorf("ctx.Value(userKey{}).(*User) failed")
		return errUnauthorized
	}

	p, err := h.Repository.GetProjectByID(ctx, params.ProjectID)
	if err != nil {
		logging.Errorf("repository.GetProjectByID failed: %v", err)
		return errInternalServerError
	}

	if !u.HasProject(p) {
		logging.Errorf("user does not have the project")
		return errProjectNotFound
	}

	if err := h.Repository.DeleteProject(ctx, p.ID); err != nil {
		logging.Errorf("repository.DeleteProject failed: %v", err)
		return errInternalServerError
	}

	return nil
}

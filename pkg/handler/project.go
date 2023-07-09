package handler

import (
	"context"

	"github.com/go-faster/errors"
	"github.com/minguu42/mtasks/gen/ogen"
	"github.com/minguu42/mtasks/pkg/entity"
	"github.com/minguu42/mtasks/pkg/logging"
	"github.com/minguu42/mtasks/pkg/ttime"
	"gorm.io/gorm"
)

// CreateProject は POST /projects に対応するハンドラ
func (h *Handler) CreateProject(ctx context.Context, req *ogen.CreateProjectReq) (*ogen.Project, error) {
	u, ok := ctx.Value(userKey{}).(*entity.User)
	if !ok {
		return nil, errUnauthorized
	}

	p, err := h.Repository.CreateProject(ctx, u.ID, req.Name)
	if err != nil {
		logging.Errorf(ctx, "repository.CreateProject failed: %v", err)
		return nil, errInternalServerError
	}

	return newProjectResponse(p), nil
}

// ListProjects は GET /projects に対応するハンドラ
func (h *Handler) ListProjects(ctx context.Context, params ogen.ListProjectsParams) (*ogen.Projects, error) {
	u, ok := ctx.Value(userKey{}).(*entity.User)
	if !ok {
		return nil, errUnauthorized
	}

	ps, err := h.Repository.GetProjectsByUserID(ctx, u.ID, string(params.Sort.Or(ogen.ListProjectsSortMinusCreatedAt)), params.Limit.Or(10)+1, params.Offset.Or(0))
	if err != nil {
		logging.Errorf(ctx, "repository.GetProjectsByUserID failed: %v", err)
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
	if !req.Name.IsSet() {
		logging.Errorf(ctx, "value contains nothing")
		return nil, errBadRequest
	}

	u, ok := ctx.Value(userKey{}).(*entity.User)
	if !ok {
		return nil, errUnauthorized
	}

	p, err := h.Repository.GetProjectByID(ctx, params.ProjectID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errProjectNotFound
		}
		logging.Errorf(ctx, "repository.GetProjectByID failed: %v", err)
		return nil, errInternalServerError
	}

	if !u.HasProject(p) {
		logging.Errorf(ctx, "user does not have the project")
		return nil, errProjectNotFound
	}

	p.Name = req.Name.Value
	p.UpdatedAt = ttime.Now(ctx)
	if err := h.Repository.UpdateProject(ctx, params.ProjectID, p.Name, p.UpdatedAt); err != nil {
		logging.Errorf(ctx, "repository.UpdateProject failed: %v", err)
		return nil, errInternalServerError
	}

	return newProjectResponse(p), nil
}

// DeleteProject は DELETE /projects/{projectID} に対応するハンドラ
func (h *Handler) DeleteProject(ctx context.Context, params ogen.DeleteProjectParams) error {
	u, ok := ctx.Value(userKey{}).(*entity.User)
	if !ok {
		return errUnauthorized
	}

	p, err := h.Repository.GetProjectByID(ctx, params.ProjectID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errProjectNotFound
		}
		logging.Errorf(ctx, "repository.GetProjectByID failed: %v", err)
		return errInternalServerError
	}

	if !u.HasProject(p) {
		logging.Errorf(ctx, "user does not have the project")
		return errProjectNotFound
	}

	if err := h.Repository.DeleteProject(ctx, p.ID); err != nil {
		logging.Errorf(ctx, "repository.DeleteProject failed: %v", err)
		return errInternalServerError
	}

	return nil
}

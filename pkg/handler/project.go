package handler

import (
	"context"

	"github.com/go-faster/errors"
	"github.com/minguu42/opepe/gen/ogen"
	"github.com/minguu42/opepe/pkg/entity"
	"github.com/minguu42/opepe/pkg/logging"
	"github.com/minguu42/opepe/pkg/repository"
	"github.com/minguu42/opepe/pkg/ttime"
)

// CreateProject は POST /projects に対応するハンドラ
func (h *Handler) CreateProject(ctx context.Context, req *ogen.CreateProjectReq) (*ogen.Project, error) {
	u, ok := ctx.Value(userKey{}).(*entity.User)
	if !ok {
		return nil, errUnauthorized
	}

	p, err := h.Repository.CreateProject(ctx, u.ID, req.Name, req.Color)
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

	limit := params.Limit.Or(defaultLimit)
	ps, err := h.Repository.GetProjectsByUserID(ctx, u.ID, string(params.Sort.Or(ogen.ListProjectsSortMinusCreatedAt)), limit+1, params.Offset.Or(defaultOffset))
	if err != nil {
		logging.Errorf(ctx, "repository.GetProjectsByUserID failed: %v", err)
		return nil, errInternalServerError
	}

	hasNext := false
	if len(ps) == limit+1 {
		hasNext = true
		ps = ps[:limit]
	}

	return &ogen.Projects{
		Projects: newProjectsResponse(ps),
		HasNext:  hasNext,
	}, nil
}

func updateProject(ctx context.Context, p *entity.Project, req *ogen.UpdateProjectReq) *entity.Project {
	p.Name = req.Name.Or(p.Name)
	p.Color = req.Color.Or(p.Color)
	p.IsArchived = req.IsArchived.Or(p.IsArchived)
	p.UpdatedAt = ttime.Now(ctx)
	return p
}

// UpdateProject は PATCH /projects/{projectID} に対応するハンドラ
func (h *Handler) UpdateProject(ctx context.Context, req *ogen.UpdateProjectReq, params ogen.UpdateProjectParams) (*ogen.Project, error) {
	u, ok := ctx.Value(userKey{}).(*entity.User)
	if !ok {
		return nil, errUnauthorized
	}

	p, err := h.Repository.GetProjectByID(ctx, params.ProjectID)
	if err != nil {
		if errors.Is(err, repository.ErrRecordNotFound) {
			return nil, errProjectNotFound
		}
		logging.Errorf(ctx, "repository.GetProjectByID failed: %v", err)
		return nil, errInternalServerError
	}

	if !u.HasProject(p) {
		logging.Errorf(ctx, "user does not have the project")
		return nil, errProjectNotFound
	}

	newProject := updateProject(ctx, p, req)
	if err := h.Repository.UpdateProject(ctx, newProject); err != nil {
		logging.Errorf(ctx, "repository.UpdateProject failed: %v", err)
		return nil, errInternalServerError
	}

	return newProjectResponse(newProject), nil
}

// DeleteProject は DELETE /projects/{projectID} に対応するハンドラ
func (h *Handler) DeleteProject(ctx context.Context, params ogen.DeleteProjectParams) error {
	u, ok := ctx.Value(userKey{}).(*entity.User)
	if !ok {
		return errUnauthorized
	}

	p, err := h.Repository.GetProjectByID(ctx, params.ProjectID)
	if err != nil {
		if errors.Is(err, repository.ErrRecordNotFound) {
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

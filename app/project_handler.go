package app

import (
	"context"
	"time"

	"github.com/minguu42/mtasks/app/logging"
	"github.com/minguu42/mtasks/app/ogen"
)

// CreateProject は POST /projects に対応するハンドラ
func (h *handler) CreateProject(ctx context.Context, req *ogen.CreateProjectReq) (ogen.CreateProjectRes, error) {
	u, ok := ctx.Value(userKey{}).(*User)
	if !ok {
		logging.Errorf("ctx.Value(userKey{}).(*User) failed")
		return &ogen.CreateProjectInternalServerError{
			Message: messageInternalServerError,
			Debug:   "ctx.Value(userKey{}).(*User) failed",
		}, nil
	}

	p, err := h.repository.CreateProject(ctx, u.ID, req.Name)
	if err != nil {
		logging.Errorf("repository.CreateProject failed: %v", err)
		return &ogen.CreateProjectInternalServerError{
			Message: messageInternalServerError,
			Debug:   err.Error(),
		}, nil
	}

	return newProjectResponse(p), nil
}

// ListProjects は GET /projects に対応するハンドラ
func (h *handler) ListProjects(ctx context.Context, params ogen.ListProjectsParams) (ogen.ListProjectsRes, error) {
	u, ok := ctx.Value(userKey{}).(*User)
	if !ok {
		logging.Errorf("ctx.Value(userKey{}).(*User) failed")
		return &ogen.ListProjectsInternalServerError{
			Message: messageInternalServerError,
			Debug:   "ctx.Value(userKey{}).(*User) failed",
		}, nil
	}

	ps, err := h.repository.GetProjectsByUserID(ctx, u.ID, string(params.Sort.Or(ogen.ListProjectsSortMinusCreatedAt)), params.Limit.Or(10)+1, params.Offset.Or(0))
	if err != nil {
		logging.Errorf("repository.GetProjectsByUserID failed: %v", err)
		return &ogen.ListProjectsInternalServerError{
			Message: messageInternalServerError,
			Debug:   err.Error(),
		}, nil
	}

	hasNext := false
	if len(ps) == params.Limit.Or(10)+1 {
		hasNext = true
		ps = ps[:params.Limit.Or(10)]
	}
	logging.Debugf("ps: %v, cap: %d, len: %d\n", ps, cap(ps), len(ps))

	return &ogen.Projects{
		Projects: newProjectsResponse(ps),
		HasNext:  hasNext,
	}, nil
}

// UpdateProject は PATCH /projects/{projectID} に対応するハンドラ
func (h *handler) UpdateProject(ctx context.Context, req *ogen.UpdateProjectReq, params ogen.UpdateProjectParams) (ogen.UpdateProjectRes, error) {
	u, ok := ctx.Value(userKey{}).(*User)
	if !ok {
		logging.Errorf("ctx.Value(userKey{}).(*User) failed")
		return &ogen.UpdateProjectInternalServerError{
			Message: messageInternalServerError,
			Debug:   "ctx.Value(userKey{}).(*User) failed",
		}, nil
	}

	p, err := h.repository.GetProjectByID(ctx, params.ProjectID)
	if err != nil {
		logging.Errorf("repository.GetProjectByID failed: %v", err)
		return &ogen.UpdateProjectInternalServerError{
			Message: messageInternalServerError,
			Debug:   err.Error(),
		}, nil
	}

	if u.ID != p.UserID {
		logging.Errorf("u.ID != p.UserID")
		return &ogen.UpdateProjectNotFound{
			Message: messageNotFound,
			Debug:   err.Error(),
		}, nil
	}

	if !req.Name.IsSet() {
		logging.Errorf("value contains nothing")
		return &ogen.UpdateProjectBadRequest{
			Message: messageBadRequest,
			Debug:   err.Error(),
		}, nil
	}

	p.Name = req.Name.Value
	p.UpdatedAt = time.Now()
	if err := h.repository.UpdateProject(ctx, params.ProjectID, p.Name, p.UpdatedAt); err != nil {
		logging.Errorf("repository.UpdateProject failed: %v", err)
		return &ogen.UpdateProjectInternalServerError{
			Message: messageInternalServerError,
			Debug:   err.Error(),
		}, nil
	}

	return newProjectResponse(p), nil
}

// DeleteProject は DELETE /projects/{projectID} に対応するハンドラ
func (h *handler) DeleteProject(ctx context.Context, params ogen.DeleteProjectParams) (ogen.DeleteProjectRes, error) {
	u, ok := ctx.Value(userKey{}).(*User)
	if !ok {
		logging.Errorf("ctx.Value(userKey{}).(*User) failed")
		return &ogen.DeleteProjectInternalServerError{
			Message: messageInternalServerError,
			Debug:   "ctx.Value(userKey{}).(*User) failed",
		}, nil
	}

	p, err := h.repository.GetProjectByID(ctx, params.ProjectID)
	if err != nil {
		logging.Errorf("repository.GetProjectByID failed: %v", err)
		return &ogen.DeleteProjectInternalServerError{
			Message: messageInternalServerError,
			Debug:   err.Error(),
		}, nil
	}

	if u.ID != p.UserID {
		logging.Errorf("u.ID != p.UserID")
		return &ogen.DeleteProjectNotFound{
			Message: messageNotFound,
			Debug:   err.Error(),
		}, nil
	}

	if err := h.repository.DeleteProject(ctx, p.ID); err != nil {
		logging.Errorf("repository.DeleteProject failed: %v", err)
		return &ogen.DeleteProjectInternalServerError{
			Message: messageInternalServerError,
			Debug:   err.Error(),
		}, nil
	}

	return &ogen.DeleteProjectNoContent{}, nil
}

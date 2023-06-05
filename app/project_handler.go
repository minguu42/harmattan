package app

import (
	"context"

	"github.com/minguu42/mtasks/app/ogen"
)

// PostProjects は POST /projects に対応するハンドラ
func (h *handler) PostProjects(_ context.Context, _ *ogen.PostProjectsReq, _ ogen.PostProjectsParams) (ogen.PostProjectsRes, error) {
	return &ogen.PostProjectsNotImplemented{}, nil
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

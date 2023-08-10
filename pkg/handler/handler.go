// Package handler はハンドラに関するパッケージ
package handler

import (
	"context"
	"errors"
	"net/http"

	"github.com/minguu42/opepe/gen/ogen"
	"github.com/minguu42/opepe/pkg/domain/idgen"
	"github.com/minguu42/opepe/pkg/domain/model"
	"github.com/minguu42/opepe/pkg/domain/repository"
)

var (
	defaultLimit  = 10
	defaultOffset = 0
)

// ハンドラが返すエラー一覧
var (
	errUnauthorized        = errors.New("ユーザの認証に失敗しました。もしくはユーザが認証されていません。")
	errTaskNotFound        = errors.New("指定されたプロジェクトが見つかりません。")
	errProjectNotFound     = errors.New("指定されたタスクが見つかりません。")
	errInternalServerError = errors.New("不明なエラーが発生しました。")
)

// Handler は ogen.Handler を満たすハンドラ
type Handler struct {
	Repository  repository.Repository
	IDGenerator idgen.IDGenerator
}

// NewError はハンドラから渡されるエラーから ogen.ErrorStatusCode を生成する
func (h *Handler) NewError(_ context.Context, err error) *ogen.ErrorStatusCode {
	var code int
	switch {
	case errors.Is(err, errUnauthorized):
		code = http.StatusUnauthorized
	case errors.Is(err, errProjectNotFound):
		code = http.StatusNotFound
	case errors.Is(err, errTaskNotFound):
		code = http.StatusNotFound
	case errors.Is(err, errInternalServerError):
		code = http.StatusInternalServerError
	default:
		code = http.StatusInternalServerError
		err = errInternalServerError
	}
	return &ogen.ErrorStatusCode{
		StatusCode: code,
		Response: ogen.Error{
			Code:    code,
			Message: err.Error(),
		},
	}
}

// newProjectResponse は entity.Project から ogen.Project を生成する
func newProjectResponse(p *model.Project) *ogen.Project {
	return &ogen.Project{
		ID:         p.ID,
		Name:       p.Name,
		Color:      p.Color,
		IsArchived: p.IsArchived,
		CreatedAt:  p.CreatedAt,
		UpdatedAt:  p.UpdatedAt,
	}
}

// newProjectsResponse は entity.Project のスライスから ogen.Project のスライスを生成する
func newProjectsResponse(ps []model.Project) []ogen.Project {
	projects := make([]ogen.Project, 0, len(ps))
	for _, p := range ps {
		projects = append(projects, *newProjectResponse(&p))
	}
	return projects
}

// newTaskResponse は entity.Task から ogen.Task を生成する
func newTaskResponse(t *model.Task) *ogen.Task {
	dueOn := ogen.OptDate{}
	if t.DueOn != nil {
		dueOn = ogen.NewOptDate(*t.DueOn)
	}
	completedAt := ogen.OptDateTime{}
	if t.CompletedAt != nil {
		completedAt = ogen.NewOptDateTime(*t.CompletedAt)
	}
	return &ogen.Task{
		ID:          t.ID,
		ProjectID:   t.ProjectID,
		Title:       t.Title,
		Content:     t.Content,
		Priority:    t.Priority,
		DueOn:       dueOn,
		CompletedAt: completedAt,
		CreatedAt:   t.CreatedAt,
		UpdatedAt:   t.UpdatedAt,
	}
}

// newTasksResponse は entity.Task のスライスから ogen.Task のスライスを生成する
func newTasksResponse(ts []model.Task) []ogen.Task {
	tasks := make([]ogen.Task, 0, len(ts))
	for _, t := range ts {
		tasks = append(tasks, *newTaskResponse(&t))
	}
	return tasks
}

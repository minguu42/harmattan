// Package handler はハンドラ、ミドルウェアに関するパッケージ
package handler

import (
	"context"
	"errors"

	"github.com/minguu42/mtasks/app"
	"github.com/minguu42/mtasks/app/ogen"
	"github.com/minguu42/mtasks/app/repository"
)

var (
	errBadRequest          = errors.New("there is an input error")
	errUnauthorized        = errors.New("user is not authenticated")
	errTaskNotFound        = errors.New("the specified task is not found")
	errProjectNotFound     = errors.New("the specified project is not found")
	errInternalServerError = errors.New("some error occurred on the server")
	errNotImplemented      = errors.New("this API operation is not yet implemented")
	errServerUnavailable   = errors.New("server is temporarily unavailable")
)

// Handler -
type Handler struct {
	Repository repository.Repository
}

// NewError -
func (h *Handler) NewError(_ context.Context, err error) *ogen.ErrorStatusCode {
	var (
		statusCode int
		message    string
	)
	switch err {
	case errBadRequest:
		statusCode = 400
		message = "入力に誤りがあります。入力をご確認ください。"
	case errUnauthorized:
		statusCode = 401
		message = "ユーザが認証されていません。ユーザの認証後にもう一度お試しください。"
	case errProjectNotFound:
		statusCode = 404
		message = "指定されたプロジェクトが見つかりません。もう一度ご確認ください。"
	case errTaskNotFound:
		statusCode = 404
		message = "指定されたタスクが見つかりません。もう一度ご確認ください。"
	case errInternalServerError:
		statusCode = 500
		message = "不明なエラーが発生しました。もう一度お試しください。"
	case errNotImplemented:
		statusCode = 501
		message = "この機能はもうすぐ使用できます。お楽しみに♪"
	case errServerUnavailable:
		statusCode = 503
		message = "サーバが一時的に利用できない状態です。時間を空けてから、もう一度お試しください。"
	default:
		statusCode = 500
		message = "不明なエラーが発生しました。もう一度お試しください。"
		err = errInternalServerError
	}

	return &ogen.ErrorStatusCode{
		StatusCode: statusCode,
		Response: ogen.Error{
			Message: message,
			Debug:   err.Error(),
		},
	}
}

// newProjectResponse は Project から ogen.Project を生成する
func newProjectResponse(p *app.Project) *ogen.Project {
	return &ogen.Project{
		ID:        p.ID,
		Name:      p.Name,
		CreatedAt: p.CreatedAt,
		UpdatedAt: p.UpdatedAt,
	}
}

// newProjectsResponse は Project のスライスから ogen.Project のスライスを生成する
func newProjectsResponse(ps []*app.Project) []ogen.Project {
	projects := make([]ogen.Project, 0, len(ps))
	for _, p := range ps {
		projects = append(projects, *newProjectResponse(p))
	}
	return projects
}

// newTaskResponse は Task から ogen.Task を生成する
func newTaskResponse(t *app.Task) *ogen.Task {
	completedAt := ogen.OptDateTime{}
	if t.CompletedAt != nil {
		completedAt = ogen.NewOptDateTime(*t.CompletedAt)
	}
	return &ogen.Task{
		ID:          t.ID,
		ProjectID:   t.ProjectID,
		Title:       t.Title,
		CompletedAt: completedAt,
		CreatedAt:   t.CreatedAt,
		UpdatedAt:   t.UpdatedAt,
	}
}

// newTasksResponse は Task のスライスから ogen.Task のスライスを生成する
func newTasksResponse(ts []*app.Task) []ogen.Task {
	tasks := make([]ogen.Task, 0, len(ts))
	for _, t := range ts {
		tasks = append(tasks, *newTaskResponse(t))
	}
	return tasks
}

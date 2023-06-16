// Package app はアプリケーションの中心に関するパッケージ
package app

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/minguu42/mtasks/app/env"
	"github.com/minguu42/mtasks/app/ogen"
)

type handler struct {
	repository repository
}

type repository interface {
	GetUserByAPIKey(ctx context.Context, apiKey string) (*User, error)

	CreateProject(ctx context.Context, userID int64, name string) (*Project, error)
	GetProjectsByUserID(ctx context.Context, userID int64, sort string, limit, offset int) ([]*Project, error)
	GetProjectByID(ctx context.Context, id int64) (*Project, error)
	UpdateProject(ctx context.Context, id int64, name string, updatedAt time.Time) error
	DeleteProject(ctx context.Context, id int64) error

	CreateTask(ctx context.Context, projectID int64, title string) (*Task, error)
	GetTasksByProjectID(ctx context.Context, projectID int64, sort string, limit, offset int) ([]*Task, error)
	GetTaskByID(ctx context.Context, id int64) (*Task, error)
	UpdateTask(ctx context.Context, id int64, completedAt *time.Time, updatedAt time.Time) error
	DeleteTask(ctx context.Context, id int64) error
}

// NewServer はサーバを初期化する
func NewServer(api *env.API, repository repository) (*http.Server, error) {
	s, err := ogen.NewServer(
		&handler{repository: repository},
		&securityHandler{repository: repository},
	)
	if err != nil {
		return nil, fmt.Errorf("ogen.NewServer failed: %w", err)
	}

	return &http.Server{
		Addr:              fmt.Sprintf("%s:%d", api.Host, api.Port),
		Handler:           logMiddleware(s),
		ReadTimeout:       10 * time.Second,
		ReadHeaderTimeout: 10 * time.Second,
		MaxHeaderBytes:    1 << 20,
	}, nil
}

var (
	errBadRequest          = errors.New("there is an input error")
	errUnauthorized        = errors.New("user is not authenticated")
	errTaskNotFound        = errors.New("the specified task is not found")
	errProjectNotFound     = errors.New("the specified project is not found")
	errInternalServerError = errors.New("some error occurred on the server")
	errNotImplemented      = errors.New("this API operation is not yet implemented")
	errServerUnavailable   = errors.New("server is temporarily unavailable")
)

func (h *handler) NewError(_ context.Context, err error) *ogen.ErrorStatusCode {
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

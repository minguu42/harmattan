// Package app はアプリケーションの中心に関するパッケージ
package app

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/minguu42/mtasks/app/env"
	"github.com/minguu42/mtasks/app/ogen"
)

type handler struct {
	repository repository
}

type securityHandler struct{}

// HandleIsAuthorized -
// TODO: これを利用して authMiddleware を記述すべき
func (s *securityHandler) HandleIsAuthorized(ctx context.Context, _ string, _ ogen.IsAuthorized) (context.Context, error) {
	return ctx, nil
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
	s, err := ogen.NewServer(&handler{
		repository: repository,
	}, &securityHandler{})
	if err != nil {
		return nil, fmt.Errorf("ogen.NewServer failed: %w", err)
	}

	return &http.Server{
		Addr:              fmt.Sprintf("%s:%d", api.Host, api.Port),
		Handler:           logMiddleware(&authMiddleware{next: s, repository: repository}),
		ReadTimeout:       10 * time.Second,
		ReadHeaderTimeout: 10 * time.Second,
		MaxHeaderBytes:    1 << 20,
	}, nil
}

const (
	messageBadRequest          = "入力に誤りがあります。今一度入力をご確認ください。"
	messageUnauthorized        = "ユーザが認証されていません。ユーザの認証後にもう一度お試しください。"
	messageNotFound            = "指定されたリソースが存在しません。"
	messageInternalServerError = "不明なエラーが発生しました。もう一度お試しください。"
	messageNotImplemented      = "この機能はもうすぐ使用できます。お楽しみに♪"
	messageServerUnavailable   = "サーバが一時的に利用できない状態です。時間を空けてから、もう一度お試しください。"
)

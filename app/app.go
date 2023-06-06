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

type repository interface {
	GetUserByAPIKey(ctx context.Context, apiKey string) (*User, error)

	CreateProject(ctx context.Context, userID int64, name string) (*Project, error)
	GetProjectsByUserID(ctx context.Context, userID int64) ([]*Project, error)
	GetProjectByID(ctx context.Context, id int64) (*Project, error)
	UpdateProject(ctx context.Context, id int64, name string, updatedAt time.Time) error
	DeleteProject(ctx context.Context, id int64) error

	CreateTask(ctx context.Context, projectID int64, title string) (*Task, error)
	GetTasksByProjectID(ctx context.Context, projectID int64) ([]*Task, error)
	GetTaskByID(ctx context.Context, id int64) (*Task, error)
	UpdateTask(ctx context.Context, id int64, completedAt *time.Time, updatedAt time.Time) error
	DeleteTask(ctx context.Context, id int64) error
}

// NewServer はサーバを初期化する
func NewServer(api *env.API, repository repository) (*http.Server, error) {
	s, err := ogen.NewServer(&handler{
		repository: repository,
	})
	if err != nil {
		return nil, fmt.Errorf("ogen.NewServer failed: %w", err)
	}

	return &http.Server{
		Addr:              fmt.Sprintf("%s:%d", api.Host, api.Port),
		Handler:           s,
		ReadTimeout:       10 * time.Second,
		ReadHeaderTimeout: 10 * time.Second,
		MaxHeaderBytes:    1 << 20,
	}, nil
}

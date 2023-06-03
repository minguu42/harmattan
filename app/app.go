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
	GetUserByToken(ctx context.Context, token string) (*User, error)

	CreateTask(ctx context.Context, userID int64, title string) (*Task, error)
	GetTasksByUserID(ctx context.Context, userID int64) ([]*Task, error)
	GetTaskByID(ctx context.Context, id int64) (*Task, error)
	UpdateTask(ctx context.Context, id int64, completedAt *time.Time) error
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

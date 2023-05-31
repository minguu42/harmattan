// Package app はモデル、ハンドラ関数、データベース関数を定義する
package app

import (
	"context"
	"database/sql"
	"fmt"
	"net/http"
	"time"

	"github.com/minguu42/mtasks/pkg/env"
	"github.com/minguu42/mtasks/pkg/ogen"
)

type Handler struct {
	Repository Repository
}

type Repository interface {
	getUserByToken(ctx context.Context, token string) (*user, error)

	createTask(ctx context.Context, userID int64, title string) (*task, error)
	getTasksByUserID(ctx context.Context, userID int64) ([]*task, error)
	getTaskByID(ctx context.Context, id int64) (*task, error)
	updateTask(ctx context.Context, id int64, completedAt *time.Time) error
	deleteTask(ctx context.Context, id int64) error
}

type Database struct {
	*sql.DB
}

// NewServer はサーバを初期化する
func NewServer(api *env.API) (*http.Server, error) {
	s, err := ogen.NewServer(&Handler{
		Repository: &Database{db},
	})
	if err != nil {
		return nil, fmt.Errorf("ogen.NewServer failed: %w", err)
	}

	return &http.Server{
		Addr:              fmt.Sprintf(":%d", api.Port),
		Handler:           s,
		ReadTimeout:       10 * time.Second,
		ReadHeaderTimeout: 10 * time.Second,
		MaxHeaderBytes:    1 << 20,
	}, nil
}

// Package repository はデータベースへの操作を抽象化する
package repository

import (
	"context"
	"time"

	"github.com/minguu42/mtasks/app"
)

type Repository interface {
	GetUserByAPIKey(ctx context.Context, apiKey string) (*app.User, error)

	CreateProject(ctx context.Context, userID int64, name string) (*app.Project, error)
	GetProjectsByUserID(ctx context.Context, userID int64, sort string, limit, offset int) ([]*app.Project, error)
	GetProjectByID(ctx context.Context, id int64) (*app.Project, error)
	UpdateProject(ctx context.Context, id int64, name string, updatedAt time.Time) error
	DeleteProject(ctx context.Context, id int64) error

	CreateTask(ctx context.Context, projectID int64, title string) (*app.Task, error)
	GetTasksByProjectID(ctx context.Context, projectID int64, sort string, limit, offset int) ([]*app.Task, error)
	GetTaskByID(ctx context.Context, id int64) (*app.Task, error)
	UpdateTask(ctx context.Context, id int64, completedAt *time.Time, updatedAt time.Time) error
	DeleteTask(ctx context.Context, id int64) error
}

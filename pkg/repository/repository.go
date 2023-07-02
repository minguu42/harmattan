// Package repository はデータベースへの操作を抽象化する
package repository

//go:generate mockgen -source=$GOFILE -destination=mock/$GOFILE -package=mock

import (
	"context"
	"time"

	"github.com/minguu42/mtasks/pkg/entity"
)

type Repository interface {
	GetUserByAPIKey(ctx context.Context, apiKey string) (*entity.User, error)

	CreateProject(ctx context.Context, userID int64, name string) (*entity.Project, error)
	GetProjectsByUserID(ctx context.Context, userID int64, sort string, limit, offset int) ([]*entity.Project, error)
	GetProjectByID(ctx context.Context, id int64) (*entity.Project, error)
	UpdateProject(ctx context.Context, id int64, name string, updatedAt time.Time) error
	DeleteProject(ctx context.Context, id int64) error

	CreateTask(ctx context.Context, projectID int64, title string) (*entity.Task, error)
	GetTasksByProjectID(ctx context.Context, projectID int64, sort string, limit, offset int) ([]*entity.Task, error)
	GetTaskByID(ctx context.Context, id int64) (*entity.Task, error)
	UpdateTask(ctx context.Context, id int64, completedAt *time.Time, updatedAt time.Time) error
	DeleteTask(ctx context.Context, id int64) error
}

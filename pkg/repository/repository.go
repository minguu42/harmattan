// Package repository はデータベースへの操作を抽象化する
package repository

//go:generate mockgen -source=$GOFILE -destination=../../gen/mock/$GOFILE -package=mock

import (
	"context"
	"time"

	"github.com/minguu42/mtasks/pkg/entity"
)

type Repository interface {
	GetUserByAPIKey(ctx context.Context, apiKey string) (*entity.User, error)

	CreateProject(ctx context.Context, userID string, name string, color string) (*entity.Project, error)
	GetProjectsByUserID(ctx context.Context, userID string, sort string, limit, offset int) ([]*entity.Project, error)
	GetProjectByID(ctx context.Context, id string) (*entity.Project, error)
	UpdateProject(ctx context.Context, p *entity.Project) error
	DeleteProject(ctx context.Context, id string) error

	CreateTask(ctx context.Context, projectID string, title string) (*entity.Task, error)
	GetTasksByProjectID(ctx context.Context, projectID string, sort string, limit, offset int) ([]*entity.Task, error)
	GetTaskByID(ctx context.Context, id string) (*entity.Task, error)
	UpdateTask(ctx context.Context, id string, completedAt *time.Time, updatedAt time.Time) error
	DeleteTask(ctx context.Context, id string) error
}

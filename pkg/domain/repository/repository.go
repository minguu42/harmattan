// Package repository はデータストアにデータを保存する
package repository

//go:generate mockgen -source=$GOFILE -destination=../../../gen/mock/$GOFILE -package=mock

import (
	"context"
	"errors"

	"github.com/minguu42/opepe/pkg/domain/model"
)

var ErrModelNotFound = errors.New("model not found in data store")

type Repository interface {
	GetUserByAPIKey(ctx context.Context, apiKey string) (*model.User, error)

	CreateProject(ctx context.Context, p *model.Project) error
	GetProjectsByUserID(ctx context.Context, userID string, limit, offset int) ([]model.Project, error)
	GetProjectByID(ctx context.Context, id string) (*model.Project, error)
	UpdateProject(ctx context.Context, p *model.Project) error
	DeleteProject(ctx context.Context, id string) error

	CreateTask(ctx context.Context, t *model.Task) error
	GetTasksByProjectID(ctx context.Context, projectID string, limit, offset int) ([]model.Task, error)
	GetTaskByID(ctx context.Context, id string) (*model.Task, error)
	UpdateTask(ctx context.Context, t *model.Task) error
	DeleteTask(ctx context.Context, id string) error
}

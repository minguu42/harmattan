package database

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/minguu42/opepe/gen/sqlc"
	"github.com/minguu42/opepe/pkg/domain/model"
	"github.com/minguu42/opepe/pkg/domain/repository"
)

func (db *DB) CreateProject(ctx context.Context, p *model.Project) error {
	if err := sqlc.New(db.sqlDB).CreateProject(ctx, sqlc.CreateProjectParams{
		ID:         p.ID,
		UserID:     p.UserID,
		Name:       p.Name,
		Color:      p.Color,
		IsArchived: p.IsArchived,
		CreatedAt:  p.CreatedAt,
		UpdatedAt:  p.UpdatedAt,
	}); err != nil {
		return fmt.Errorf("q.CreateProject failed: %w", err)
	}
	return nil
}

func (db *DB) GetProjectsByUserID(ctx context.Context, userID string, limit, offset int) ([]model.Project, error) {
	projects, err := sqlc.New(db.sqlDB).GetProjectsByUserID(ctx, sqlc.GetProjectsByUserIDParams{
		UserID: userID,
		Limit:  int32(limit),
		Offset: int32(offset),
	})
	if err != nil {
		return nil, fmt.Errorf("q.GetProjectsByUserID failed: %w", err)
	}

	ps := make([]model.Project, 0, len(projects))
	for _, p := range projects {
		ps = append(ps, model.Project{
			ID:         p.ID,
			UserID:     p.UserID,
			Name:       p.Name,
			Color:      p.Color,
			IsArchived: p.IsArchived,
			CreatedAt:  p.CreatedAt,
			UpdatedAt:  p.UpdatedAt,
		})
	}
	return ps, nil
}

func (db *DB) GetProjectByID(ctx context.Context, id string) (*model.Project, error) {
	p, err := sqlc.New(db.sqlDB).GetProjectByID(ctx, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, repository.ErrModelNotFound
		}
		return nil, fmt.Errorf("q.GetProjectByID failed: %w", err)
	}
	return &model.Project{
		ID:         p.ID,
		UserID:     p.UserID,
		Name:       p.Name,
		Color:      p.Color,
		IsArchived: p.IsArchived,
		CreatedAt:  p.CreatedAt,
		UpdatedAt:  p.UpdatedAt,
	}, nil
}

func (db *DB) UpdateProject(ctx context.Context, p *model.Project) error {
	if err := sqlc.New(db.sqlDB).UpdateProject(ctx, sqlc.UpdateProjectParams{
		Name:  p.Name,
		Color: p.Color,
		ID:    p.ID,
	}); err != nil {
		return fmt.Errorf("q.UpdateProject failed: %w", err)
	}
	return nil
}

func (db *DB) DeleteProject(ctx context.Context, id string) error {
	if err := sqlc.New(db.sqlDB).DeleteProject(ctx, id); err != nil {
		return fmt.Errorf("q.DeleteProject failed: %w", err)
	}
	return nil
}

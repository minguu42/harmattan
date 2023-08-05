package database

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/minguu42/opepe/gen/sqlc"
	"github.com/minguu42/opepe/pkg/entity"
	"github.com/minguu42/opepe/pkg/repository"
	"github.com/minguu42/opepe/pkg/ttime"
)

func (db *DB) SaveProject(ctx context.Context, p *entity.Project) error {
	now := ttime.Now(ctx)
	if p.ID != "" {
		if err := sqlc.New(db._db).UpdateProject(ctx, sqlc.UpdateProjectParams{
			Name:  p.Name,
			Color: p.Color,
			ID:    p.ID,
		}); err != nil {
			return fmt.Errorf("q.UpdateProject failed: %w", err)
		}
		return nil
	}

	if err := sqlc.New(db._db).CreateProject(ctx, sqlc.CreateProjectParams{
		ID:         db.idGenerator.Generate(),
		UserID:     p.UserID,
		Name:       p.Name,
		Color:      p.Color,
		IsArchived: false,
		CreatedAt:  now,
		UpdatedAt:  now,
	}); err != nil {
		return fmt.Errorf("q.CreateProject failed: %w", err)
	}
	return nil
}

func (db *DB) GetProjectsByUserID(ctx context.Context, userID string, limit, offset int) ([]*entity.Project, error) {
	projects, err := sqlc.New(db._db).GetProjectsByUserID(ctx, sqlc.GetProjectsByUserIDParams{
		UserID: userID,
		Limit:  int32(limit),
		Offset: int32(offset),
	})
	if err != nil {
		return nil, fmt.Errorf("q.GetProjectsByUserID failed: %w", err)
	}

	ps := make([]*entity.Project, 0, len(projects))
	for _, p := range projects {
		ps = append(ps, &entity.Project{
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

func (db *DB) GetProjectByID(ctx context.Context, id string) (*entity.Project, error) {
	p, err := sqlc.New(db._db).GetProjectByID(ctx, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, repository.ErrRecordNotFound
		}
		return nil, fmt.Errorf("q.GetProjectByID failed: %w", err)
	}
	return &entity.Project{
		ID:         p.ID,
		UserID:     p.UserID,
		Name:       p.Name,
		Color:      p.Color,
		IsArchived: p.IsArchived,
		CreatedAt:  p.CreatedAt,
		UpdatedAt:  p.UpdatedAt,
	}, nil
}

func (db *DB) DeleteProject(ctx context.Context, id string) error {
	if err := sqlc.New(db._db).DeleteProject(ctx, id); err != nil {
		return fmt.Errorf("q.DeleteProject failed: %w", err)
	}
	return nil
}

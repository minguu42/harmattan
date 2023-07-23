package database

import (
	"context"
	"errors"
	"fmt"

	"github.com/minguu42/opepe/pkg/entity"
	"github.com/minguu42/opepe/pkg/repository"
	"github.com/minguu42/opepe/pkg/ttime"
	"gorm.io/gorm"
)

func (db *DB) CreateProject(ctx context.Context, userID string, name string, color string) (*entity.Project, error) {
	now := ttime.Now(ctx)
	p := entity.Project{
		ID:         db.idGenerator.Generate(),
		UserID:     userID,
		Name:       name,
		Color:      color,
		IsArchived: false,
		CreatedAt:  now,
		UpdatedAt:  now,
	}

	if err := db.conn(ctx).Create(&p).Error; err != nil {
		return nil, fmt.Errorf("gormDB.Create failed: %w", err)
	}
	return &p, nil
}

func (db *DB) GetProjectByID(ctx context.Context, id string) (*entity.Project, error) {
	var p entity.Project
	if err := db.conn(ctx).First(&p, "id = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, repository.ErrRecordNotFound
		}
		return nil, fmt.Errorf("gormDB.First failed: %w", err)
	}
	return &p, nil
}

func (db *DB) GetProjectsByUserID(ctx context.Context, userID string, sort string, limit, offset int) ([]*entity.Project, error) {
	ps := make([]*entity.Project, 0, limit)
	if err := db.conn(ctx).Where("user_id = ?", userID).
		Order(generateOrderByClause(sort)).Limit(limit).Offset(offset).Find(&ps).Error; err != nil {
		return nil, fmt.Errorf("gormDB.Find failed: %w", err)
	}
	return ps, nil
}

func (db *DB) UpdateProject(ctx context.Context, p *entity.Project) error {
	if err := db.conn(ctx).Model(&entity.Project{ID: p.ID}).Updates(entity.Project{
		Name:       p.Name,
		Color:      p.Color,
		IsArchived: p.IsArchived,
		UpdatedAt:  p.UpdatedAt,
	}).Error; err != nil {
		return fmt.Errorf("gormDB.Updates failed: %w", err)
	}
	return nil
}

func (db *DB) DeleteProject(ctx context.Context, id string) error {
	if err := db.conn(ctx).Delete(&entity.Project{ID: id}).Error; err != nil {
		return fmt.Errorf("gormDB.Delete failed: %w", err)
	}
	return nil
}

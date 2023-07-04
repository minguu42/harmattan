package database

import (
	"context"
	"fmt"
	"time"

	"github.com/minguu42/mtasks/pkg/entity"
)

func (db *DB) CreateProject(ctx context.Context, userID int64, name string) (*entity.Project, error) {
	now := time.Now()
	p := entity.Project{
		UserID:    userID,
		Name:      name,
		CreatedAt: now,
		UpdatedAt: now,
	}

	if err := db.conn(ctx).Create(&p).Error; err != nil {
		return nil, fmt.Errorf("gormDB.Create failed: %w", err)
	}
	return &p, nil
}

func (db *DB) GetProjectByID(ctx context.Context, id int64) (*entity.Project, error) {
	var p entity.Project
	if err := db.conn(ctx).First(&p, id).Error; err != nil {
		return nil, fmt.Errorf("gormDB.First failed: %w", err)
	}
	return &p, nil
}

func (db *DB) GetProjectsByUserID(ctx context.Context, userID int64, sort string, limit, offset int) ([]*entity.Project, error) {
	ps := make([]*entity.Project, 0, limit)
	if err := db.conn(ctx).Where("user_id = ?", userID).
		Order(generateOrderByClause(sort)).Limit(limit).Offset(offset).Find(&ps).Error; err != nil {
		return nil, fmt.Errorf("gormDB.Find failed: %w", err)
	}
	return ps, nil
}

func (db *DB) UpdateProject(ctx context.Context, id int64, name string, updatedAt time.Time) error {
	p := entity.Project{ID: id}
	if err := db.conn(ctx).Model(&p).Updates(entity.Project{Name: name, UpdatedAt: updatedAt}).Error; err != nil {
		return fmt.Errorf("gormDB.Updates failed: %w", err)
	}
	return nil
}

func (db *DB) DeleteProject(ctx context.Context, id int64) error {
	if err := db.conn(ctx).Delete(&entity.Project{}, id).Error; err != nil {
		return fmt.Errorf("gormDB.Delete failed: %w", err)
	}
	return nil
}

package database

import (
	"context"
	"errors"
	"fmt"

	"github.com/minguu42/opepe/pkg/entity"
	"github.com/minguu42/opepe/pkg/repository"
	"gorm.io/gorm"
)

func (db *DB) GetUserByAPIKey(ctx context.Context, apiKey string) (*entity.User, error) {
	var u entity.User
	if err := db.conn(ctx).Where("api_key = ?", apiKey).First(&u).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, repository.ErrRecordNotFound
		}
		return nil, fmt.Errorf("gormDB.Find failed: %w", err)
	}
	return &u, nil
}

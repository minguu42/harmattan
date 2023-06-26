package database

import (
	"context"
	"fmt"

	"github.com/minguu42/mtasks/pkg/entity"
)

func (db *DB) GetUserByAPIKey(ctx context.Context, apiKey string) (*entity.User, error) {
	var u entity.User
	if err := db.gormDB.WithContext(ctx).Where("api_key = ?", apiKey).First(&u).Error; err != nil {
		return nil, fmt.Errorf("gormDB.Find failed: %w", err)
	}
	return &u, nil
}

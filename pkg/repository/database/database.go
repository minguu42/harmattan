// Package database はデータベースに関するパッケージ
package database

import (
	"fmt"
	"strings"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// DB は repository.Repository インタフェースを実装するデータベース
type DB struct {
	gormDB *gorm.DB
}

// Close は新しいクエリの実行を辞め、データベースとの接続を閉じる
func (db *DB) Close() error {
	sqlDB, err := db.gormDB.DB()
	if err != nil {
		return fmt.Errorf("gormDB.DB failed: %w", err)
	}

	return sqlDB.Close()
}

// Open はデータベースとの接続を確立する
func Open(dsn string) (*DB, error) {
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("gorm.Open failed: %w", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("db.DB failed: %w", err)
	}
	sqlDB.SetConnMaxLifetime(3 * time.Minute)
	sqlDB.SetMaxOpenConns(10)
	sqlDB.SetMaxIdleConns(10)

	return &DB{gormDB: db}, nil
}

// generateOrderByClause は sort クエリから ORDER BY 句の値を生成する
// 例: 'createdAt' -> 'created_at ASC'、'-createdAt' -> 'created_at DESC'
func generateOrderByClause(sort string) string {
	m := map[string]string{
		"createdAt": "created_at",
		"updatedAt": "updated_at",
	}
	if strings.HasPrefix(sort, "-") {
		return fmt.Sprintf("%s DESC", m[strings.TrimPrefix(sort, "-")])
	}
	return fmt.Sprintf("%s ASC", m[sort])
}

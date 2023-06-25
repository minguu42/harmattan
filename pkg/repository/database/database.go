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
	*gorm.DB
}

// Close は新しいクエリの実行を辞め、データベースとの接続を閉じる
func (db *DB) Close() error {
	sqlDB, err := db.DB.DB()
	if err != nil {
		return fmt.Errorf("db.DB failed: %w", err)
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

	return &DB{db}, nil
}

// generateOrderByClause は sort クエリから ORDER BY 句の値を生成する
// 例: createdAt -> createdAt ASC、-createdAt -> createdAt DESC
func generateOrderByClause(sort string) string {
	if strings.HasPrefix(sort, "-") {
		return fmt.Sprintf("%s DESC", strings.TrimPrefix(sort, "-"))
	}
	return fmt.Sprintf("%s ASC", sort)
}

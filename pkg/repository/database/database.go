// Package database はデータベースに関するパッケージ
package database

import (
	"context"
	"fmt"
	"strings"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// DB は repository.Repository インタフェースを実装するデータベース
type DB struct {
	_db *gorm.DB
	_tx *gorm.DB
}

// conn はデータベースへのコネクションを返す
// NOTE: DB 構造体の _db, _tx フィールドは直接使用せず、このメソッドの戻り値を使用する
func (db *DB) conn(ctx context.Context) *gorm.DB {
	if db._tx != nil {
		return db._tx.WithContext(ctx)
	}
	return db._db.WithContext(ctx)
}

// Begin はトランザクションを開始する
func (db *DB) Begin() error {
	tx := db._db.Begin()
	if err := tx.Error; err != nil {
		return err
	}
	db._tx = tx
	return nil
}

// Rollback はトランザクションを終了し、トランザクション中の変更を元に戻す
func (db *DB) Rollback() {
	db._tx.Rollback()
	db._tx = nil
}

// Close は新しいクエリの実行を辞め、データベースとの接続を閉じる
func (db *DB) Close() error {
	sqlDB, err := db._db.DB()
	if err != nil {
		return fmt.Errorf("_db.DB failed: %w", err)
	}

	return sqlDB.Close()
}

// DSN はデータベースとの接続に使用する Data Source Name を生成する
func DSN(user, password, host string, port int, database string) string {
	return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=true",
		user,
		password,
		host,
		port,
		database,
	)
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

	return &DB{_db: db}, nil
}

// generateOrderByClause は sort クエリから ORDER BY 句の値を生成する
// 例: 'createdAt' -> 'created_at ASC'、'-createdAt' -> 'created_at DESC'
func generateOrderByClause(sort string) string {
	m := map[string]string{
		"name":      "name",
		"createdAt": "created_at",
		"updatedAt": "updated_at",
	}
	if strings.HasPrefix(sort, "-") {
		return fmt.Sprintf("%s DESC", m[strings.TrimPrefix(sort, "-")])
	}
	return fmt.Sprintf("%s ASC", m[sort])
}

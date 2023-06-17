// Package database はデータベースに関するパッケージ
package database

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/minguu42/mtasks/pkg/logging"
)

// DB は repository インタフェースを満たすデータベース
type DB struct {
	*sql.DB
}

// Open はデータベースとの接続を確立する
func Open(ctx context.Context, dsn string) (*DB, error) {
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, fmt.Errorf("sql.Open failed: %w", err)
	}

	maxFailureTimes := 2
	for i := 0; i <= maxFailureTimes; i++ {
		if err := db.PingContext(ctx); err == nil {
			break
		} else if i == maxFailureTimes {
			return nil, fmt.Errorf("db.PingContext failed: %w", err)
		}
		logging.Infof("db.PingContext failed. try again after 15 seconds")
		time.Sleep(15 * time.Second)
	}

	db.SetConnMaxLifetime(3 * time.Minute)
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(10)
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

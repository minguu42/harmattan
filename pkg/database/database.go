// Package database はデータベースに関するパッケージ
package database

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/minguu42/mtasks/pkg/logging"
)

// Open はデータベースとの接続が確立する
func Open(ctx context.Context, dsn string) (*sql.DB, error) {
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
	return db, nil
}

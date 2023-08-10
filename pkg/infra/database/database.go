// Package database はデータベースに関するパッケージ
package database

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/minguu42/opepe/pkg/logging"
)

// DB は repository.Repository インタフェースをじっそうs
type DB struct {
	sqlDB *sql.DB
}

// Close は新しいクエリの実行を停止し、データベースとの接続を切る
func (db *DB) Close() error {
	return db.sqlDB.Close()
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
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, fmt.Errorf("sql.Open failed: %w", err)
	}
	db.SetConnMaxLifetime(3 * time.Minute)
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(10)

	maxFailureTimes := 2
	for i := 0; i <= maxFailureTimes; i++ {
		if err := db.Ping(); err != nil && i != maxFailureTimes {
			logging.Infof(context.Background(), "db.Ping failed. try again after 15 seconds")
			time.Sleep(15 * time.Second)
			continue
		} else if err != nil && i == maxFailureTimes {
			return nil, fmt.Errorf("db.Ping failed: %w", err)
		}
		break
	}

	return &DB{sqlDB: db}, nil
}

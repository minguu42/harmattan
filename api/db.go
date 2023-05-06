package api

import (
	"database/sql"
	"fmt"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

// OpenDB はデータベースとの接続を確立した sql.DB を返す
func OpenDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, fmt.Errorf("sql.Open failed: %w", err)
	}

	maxFailureTimes := 2
	for {
		if err := db.Ping(); err == nil {
			break
		} else if maxFailureTimes <= 0 {
			return nil, fmt.Errorf("db.Ping failed: %w", err)
		}

		Infof("db.Ping failed. try again after 15 seconds")
		time.Sleep(15 * time.Second)
		maxFailureTimes--
	}

	db.SetConnMaxLifetime(3 * time.Minute)
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(10)
	return db, nil
}

// DSN はDSNを生成する
func DSN(user, password, host string, port int, dbname string) string {
	return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=true",
		user,
		password,
		host,
		port,
		dbname,
	)
}

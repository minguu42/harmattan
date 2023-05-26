package api

import (
	"database/sql"
	"fmt"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/minguu42/mtasks/pkg/logging"
)

var db *sql.DB

// OpenDB はデータベースとの接続を確立する
func OpenDB(dsn string) error {
	var err error
	db, err = sql.Open("mysql", dsn)
	if err != nil {
		return fmt.Errorf("sql.Open failed: %w", err)
	}

	maxFailureTimes := 2
	for {
		if err := db.Ping(); err == nil {
			break
		} else if maxFailureTimes <= 0 {
			return fmt.Errorf("db.Ping failed: %w", err)
		}

		logging.Infof("db.Ping failed. try again after 15 seconds")
		time.Sleep(15 * time.Second)
		maxFailureTimes--
	}

	db.SetConnMaxLifetime(3 * time.Minute)
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(10)
	return nil
}

// CloseDB はデータベースとの接続を終了する
func CloseDB() {
	_ = db.Close()
}

// DSN はデータベースとの接続に使用する Data Source Name を生成する
func DSN(user, password, host string, port int, dbname string) string {
	return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=true",
		user,
		password,
		host,
		port,
		dbname,
	)
}

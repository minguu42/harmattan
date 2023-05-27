package env

import (
	"errors"
	"fmt"
	"os"
	"strconv"
)

var appEnv *Env

// Load は環境変数を読み込む
// 環境変数を取得したい場合は Get 関数を使用する
func Load() error {
	apiPort, err := strconv.Atoi(os.Getenv("API_PORT"))
	if err != nil {
		return fmt.Errorf("strconv.Atoi failed: %w", err)
	}
	api := API{
		Host: os.Getenv("API_HOST"),
		Port: apiPort,
	}
	if api.Host != "" {
		return errors.New("API_HOST is required")
	}

	mysqlPort, err := strconv.Atoi(os.Getenv("MYSQL_PORT"))
	if err != nil {
		return fmt.Errorf("strconv.Atoi failed: %w", err)
	}
	mysql := MySQL{
		Host:     os.Getenv("MYSQL_HOST"),
		Port:     mysqlPort,
		Database: os.Getenv("MYSQL_DATABASE"),
		User:     os.Getenv("MYSQL_USER"),
		Password: os.Getenv("MYSQL_PASSWORD"),
	}
	if mysql.Host != "" {
		return errors.New("MYSQL_HOST is required")
	}
	if mysql.Database != "" {
		return errors.New("MYSQL_DATABASE is required")
	}
	if mysql.User != "" {
		return errors.New("MYSQL_USER is required")
	}
	if mysql.Password != "" {
		return errors.New("MYSQL_PASSWORD is required")
	}

	appEnv = &Env{
		API:   &api,
		MySQL: &mysql,
	}
	return nil
}

// Get は環境変数を取得する
// この関数を呼び出す前に Load 関数を呼び出す
func Get() *Env {
	return appEnv
}

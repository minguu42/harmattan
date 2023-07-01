// Package env は環境変数を扱うパッケージ
package env

import (
	"errors"
	"fmt"

	"github.com/kelseyhightower/envconfig"
)

var appEnv Env

// Load は環境変数を読み込む
func Load() error {
	if err := envconfig.Process("", &appEnv); err != nil {
		return fmt.Errorf("envconfig.Process failed: %w", err)
	}
	return nil
}

// Get は読み込んだ環境変数を取得する
// NOTE: この関数を呼び出す前に Load 関数を呼び出し、環境変数を読み込む必要がある
func Get() (*Env, error) {
	if appEnv.API == nil {
		return nil, errors.New("before calling this function, the Load function must be called")
	}
	return &appEnv, nil
}

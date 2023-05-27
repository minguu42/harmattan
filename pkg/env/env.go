// Package env は環境変数を扱うパッケージ
package env

import (
	"fmt"

	"github.com/kelseyhightower/envconfig"
)

var appEnv Env

// Load は環境変数を読み込み、 *Env を返す
func Load() (*Env, error) {
	if err := envconfig.Process("", &appEnv); err != nil {
		return nil, fmt.Errorf("envconfig.Process failed: %w", err)
	}
	return &appEnv, nil
}

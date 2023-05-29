// Package server はサーバに関するパッケージ
package server

import (
	"fmt"
	"net/http"
	"time"

	"github.com/minguu42/mtasks/pkg/app"
	"github.com/minguu42/mtasks/pkg/env"
	"github.com/minguu42/mtasks/pkg/ogen"
)

// NewServer はサーバの初期化する
func NewServer(api *env.API) (*http.Server, error) {
	s, err := ogen.NewServer(&app.Handler{})
	if err != nil {
		return nil, fmt.Errorf("ogen.NewServer failed: %w", err)
	}

	return &http.Server{
		Addr:              fmt.Sprintf(":%d", api.Port),
		Handler:           s,
		ReadTimeout:       10 * time.Second,
		ReadHeaderTimeout: 10 * time.Second,
		MaxHeaderBytes:    1 << 20,
	}, nil
}

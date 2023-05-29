// Package server はサーバに関するパッケージ
package server

import (
	"fmt"
	"net/http"
	"time"

	"github.com/minguu42/mtasks/pkg/ogen"

	"github.com/minguu42/mtasks/pkg/env"

	"github.com/minguu42/mtasks/pkg/app"
)

// NewServer はルーティングの設定、サーバの初期化を行う
func NewServer(api *env.API) *http.Server {
	s, _ := ogen.NewServer(&app.Handler{})

	return &http.Server{
		Addr:              fmt.Sprintf(":%d", api.Port),
		Handler:           s,
		ReadTimeout:       10 * time.Second,
		ReadHeaderTimeout: 10 * time.Second,
		MaxHeaderBytes:    1 << 20,
	}
}

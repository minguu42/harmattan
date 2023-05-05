// Package server はルーティングを含むサーバの設定を記述するパッケージ
package server

import (
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/minguu42/mtasks/pkg/app"
)

func New() *http.Server {
	r := chi.NewRouter()
	r.Get("/health", app.GetHealth())

	return &http.Server{
		Addr:           ":8080",
		Handler:        r,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
}

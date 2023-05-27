// Package server はサーバに関するパッケージ
package server

import (
	"fmt"
	"net/http"
	"time"

	"github.com/minguu42/mtasks/pkg/env"

	"github.com/go-chi/chi/v5"
	"github.com/minguu42/mtasks/pkg/app"
)

// NewServer はルーティングの設定、サーバの初期化を行う
func NewServer(api *env.API) *http.Server {
	r := chi.NewRouter()

	r.Get("/health", app.GetHealth)
	r.Route("/tasks", func(r chi.Router) {
		r.Post("/", app.PostTasks)
		r.Get("/", app.GetTasks)
		r.Patch("/{taskID}", app.PatchTask)
		r.Delete("/{taskID}", app.DeleteTask)
	})

	return &http.Server{
		Addr:              fmt.Sprintf(":%d", api.Port),
		Handler:           r,
		ReadTimeout:       10 * time.Second,
		ReadHeaderTimeout: 10 * time.Second,
		MaxHeaderBytes:    1 << 20,
	}
}

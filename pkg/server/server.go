package server

import (
	"github.com/go-chi/chi/v5"
	"github.com/minguu42/mtasks/pkg/app"
	"net/http"
	"time"
)

// NewServer はルーティングの設定、サーバの初期化を行う
func NewServer() *http.Server {
	r := chi.NewRouter()

	r.Get("/health", app.GetHealth)
	r.Route("/tasks", func(r chi.Router) {
		r.Post("/", app.PostTasks)
		r.Get("/", app.GetTasks)
		r.Patch("/{taskID}", app.PatchTask)
		r.Delete("/{taskID}", app.DeleteTask)
	})

	return &http.Server{
		Addr:              ":8080",
		Handler:           r,
		ReadTimeout:       10 * time.Second,
		ReadHeaderTimeout: 10 * time.Second,
		MaxHeaderBytes:    1 << 20,
	}
}

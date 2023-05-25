// Package api は全てのコードを持つ
package api

import (
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
)

// NewServer は http.Server を返す
func NewServer() *http.Server {
	r := chi.NewRouter()
	r.Get("/health", getHealth())
	r.Post("/tasks", postTasks)

	return &http.Server{
		Addr:              ":8080",
		Handler:           r,
		ReadTimeout:       10 * time.Second,
		ReadHeaderTimeout: 10 * time.Second,
		MaxHeaderBytes:    1 << 20,
	}
}

// Package api は全てのコードを持つ
package api

import (
	"database/sql"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
)

// NewServer は http.Server を返す
func NewServer(_ *sql.DB) *http.Server {
	r := chi.NewRouter()
	r.Get("/health", getHealth())

	return &http.Server{
		Addr:              ":8080",
		Handler:           r,
		ReadTimeout:       10 * time.Second,
		ReadHeaderTimeout: 10 * time.Second,
		MaxHeaderBytes:    1 << 20,
	}
}

package ttime

import (
	"context"
	"net/http"
	"time"
)

type timeKey struct{}

// MiddlewareFixTime はリクエストのコンテキストに固定時刻を含めるミドルウェア
func MiddlewareFixTime(next http.Handler, tm time.Time) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		context.WithValue(r.Context(), timeKey{}, tm)
		next.ServeHTTP(w, r)
	}
}

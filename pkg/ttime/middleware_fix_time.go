package ttime

import (
	"context"
	"net/http"
	"time"
)

type TimeKey struct{}

// MiddlewareFixTime はリクエストのコンテキストに固定時刻を含めるミドルウェア
func MiddlewareFixTime(next http.Handler, tm time.Time) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		next.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), TimeKey{}, tm)))
	}
}

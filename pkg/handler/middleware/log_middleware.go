package middleware

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/minguu42/opepe/pkg/logging"
	"github.com/minguu42/opepe/pkg/ttime"
)

// Log はリクエストの結果をログに出力するミドルウェア
func Log(next http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		t1 := ttime.Now(ctx)
		next.ServeHTTP(w, r)
		t2 := ttime.Now(ctx)

		quote := "\x1b[34m\"\x1b[0m"
		method := fmt.Sprintf("\x1b[35m%s\x1b[0m", r.Method)
		scheme := "http"
		if r.TLS != nil {
			scheme = "https"
		}
		url := fmt.Sprintf("%s://%s%s", scheme, r.Host, r.RequestURI)
		proto := fmt.Sprintf("\x1b[36m%s\x1b[0m", r.Proto)
		remoteAddr := fmt.Sprintf("\x1b[33m%s\x1b[0m", r.RemoteAddr)
		t := fmt.Sprintf("\x1b[32m%sμs\x1b[0m", strconv.FormatInt(t2.Sub(t1).Microseconds(), 10))
		logging.Infof(r.Context(), "%s%s %s %s%s from %s in %s", quote, method, url, proto, quote, remoteAddr, t)
	}
}

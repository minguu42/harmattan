package handler

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/minguu42/mtasks/pkg/logging"
)

// MiddlewareLog -
func MiddlewareLog(next http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		t1 := time.Now()
		next.ServeHTTP(w, r)
		t2 := time.Now()

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
		logging.Infof("%s%s %s %s%s from %s in %s", quote, method, url, proto, quote, remoteAddr, t)
	}
}

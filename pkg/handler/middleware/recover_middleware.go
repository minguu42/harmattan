package middleware

import (
	"encoding/json"
	"fmt"
	"net/http"
	"runtime"

	"github.com/minguu42/opepe/gen/ogen"
	"github.com/minguu42/opepe/pkg/logging"
)

// Recover は
func Recover(next http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if v := recover(); v != nil {
				if v == http.ErrAbortHandler {
					panic(v)
				}

				w.WriteHeader(http.StatusInternalServerError)
				logging.Errorf(r.Context(), fmt.Sprintf("%s\n%s", v, stacktrace()))
				_ = json.NewEncoder(w).Encode(ogen.Error{
					Code:    500,
					Message: "サーバ側で何らかのエラーが発生しました",
				})
			}
		}()
		next.ServeHTTP(w, r)
	}
}

func stacktrace() string {
	m := ""
	for depth := 0; ; depth++ {
		_, file, line, ok := runtime.Caller(depth)
		if !ok {
			break
		}
		m += fmt.Sprintf("%d:%v:%d\n", depth, file, line)
	}
	return m
}

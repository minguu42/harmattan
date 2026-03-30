package api

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"runtime"
	"time"

	"github.com/minguu42/harmattan/internal/api/apierror"
	"github.com/minguu42/harmattan/internal/atel"
	"github.com/minguu42/harmattan/internal/lib/clock"
	"github.com/minguu42/harmattan/internal/lib/errtrace"
	"github.com/ogen-go/ogen/middleware"
)

type requestStartKey struct{}

// setRequestStart はリクエストの開始時刻をコンテキストに付与する
// ogenミドルウェアではセキュリティハンドラの後にしか実行されず正確な開始時刻が分からないためhttpミドルウェアを利用している
func setRequestStart(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		next.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), requestStartKey{}, clock.Now(r.Context()))))
	})
}

// attachTraceID は認証不要のエンドポイント用にトレースIDをロガーに付与する
// 認証が必要なエンドポイントではセキュリティハンドラで先に付与しているが、重複しても影響はない
func attachTraceID() middleware.Middleware {
	return func(req middleware.Request, next middleware.Next) (middleware.Response, error) {
		req.SetContext(atel.ContextWithTracedLogger(req.Context))
		return next(req)
	}
}

const slowRequestThreshold = 1 * time.Second

func accessLog() middleware.Middleware {
	return func(req middleware.Request, next middleware.Next) (middleware.Response, error) {
		// CheckHealthオペレーションのログは出さない
		if req.OperationID == "CheckHealth" {
			return next(req)
		}

		start := clock.Now(req.Context)
		if t, ok := req.Context.Value(requestStartKey{}).(time.Time); ok {
			start = t
		}
		resp, err := next(req)
		duration := clock.Now(req.Context).Sub(start)

		status := 200
		if err != nil {
			status = apierror.ToError(err).Status()
		}

		atel.AccessLog(req.Context, &atel.AccessFields{
			Status:      status,
			Duration:    duration,
			OperationID: req.OperationID,
			Method:      req.Raw.Method,
			URL:         req.Raw.URL.String(),
			Body:        req.Body,
			IPAddress:   req.Raw.RemoteAddr,
			UserAgent:   req.Raw.UserAgent(),
		})
		if status >= 500 {
			atel.AccessErrorLog(req.Context, req.OperationID, err)
		}
		if duration >= slowRequestThreshold {
			atel.AccessSlowLog(req.Context, req.OperationID, status, duration)
		}
		return resp, err
	}
}

func recovery() middleware.Middleware {
	return func(req middleware.Request, next middleware.Next) (resp middleware.Response, err error) {
		defer func() {
			if r := recover(); r != nil {
				message := "panic: "
				switch v := r.(type) {
				case string:
					message += v
				case fmt.Stringer:
					message += v.String()
				case error:
					message += v.Error()
				default:
					message += fmt.Sprintf("%v", v)
				}

				pc := make([]uintptr, errtrace.MaxStackDepth)
				n := runtime.Callers(2, pc)

				err = errtrace.FromStack(apierror.UnknownError(errors.New(message)), pc[:n:n])
			}
		}()

		resp, err = next(req)
		return
	}
}

package middleware

import (
	"time"

	"github.com/minguu42/harmattan/api/applog"
	"github.com/minguu42/harmattan/internal/lib/clock"
	"github.com/ogen-go/ogen/middleware"
)

func AccessLog(l *applog.Logger) middleware.Middleware {
	return func(req middleware.Request, next middleware.Next) (middleware.Response, error) {
		if req.OperationID == "checkHealth" {
			return next(req)
		}

		start := clock.Now(req.Context)
		resp, err := next(req)

		l.Access(req.Context, &applog.AccessFields{
			ExecutionTime: time.Since(start),
			Err:           err,
			OperationID:   req.OperationID,
			Method:        req.Raw.Method,
			URL:           req.Raw.URL.String(),
			Body:          req.Body,
			IPAddress:     req.Raw.RemoteAddr,
		})
		return resp, err
	}
}

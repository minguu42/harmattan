package middleware

import (
	"time"

	"github.com/minguu42/harmattan/internal/api/usecase"
	"github.com/minguu42/harmattan/internal/atel"
	"github.com/minguu42/harmattan/internal/lib/clock"
	"github.com/ogen-go/ogen/middleware"
)

func AccessLog() middleware.Middleware {
	return func(req middleware.Request, next middleware.Next) (middleware.Response, error) {
		// CheckHealthオペレーションのログは出さない
		if req.OperationID == "CheckHealth" {
			return next(req)
		}

		start := clock.Now(req.Context)
		resp, err := next(req)

		status := 200
		if err != nil {
			status = usecase.ToError(err).Status()
		}

		atel.AccessLog(req.Context, &atel.AccessFields{
			Status:      status,
			Duration:    time.Since(start),
			OperationID: req.OperationID,
			Method:      req.Raw.Method,
			URL:         req.Raw.URL.String(),
			Body:        req.Body,
			IPAddress:   req.Raw.RemoteAddr,
			UserAgent:   req.Raw.UserAgent(),
		})
		if status >= 500 {
			atel.ErrorLog(req.Context, "Unexpected error occurred", err)
		}
		return resp, err
	}
}

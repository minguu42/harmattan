package middleware

import (
	"time"

	"github.com/minguu42/harmattan/internal/api/apierror"
	"github.com/minguu42/harmattan/internal/atel"
	"github.com/minguu42/harmattan/internal/lib/clock"
	"github.com/ogen-go/ogen/middleware"
)

const slowRequestThreshold = 1 * time.Second

func AccessLog() middleware.Middleware {
	return func(req middleware.Request, next middleware.Next) (middleware.Response, error) {
		// CheckHealthオペレーションのログは出さない
		if req.OperationID == "CheckHealth" {
			return next(req)
		}

		start := clock.Now(req.Context)
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

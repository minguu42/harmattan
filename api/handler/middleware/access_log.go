package middleware

import (
	"time"

	"github.com/minguu42/harmattan/api/usecase"
	"github.com/minguu42/harmattan/internal/alog"
	"github.com/minguu42/harmattan/internal/lib/clock"
	"github.com/ogen-go/ogen/middleware"
)

func AccessLog(l *alog.Logger) middleware.Middleware {
	return func(req middleware.Request, next middleware.Next) (middleware.Response, error) {
		if req.OperationID == "CheckHealth" {
			return next(req)
		}

		status := 200
		var errorInfo *alog.ErrorInfo
		start := clock.Now(req.Context)

		resp, err := next(req)
		if err != nil {
			appErr := usecase.ToError(err)
			status = appErr.Status()
			errorInfo = &alog.ErrorInfo{
				ErrorMessage: appErr.Error(),
				StackTrace:   appErr.Stacktrace(),
			}
		}

		l.Access(req.Context, &alog.AccessFields{
			Status:        status,
			ErrorInfo:     errorInfo,
			ExecutionTime: time.Since(start),
			OperationID:   req.OperationID,
			Method:        req.Raw.Method,
			URL:           req.Raw.URL.String(),
			Body:          req.Body,
			IPAddress:     req.Raw.RemoteAddr,
		})
		return resp, err
	}
}

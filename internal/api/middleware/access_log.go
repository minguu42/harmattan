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
		start := clock.Now(req.Context)
		resp, err := next(req)

		status := 200
		if err != nil {
			status = usecase.ToError(err).Status()
		}

		// CheckHealthオペレーションの正常系のアクセスログは出さない
		if req.OperationID == "CheckHealth" && status < 500 {
			return resp, err
		}

		atel.AccessLog(req.Context, &atel.AccessFields{
			Status:        status,
			Err:           err,
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

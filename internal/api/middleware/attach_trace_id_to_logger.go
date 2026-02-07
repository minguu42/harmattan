package middleware

import (
	"github.com/minguu42/harmattan/internal/atel"
	"github.com/ogen-go/ogen/middleware"
)

func AttachTraceIDToLogger() middleware.Middleware {
	return func(req middleware.Request, next middleware.Next) (middleware.Response, error) {
		req.SetContext(atel.ContextWithTracedLogger(req.Context))
		return next(req)
	}
}

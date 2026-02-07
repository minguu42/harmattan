package middleware

import (
	"github.com/minguu42/harmattan/internal/alog"
	"github.com/minguu42/harmattan/internal/atel"
	"github.com/ogen-go/ogen/middleware"
)

func AttachRequestIDToLogger() middleware.Middleware {
	return func(req middleware.Request, next middleware.Next) (middleware.Response, error) {
		if traceID := atel.TraceIDFromContext(req.Context); traceID != "" {
			req.SetContext(alog.ContextWithRequestID(req.Context, traceID))
		}
		return next(req)
	}
}

package middleware

import (
	"github.com/minguu42/harmattan/internal/applog"
	"github.com/minguu42/harmattan/internal/lib/idgen"
	"github.com/ogen-go/ogen/middleware"
)

func AttachRequestIDToLogger(l *applog.Logger) middleware.Middleware {
	return func(req middleware.Request, next middleware.Next) (middleware.Response, error) {
		req.SetContext(l.ContextWithRequestID(req.Context, idgen.ULID(req.Context)))
		return next(req)
	}
}

package applog

import (
	"context"
	"log/slog"
	"os"
	"time"

	"github.com/minguu42/harmattan/api/apperr"
	"github.com/minguu42/harmattan/internal/auth"
	"github.com/minguu42/harmattan/internal/lib/logutil"
)

type Logger struct {
	base *slog.Logger
}

func New(indented bool) *Logger {
	opts := &slog.HandlerOptions{ReplaceAttr: func(_ []string, a slog.Attr) slog.Attr {
		if a.Key == slog.MessageKey {
			a.Key = "message"
		}
		return logutil.MaskAttr(a)
	}}
	if indented {
		return &Logger{base: slog.New(logutil.NewJSONIndentHandler(os.Stdout, opts))}
	}
	return &Logger{base: slog.New(slog.NewJSONHandler(os.Stdout, opts))}
}

type loggerKey struct{}

func (l *Logger) logger(ctx context.Context) *slog.Logger {
	if logger, ok := ctx.Value(loggerKey{}).(*slog.Logger); ok {
		return logger
	}
	return l.base
}

func (l *Logger) ContextWithRequestID(ctx context.Context, requestID string) context.Context {
	logger, ok := ctx.Value(loggerKey{}).(*slog.Logger)
	if !ok {
		logger = l.base
	}
	return context.WithValue(ctx, loggerKey{}, logger.With(slog.String("request_id", requestID)))
}

func (l *Logger) Event(ctx context.Context, msg string) {
	l.logger(ctx).Log(ctx, slog.LevelInfo, msg)
}

func (l *Logger) Error(ctx context.Context, msg string) {
	l.logger(ctx).Log(ctx, slog.LevelError, msg)
}

type AccessFields struct {
	ExecutionTime time.Duration
	Err           error

	OperationID string
	Method      string
	URL         string
	Body        any
	IPAddress   string
}

func (l *Logger) Access(ctx context.Context, fields *AccessFields) {
	message := "Request accepted"
	executionTime := slog.Int64("execution_time", fields.ExecutionTime.Milliseconds())
	var request slog.Attr
	if user, ok := auth.UserFromContext(ctx); ok {
		request = slog.Group("request",
			slog.String("user_id", string(user.ID)),
			slog.String("operation_id", fields.OperationID),
			slog.String("method", fields.Method),
			slog.String("url", fields.URL),
			slog.Any("body", fields.Body),
			slog.String("ip_address", fields.IPAddress),
		)
	} else {
		request = slog.Group("request",
			slog.String("operation_id", fields.OperationID),
			slog.String("method", fields.Method),
			slog.String("url", fields.URL),
			slog.Any("body", fields.Body),
			slog.String("ip_address", fields.IPAddress),
		)
	}

	if fields.Err == nil {
		l.logger(ctx).LogAttrs(ctx, slog.LevelInfo, message,
			slog.Int("status_code", 200),
			executionTime,
			request,
		)
		return
	}

	appErr := apperr.ToError(fields.Err)
	status := appErr.StatusCode()
	level := slog.LevelWarn
	if status >= 500 {
		level = slog.LevelError
	}
	if stacktrace := appErr.Stacktrace(); len(stacktrace) > 0 {
		l.logger(ctx).LogAttrs(ctx, level, message,
			slog.Int("status_code", status),
			slog.String("error_message", appErr.Error()),
			slog.Any("stacktrace", stacktrace),
			executionTime,
			request,
		)
		return
	}
	l.logger(ctx).LogAttrs(ctx, level, message,
		slog.Int("status_code", status),
		slog.String("error_message", appErr.Error()),
		executionTime,
		request,
	)
}

package applog

import (
	"context"
	"log/slog"
	"os"
	"strings"
	"time"

	"github.com/minguu42/harmattan/api/apperr"
	"github.com/minguu42/harmattan/lib/slogdebug"
)

type Logger struct {
	base *slog.Logger
}

func New() *Logger {
	opts := &slog.HandlerOptions{
		ReplaceAttr: func(_ []string, a slog.Attr) slog.Attr {
			if a.Key == slog.MessageKey {
				a.Key = "message"
			}
			return a
		},
	}
	if strings.ToLower(os.Getenv("USE_DEBUG_LOGGER")) == "true" {
		return &Logger{base: slog.New(slogdebug.NewJSONIndentHandler(os.Stdout, opts))}
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
	RemoteAddr  string
}

func (l *Logger) Access(ctx context.Context, fields *AccessFields) {
	message := "Request accepted"
	executionTime := slog.Int64("execution_time", fields.ExecutionTime.Milliseconds())
	request := slog.Group("request",
		slog.String("operation_id", fields.OperationID),
		slog.String("method", fields.Method),
		slog.String("url", fields.URL),
		slog.String("remote_addr", fields.RemoteAddr),
	)
	if fields.Err == nil {
		l.logger(ctx).LogAttrs(ctx, slog.LevelInfo, message, executionTime, request)
		return
	}

	appErr := apperr.ToError(fields.Err)
	status := appErr.APIError().StatusCode
	level := slog.LevelWarn
	if status >= 500 {
		level = slog.LevelError
	}
	l.logger(ctx).LogAttrs(ctx, level, message, executionTime, request,
		slog.Int("status_code", status),
		slog.String("error_message", appErr.Error()),
	)
}

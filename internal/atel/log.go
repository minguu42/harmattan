package atel

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"os"
	"runtime"
	"time"

	"github.com/minguu42/harmattan/internal/domain"
	"github.com/minguu42/harmattan/internal/lib/errtrace"
)

func EventLog(ctx context.Context, message string) {
	logger(ctx).base.Log(ctx, slog.LevelInfo, message)
}

func ErrorLog(ctx context.Context, message string, err error) {
	logger(ctx).base.LogAttrs(ctx, slog.LevelError, message, slog.Any("error", err))
}

type AccessFields struct {
	Status      int
	Duration    time.Duration
	OperationID string
	Method      string
	URL         string
	Body        any
	IPAddress   string
	UserAgent   string
}

func AccessLog(ctx context.Context, fields *AccessFields) {
	attrs := make([]slog.Attr, 0, 9)
	attrs = append(attrs,
		slog.Int("response.status_code", fields.Status),
		slog.Int64("response.duration", fields.Duration.Milliseconds()),
		slog.String("request.operation", fields.OperationID),
		slog.String("request.method", fields.Method),
		slog.String("request.url", fields.URL),
	)
	if fields.Body != nil {
		attrs = append(attrs, slog.Any("request.body", fields.Body))
	}
	if user, err := domain.UserFromContext(ctx); err == nil {
		attrs = append(attrs, slog.String("request.user_id", string(user.ID)))
	}
	attrs = append(attrs,
		slog.String("request.ip_address", fields.IPAddress),
		slog.String("request.user_agent", fields.UserAgent),
	)
	logger(ctx).base.LogAttrs(ctx, slog.LevelInfo, "Request processed", attrs...)
}

func AccessErrorLog(ctx context.Context, operationID string, err error) {
	logger(ctx).base.LogAttrs(ctx, slog.LevelError, "Unexpected error occurred",
		slog.String("request.operation", operationID),
		slog.Any("error", err),
	)
}

func AccessSlowLog(ctx context.Context, operationID string, status int, duration time.Duration) {
	logger(ctx).base.LogAttrs(ctx, slog.LevelWarn, "Slow request detected",
		slog.String("request.operation", operationID),
		slog.Int("response.status_code", status),
		slog.Int64("response.duration", duration.Milliseconds()),
	)
}

func SQLLog(ctx context.Context, fc func() (sql string, rowsAffected int64)) {
	level := slog.LevelDebug
	if !logger(ctx).base.Enabled(ctx, level) {
		return
	}

	query, _ := fc()
	var loc string
	if _, file, line, ok := runtime.Caller(4); ok {
		loc = fmt.Sprintf("%s:%d", file, line)
	}
	logger(ctx).base.LogAttrs(ctx, level, query, slog.String("location", loc))
}

func Capture(ctx context.Context, message string) func(func() error) {
	pc := make([]uintptr, errtrace.MaxStackDepth)
	n := runtime.Callers(2, pc)

	return func(f func() error) {
		if f == nil {
			return
		}
		err := f()
		if err == nil {
			return
		}

		// os.ErrClosedは2重にクローズ処理を行った場合に返され、処理的には問題がないため無視する。
		if errors.Is(err, os.ErrClosed) {
			return
		}

		logger(ctx).base.LogAttrs(ctx, slog.LevelWarn, message, slog.Any("error", errtrace.FromStack(err, pc[:n:n])))
	}
}

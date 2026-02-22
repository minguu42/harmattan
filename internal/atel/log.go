package atel

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"os"
	"runtime"
	"time"

	"github.com/minguu42/harmattan/internal/auth"
	"github.com/minguu42/harmattan/internal/lib/errtrace"
)

const (
	reset = "\033[0m"
	green = "\033[32m"
	cyan  = "\033[36m"
)

func EventLog(ctx context.Context, message string) {
	logger(ctx).base.Log(ctx, slog.LevelInfo, message)
}

func ErrorLog(ctx context.Context, message string, err error) {
	logger(ctx).base.LogAttrs(ctx, slog.LevelError, message, errorToAttrs(err)...)
}

func errorToAttrs(err error) []slog.Attr {
	if err == nil {
		return nil
	}

	if serr := new(errtrace.StackError); errors.As(err, &serr) {
		return []slog.Attr{
			slog.String("error.message", serr.Error()),
			slog.Any("error.frames", serr.Frames()),
		}
	}
	return []slog.Attr{slog.String("error.message", err.Error())}
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
	if user, ok := auth.UserFromContext(ctx); ok {
		attrs = append(attrs, slog.String("request.user_id", string(user.ID)))
	}
	attrs = append(attrs,
		slog.String("request.ip_address", fields.IPAddress),
		slog.String("request.user_agent", fields.UserAgent),
	)

	logger(ctx).base.LogAttrs(ctx, slog.LevelInfo, "Request processed", attrs...)
}

func AccessErrorLog(ctx context.Context, operationID string, err error) {
	attrs := make([]slog.Attr, 0, 3)
	attrs = append(attrs, slog.String("request.operation", operationID))
	attrs = append(attrs, errorToAttrs(err)...)
	logger(ctx).base.LogAttrs(ctx, slog.LevelError, "Unexpected error occurred", attrs...)
}

func AccessSlowLog(ctx context.Context, operationID string, status int, duration time.Duration) {
	logger(ctx).base.LogAttrs(ctx, slog.LevelWarn, "Slow request detected",
		slog.String("request.operation", operationID),
		slog.Int("response.status_code", status),
		slog.Int64("response.duration", duration.Milliseconds()),
	)
}

func SQLLog(ctx context.Context, begin time.Time, fc func() (sql string, rowsAffected int64)) {
	if logger(ctx).level.Level() > slog.LevelDebug {
		return
	}

	query, _ := fc()
	ms := float64(time.Since(begin)) / float64(time.Millisecond)
	var loc string
	if _, file, line, ok := runtime.Caller(4); ok {
		loc = fmt.Sprintf("%s:%d", file, line)
	}

	if logger(ctx).prettyPrint {
		fmt.Printf("\n-- %s[%.3fms]%s %s%s%s\n%s\n",
			green, ms, reset,
			cyan, loc, reset,
			query,
		)
	} else {
		logger(ctx).base.LogAttrs(ctx, slog.LevelDebug, "",
			slog.String("elapsed(ms)", fmt.Sprintf("%.3f", ms)),
			slog.String("location", loc),
			slog.String("query", query),
		)
	}
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

		logger(ctx).base.LogAttrs(ctx, slog.LevelWarn, message, errorToAttrs(errtrace.FromStack(err, pc[:n:n]))...)
	}
}

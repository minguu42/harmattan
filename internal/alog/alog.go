package alog

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

func Event(ctx context.Context, message string) {
	logger(ctx).base.Log(ctx, slog.LevelInfo, message)
}

func Error(ctx context.Context, message string, err error) {
	logger(ctx).base.LogAttrs(ctx, slog.LevelError, message, slog.Any("error", err))
}

type AccessFields struct {
	Status        int
	Err           error
	ExecutionTime time.Duration

	OperationID string
	Method      string
	URL         string
	Body        any
	IPAddress   string
}

func Access(ctx context.Context, fields *AccessFields) {
	level := slog.LevelInfo
	switch {
	case 400 <= fields.Status && fields.Status < 500:
		level = slog.LevelWarn
	case 500 <= fields.Status && fields.Status < 600:
		level = slog.LevelError
	}

	attrs := make([]slog.Attr, 0, 5)
	attrs = append(attrs, slog.Int("status_code", fields.Status))
	if fields.Err != nil {
		if serr := new(errtrace.StackError); errors.As(fields.Err, &serr) {
			attrs = append(attrs, slog.Any("error", serr))
		} else {
			attrs = append(attrs, slog.Group("error", slog.String("message", fields.Err.Error())))
		}
	}
	attrs = append(attrs, slog.Int64("execution_time", fields.ExecutionTime.Milliseconds()))

	requestAttrs := make([]slog.Attr, 0, 6)
	if user, ok := auth.UserFromContext(ctx); ok {
		requestAttrs = append(requestAttrs, slog.String("user_id", string(user.ID)))
	}
	requestAttrs = append(requestAttrs,
		slog.String("operation_id", fields.OperationID),
		slog.String("method", fields.Method),
		slog.String("url", fields.URL),
	)
	if fields.Body != nil {
		requestAttrs = append(requestAttrs, slog.Any("body", fields.Body))
	}
	requestAttrs = append(requestAttrs, slog.String("ip_address", fields.IPAddress))
	attrs = append(attrs, slog.GroupAttrs("request", requestAttrs...))

	logger(ctx).base.LogAttrs(ctx, level, "Request accepted", attrs...)
}

func Capture(ctx context.Context, message string) func(func() error) {
	pc := make([]uintptr, errtrace.MaxStackDepth)
	n := runtime.Callers(2, pc)

	return func(f func() error) {
		err := f()
		if err == nil {
			return
		}

		// os.ErrClosedは2重にクローズ処理を行った場合に返され、処理的には問題がないため無視する。
		if errors.Is(err, os.ErrClosed) {
			return
		}

		Error(ctx, message, errtrace.FromStack(err, pc[:n:n]))
	}
}

func GORMTrace(ctx context.Context, begin time.Time, fc func() (sql string, rowsAffected int64)) {
	if logger(ctx).level.Level() < slog.LevelDebug {
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

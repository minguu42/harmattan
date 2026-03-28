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

const (
	reset   = "\033[0m"
	gray    = "\033[90m"
	red     = "\033[31m"
	green   = "\033[32m"
	yellow  = "\033[33m"
	blue    = "\033[34m"
	magenta = "\033[35m"
	cyan    = "\033[36m"
)

func EventLog(ctx context.Context, message string) {
	level := slog.LevelInfo
	if logger(ctx).prettyPrint && logger(ctx).base.Enabled(ctx, level) {
		fmt.Printf("%s %s %s\n", prettyTimestamp(), prettyLevel(level), message)
		return
	}
	logger(ctx).base.Log(ctx, level, message)
}

func ErrorLog(ctx context.Context, message string, err error) {
	level := slog.LevelError
	if logger(ctx).prettyPrint && logger(ctx).base.Enabled(ctx, level) {
		fmt.Printf("%s %s %s\n%+v\n", prettyTimestamp(), prettyLevel(level), message, err)
		return
	}
	logger(ctx).base.LogAttrs(ctx, level, message, errorToAttrs(err)...)
}

func errorToAttrs(err error) []slog.Attr {
	if err == nil {
		return nil
	}

	if stackErr, ok := errors.AsType[*errtrace.StackError](err); ok {
		attrs := []slog.Attr{
			slog.String("error.message", stackErr.Error()),
			slog.Any("error.frames", stackErr.Frames()),
		}
		if errAttrs := stackErr.Attrs(); len(errAttrs) > 0 {
			attrs = append(attrs, slog.GroupAttrs("error.attrs", errAttrs...))
		}
		return attrs
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
	level := slog.LevelInfo
	if logger(ctx).prettyPrint && logger(ctx).base.Enabled(ctx, level) {
		status := fmt.Sprintf("%d", fields.Status)
		switch {
		case fields.Status >= 200 && fields.Status < 300:
			status = fmt.Sprintf("%s%d%s", green, fields.Status, reset)
		case fields.Status >= 400 && fields.Status < 500:
			status = fmt.Sprintf("%s%d%s", yellow, fields.Status, reset)
		case fields.Status >= 500:
			status = fmt.Sprintf("%s%d%s", red, fields.Status, reset)
		}
		fmt.Printf("%s %s %s %s %s\n", prettyTimestamp(), prettyLevel(level), fields.OperationID, status, prettyDuration(fields.Duration))
		return
	}

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
	logger(ctx).base.LogAttrs(ctx, level, "Request processed", attrs...)
}

func AccessErrorLog(ctx context.Context, operationID string, err error) {
	level := slog.LevelError
	if logger(ctx).prettyPrint && logger(ctx).base.Enabled(ctx, level) {
		fmt.Printf("%s %s %s\n%+v\n", prettyTimestamp(), prettyLevel(level), operationID, err)
		return
	}

	attrs := make([]slog.Attr, 0, 3)
	attrs = append(attrs, slog.String("request.operation", operationID))
	attrs = append(attrs, errorToAttrs(err)...)
	logger(ctx).base.LogAttrs(ctx, level, "Unexpected error occurred", attrs...)
}

func AccessSlowLog(ctx context.Context, operationID string, status int, duration time.Duration) {
	level := slog.LevelWarn
	if logger(ctx).prettyPrint && logger(ctx).base.Enabled(ctx, level) {
		fmt.Printf("%s %s %s %s (slow)\n", prettyTimestamp(), prettyLevel(level), operationID, prettyDuration(duration))
		return
	}
	logger(ctx).base.LogAttrs(ctx, level, "Slow request detected",
		slog.String("request.operation", operationID),
		slog.Int("response.status_code", status),
		slog.Int64("response.duration", duration.Milliseconds()),
	)
}

func SQLLog(ctx context.Context, begin time.Time, fc func() (sql string, rowsAffected int64)) {
	level := slog.LevelDebug
	if !logger(ctx).base.Enabled(ctx, level) {
		return
	}

	query, _ := fc()
	duration := time.Since(begin)
	var loc string
	if _, file, line, ok := runtime.Caller(4); ok {
		loc = fmt.Sprintf("%s:%d", file, line)
	}

	if logger(ctx).prettyPrint {
		fmt.Printf("%s %s %s%s%s %s\n%s\n", prettyTimestamp(), prettyLevel(level), blue, loc, reset, prettyDuration(duration), query)
		return
	}
	logger(ctx).base.LogAttrs(ctx, level, "",
		slog.Int64("duration", duration.Milliseconds()),
		slog.String("location", loc),
		slog.String("query", query),
	)
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

func prettyTimestamp() string {
	return fmt.Sprintf("%s%s%s", gray, time.Now().Format("2006/01/02 15:04:05"), reset)
}

func prettyLevel(level slog.Level) string {
	switch level {
	case slog.LevelDebug:
		return fmt.Sprintf("%s%s%s", cyan, level.String(), reset)
	case slog.LevelInfo:
		return fmt.Sprintf("%s%s%s", green, level.String(), reset)
	case slog.LevelWarn:
		return fmt.Sprintf("%s%s%s", yellow, level.String(), reset)
	case slog.LevelError:
		return fmt.Sprintf("%s%s%s", red, level.String(), reset)
	default:
		return level.String()
	}
}

func prettyDuration(duration time.Duration) string {
	ms := float64(duration) / float64(time.Millisecond)
	return fmt.Sprintf("%s%.2fms%s", magenta, ms, reset)
}

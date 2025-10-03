package alog

import (
	"context"
	"log/slog"
	"os"
	"time"

	"github.com/minguu42/harmattan/internal/auth"
)

type Level string

const (
	LevelDebug  = "debug"
	LevelInfo   = "info"
	LevelSilent = "silent"
)

func (l Level) Level() slog.Level {
	switch l {
	case LevelDebug:
		return slog.LevelDebug
	case LevelInfo:
		return slog.LevelInfo
	case LevelSilent:
		return 12
	default:
		return slog.LevelInfo
	}
}

type Logger struct {
	base  *slog.Logger
	Level Level
}

func New(level Level, indented bool) *Logger {
	opts := &slog.HandlerOptions{
		Level: level,
		ReplaceAttr: func(_ []string, a slog.Attr) slog.Attr {
			if a.Key == slog.MessageKey {
				a.Key = "message"
			}
			return MaskAttr(a)
		}}
	if indented {
		return &Logger{base: slog.New(NewJSONIndentHandler(os.Stdout, opts)), Level: level}
	}
	return &Logger{base: slog.New(slog.NewJSONHandler(os.Stdout, opts)), Level: level}
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
	StatusCode    int
	ErrorInfo     *ErrorInfo
	ExecutionTime time.Duration

	OperationID string
	Method      string
	URL         string
	Body        any
	IPAddress   string
}

type ErrorInfo struct {
	ErrorMessage string
	StackTrace   []string
}

func (l *Logger) Access(ctx context.Context, fields *AccessFields) {
	level := slog.LevelInfo
	switch {
	case 400 <= fields.StatusCode && fields.StatusCode < 500:
		level = slog.LevelWarn
	case 500 <= fields.StatusCode && fields.StatusCode < 600:
		level = slog.LevelError
	}

	attrs := make([]slog.Attr, 0, 5)
	attrs = append(attrs, slog.Int("status_code", fields.StatusCode))
	if fields.ErrorInfo != nil {
		attrs = append(attrs, slog.String("error_message", fields.ErrorInfo.ErrorMessage))
		if len(fields.ErrorInfo.StackTrace) != 0 {
			attrs = append(attrs, slog.Any("stacktrace", fields.ErrorInfo.StackTrace))
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

	l.logger(ctx).LogAttrs(ctx, level, "Request accepted", attrs...)
}

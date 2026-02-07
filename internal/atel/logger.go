package atel

import (
	"context"
	"io"
	"log/slog"
	"os"

	"go.opentelemetry.io/otel/trace"
)

type Level string

const (
	LevelDebug Level = "debug"
	LevelInfo  Level = "info"
	LevelWarn  Level = "warn"
	LevelError Level = "error"
)

func (l Level) Level() slog.Level {
	switch l {
	case LevelDebug:
		return slog.LevelDebug
	case LevelInfo:
		return slog.LevelInfo
	case LevelWarn:
		return slog.LevelWarn
	case LevelError:
		return slog.LevelError
	default:
		return slog.LevelInfo
	}
}

var globalLogger = New(os.Stdout, LevelInfo, false)

func SetLogger(l *Logger) {
	globalLogger = l
}

type Logger struct {
	base        *slog.Logger
	level       Level
	prettyPrint bool
}

func New(w io.Writer, level Level, prettyPrint bool) *Logger {
	opts := &slog.HandlerOptions{
		Level: level.Level(),
		ReplaceAttr: func(_ []string, a slog.Attr) slog.Attr {
			if a.Key == slog.MessageKey {
				a.Key = "message"
			}
			return MaskAttr(a)
		}}
	if prettyPrint {
		return &Logger{base: slog.New(NewJSONIndentHandler(w, opts)), level: level, prettyPrint: prettyPrint}
	}
	return &Logger{base: slog.New(slog.NewJSONHandler(w, opts)), level: level, prettyPrint: prettyPrint}
}

type loggerKey struct{}

func logger(ctx context.Context) *Logger {
	if l, ok := ctx.Value(loggerKey{}).(*Logger); ok {
		return l
	}
	return globalLogger
}

func ContextWithTracedLogger(ctx context.Context) context.Context {
	spanContext := trace.SpanContextFromContext(ctx)
	if !spanContext.IsValid() {
		return ctx
	}

	l := logger(ctx)
	return context.WithValue(ctx, loggerKey{}, &Logger{
		base:        l.base.With(slog.String("request_id", spanContext.TraceID().String())),
		level:       l.level,
		prettyPrint: l.prettyPrint,
	})
}

package atel

import (
	"context"
	"io"
	"log/slog"
	"os"

	"go.opentelemetry.io/otel/trace"
)

var globalLogger = New(os.Stdout, slog.LevelInfo, false)

func SetLogger(l *Logger) {
	globalLogger = l
}

type Logger struct {
	base        *slog.Logger
	prettyPrint bool
	traced      bool
}

func New(w io.Writer, level slog.Level, prettyPrint bool) *Logger {
	opts := &slog.HandlerOptions{
		Level: level,
		ReplaceAttr: func(_ []string, a slog.Attr) slog.Attr {
			if a.Key == slog.MessageKey {
				a.Key = "message"
			}
			return MaskAttr(a)
		}}
	if prettyPrint {
		return &Logger{base: slog.New(NewJSONIndentHandler(w, opts)), prettyPrint: prettyPrint}
	}
	return &Logger{base: slog.New(slog.NewJSONHandler(w, opts)), prettyPrint: prettyPrint}
}

type loggerKey struct{}

func logger(ctx context.Context) *Logger {
	if l, ok := ctx.Value(loggerKey{}).(*Logger); ok {
		return l
	}
	return globalLogger
}

func ContextWithTracedLogger(ctx context.Context) context.Context {
	if l, ok := ctx.Value(loggerKey{}).(*Logger); ok && l.traced {
		return ctx
	}

	spanContext := trace.SpanContextFromContext(ctx)
	if !spanContext.IsValid() {
		return ctx
	}

	l := logger(ctx)
	return context.WithValue(ctx, loggerKey{}, &Logger{
		base:        l.base.With(slog.String("trace_id", spanContext.TraceID().String())),
		prettyPrint: l.prettyPrint,
		traced:      true,
	})
}

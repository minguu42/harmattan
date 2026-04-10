package atel

import (
	"context"
	"errors"
	"io"
	"log/slog"
	"os"

	"github.com/minguu42/harmattan/internal/lib/errtrace"
	"go.opentelemetry.io/otel/trace"
)

var globalLogger = New(os.Stdout, slog.LevelInfo, false)

func SetLogger(l *Logger) {
	globalLogger = l
}

type Logger struct {
	base   *slog.Logger
	traced bool
}

func New(w io.Writer, level slog.Level, prettyPrint bool) *Logger {
	if prettyPrint {
		ignoreKeys := []string{
			"trace_id",
			"request.user_id",
			"request.user_agent",
			"request.ip_address",
		}
		return &Logger{base: slog.New(NewColoredTextHandler(w, level, false, ignoreKeys))}
	}
	opts := &slog.HandlerOptions{
		Level: level,
		ReplaceAttr: func(_ []string, a slog.Attr) slog.Attr {
			if a.Key == slog.MessageKey {
				a.Key = "message"
			}
			if a.Key == "error" {
				return expandErrorAttr(a)
			}
			return MaskAttr(a)
		},
	}
	return &Logger{base: slog.New(slog.NewJSONHandler(w, opts))}
}

func expandErrorAttr(attr slog.Attr) slog.Attr {
	err, ok := attr.Value.Any().(error)
	if !ok {
		return attr
	}

	if stackErr, ok := errors.AsType[*errtrace.StackError](err); ok {
		attrs := []slog.Attr{
			slog.String("message", stackErr.Error()),
			slog.Any("frames", stackErr.Frames()),
		}
		if errAttrs := stackErr.Attrs(); len(errAttrs) > 0 {
			attrs = append(attrs, slog.GroupAttrs("attrs", errAttrs...))
		}
		return slog.GroupAttrs("error", attrs...)
	}
	return slog.GroupAttrs("error", slog.String("message", err.Error()))
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

	return context.WithValue(ctx, loggerKey{}, &Logger{
		base:   logger(ctx).base.With(slog.String("trace_id", spanContext.TraceID().String())),
		traced: true,
	})
}

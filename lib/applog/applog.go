package applog

import (
	"context"
	"log/slog"
	"os"
	"strings"

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

package applog

import (
	"context"
	"log/slog"
	"os"

	"github.com/minguu42/harmattan/lib/slogdebug"
)

var applicationLogger *slog.Logger

func init() {
	opts := &slog.HandlerOptions{
		ReplaceAttr: func(_ []string, a slog.Attr) slog.Attr {
			if a.Key == slog.MessageKey {
				a.Key = "message"
			}
			return a
		},
	}
	if os.Getenv("USE_DEBUG_LOGGER") == "true" {
		applicationLogger = slog.New(slogdebug.NewJSONIndentHandler(os.Stdout, opts))
	} else {
		applicationLogger = slog.New(slog.NewJSONHandler(os.Stdout, opts))
	}
}

type loggerKey struct{}

func logger(ctx context.Context) *slog.Logger {
	v, ok := ctx.Value(loggerKey{}).(*slog.Logger)
	if ok {
		return v
	}
	return applicationLogger
}

func Event(ctx context.Context, msg string) {
	logger(ctx).Log(ctx, slog.LevelInfo, msg)
}

func Error(ctx context.Context, msg string) {
	logger(ctx).Log(ctx, slog.LevelError, msg)
}

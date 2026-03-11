package api

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"os"

	"github.com/minguu42/harmattan/internal/atel"
	"github.com/minguu42/harmattan/internal/auth"
	"github.com/minguu42/harmattan/internal/database"
	"github.com/minguu42/harmattan/internal/lib/errtrace"
	"go.opentelemetry.io/otel/sdk/trace"
)

type Factory struct {
	Auth                   *auth.Authenticator
	DB                     *database.Client
	ShutdownTracerProvider func() error
}

func NewFactory(ctx context.Context, conf *Config) (*Factory, error) {
	authn, err := auth.NewAuthenticator(conf.IDTokenSecret, conf.IDTokenExpiration)
	if err != nil {
		return nil, errtrace.Wrap(err)
	}

	db, err := database.NewClient(ctx, &database.Config{
		DSN: database.DSN{
			Host:     conf.DBHost,
			Port:     conf.DBPort,
			Database: conf.DBDatabase,
			User:     conf.DBUser,
			Password: conf.DBPassword,
		},
		MaxOpenConns:    conf.DBMaxOpenConns,
		MaxIdleConns:    conf.DBMaxIdleConns,
		ConnMaxLifetime: conf.DBConnMaxLifetime,
	})
	if err != nil {
		return nil, errtrace.Wrap(err)
	}

	var level slog.Level
	switch conf.LogLevel {
	case "debug":
		level = slog.LevelDebug
	case "info":
		level = slog.LevelInfo
	case "warn":
		level = slog.LevelWarn
	case "error":
		level = slog.LevelError
	default:
		return nil, errtrace.Wrap(fmt.Errorf("invalid log level: %s", conf.LogLevel))
	}
	atel.SetLogger(atel.New(os.Stdout, level, conf.LogPrettyPrint))

	var exporter trace.SpanExporter
	switch conf.TraceExporter {
	case "otlp":
		exporter, err = atel.NewOTLPExporter(ctx, conf.TraceCollectorHost, conf.TraceCollectorPort)
		if err != nil {
			return nil, errtrace.Wrap(err)
		}
	case "stdout":
		exporter, err = atel.NewStdoutExporter()
		if err != nil {
			return nil, errtrace.Wrap(err)
		}
	}
	shutdown, err := atel.SetupTracerProvider(ctx, exporter)
	if err != nil {
		return nil, errtrace.Wrap(err)
	}
	return &Factory{
		Auth:                   authn,
		DB:                     db,
		ShutdownTracerProvider: shutdown,
	}, nil
}

func (f *Factory) Close() error {
	dbErr := f.DB.Close()
	traceErr := f.ShutdownTracerProvider()
	return errtrace.Wrap(errors.Join(dbErr, traceErr))
}

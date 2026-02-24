package api

import (
	"context"
	"errors"
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

	db, err := database.NewClient(ctx, conf.DB)
	if err != nil {
		return nil, errtrace.Wrap(err)
	}

	atel.SetLogger(atel.New(os.Stdout, conf.LogLevel, conf.LogPrettyPrint))

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

package atel

import (
	"context"
	"time"

	"github.com/minguu42/harmattan/internal/lib/errtrace"
	"go.opentelemetry.io/contrib/detectors/aws/ecs"
	"go.opentelemetry.io/contrib/propagators/aws/xray"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/exporters/stdout/stdouttrace"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
)

func SetupTracerProvider(ctx context.Context, exporter trace.SpanExporter) (func(context.Context) error, error) {
	res, err := newResource(ctx)
	if err != nil {
		return nil, errtrace.Wrap(err)
	}

	provider := trace.NewTracerProvider(
		trace.WithIDGenerator(xray.NewIDGenerator()),
		trace.WithResource(res),
	)
	if exporter != nil {
		provider = trace.NewTracerProvider(
			trace.WithBatcher(exporter, trace.WithBatchTimeout(time.Second)),
			trace.WithIDGenerator(xray.NewIDGenerator()),
			trace.WithResource(res),
		)
	}
	otel.SetTracerProvider(provider)
	otel.SetTextMapPropagator(xray.Propagator{})
	return provider.Shutdown, nil
}

func NewOTLPExporter(ctx context.Context) (trace.SpanExporter, error) {
	exporter, err := otlptracegrpc.New(ctx,
		otlptracegrpc.WithInsecure(),
		otlptracegrpc.WithEndpoint("otel-collector:4317"),
	)
	if err != nil {
		return nil, errtrace.Wrap(err)
	}
	return exporter, nil
}

func NewStdoutExporter() (trace.SpanExporter, error) {
	exporter, err := stdouttrace.New(stdouttrace.WithPrettyPrint())
	if err != nil {
		return nil, errtrace.Wrap(err)
	}
	return exporter, nil
}

func newResource(ctx context.Context) (*resource.Resource, error) {
	service := resource.NewSchemaless(attribute.String("service.name", "harmattan-api"))
	base, err := resource.Merge(resource.Default(), service)
	if err != nil {
		return nil, errtrace.Wrap(err)
	}

	ecsResource, err := ecs.NewResourceDetector().Detect(ctx)
	if err != nil {
		return nil, errtrace.Wrap(err)
	}
	merged, err := resource.Merge(base, ecsResource)
	if err != nil {
		return nil, errtrace.Wrap(err)
	}
	return merged, nil
}

package atel_test

import (
	"context"
	"testing"

	"github.com/minguu42/harmattan/internal/atel"
	"go.opentelemetry.io/otel"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
)

func TestTraceIDFromContext(t *testing.T) {
	ctx := context.Background()

	traceID := atel.TraceIDFromContext(ctx)
	if traceID != "" {
		t.Errorf("expected empty string, got %s", traceID)
	}

	provider := sdktrace.NewTracerProvider()
	otel.SetTracerProvider(provider)
	defer func() {
		_ = provider.Shutdown(ctx)
	}()

	tracer := otel.Tracer("test")
	ctx, span := tracer.Start(ctx, "test-span")
	defer span.End()

	traceID = atel.TraceIDFromContext(ctx)
	if traceID == "" {
		t.Error("expected non-empty trace ID")
	}
	if len(traceID) != 32 {
		t.Errorf("expected 32 characters, got %d", len(traceID))
	}
}

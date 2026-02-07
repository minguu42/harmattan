package atel

import (
	"context"

	"go.opentelemetry.io/otel/trace"
)

func TraceIDFromContext(ctx context.Context) string {
	spanContext := trace.SpanContextFromContext(ctx)
	if !spanContext.IsValid() {
		return ""
	}
	return spanContext.TraceID().String()
}

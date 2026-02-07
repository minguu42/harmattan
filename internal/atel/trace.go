package atel

import (
	"context"

	"go.opentelemetry.io/otel/trace"
)

func TraceIDFromContext(ctx context.Context) string {
	spanCtx := trace.SpanContextFromContext(ctx)
	if !spanCtx.IsValid() {
		return ""
	}
	return spanCtx.TraceID().String()
}

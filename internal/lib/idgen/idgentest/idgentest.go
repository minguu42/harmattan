package idgentest

import (
	"context"
	"testing"

	"github.com/minguu42/harmattan/internal/lib/idgen/internal"
)

type ulidKey struct{}

func init() {
	if testing.Testing() {
		internal.ULID = ulidForTest
	}
}

func ulidForTest(ctx context.Context) string {
	if ulid, ok := ctx.Value(ulidKey{}).(string); ok {
		return ulid
	}
	return internal.DefaultULID(ctx)
}

func WithFixedULID(ctx context.Context, ulid string) context.Context {
	return context.WithValue(ctx, ulidKey{}, ulid)
}

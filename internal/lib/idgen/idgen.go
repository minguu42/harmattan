package idgen

import (
	"context"
	"testing"

	"github.com/oklog/ulid/v2"
)

var internalULID = func(_ context.Context) string { return ulid.Make().String() }

func init() {
	if testing.Testing() {
		internalULID = ulidForTest
	}
}

func ULID(ctx context.Context) string {
	return internalULID(ctx)
}

type ulidKey struct{}

func ulidForTest(ctx context.Context) string {
	if v, ok := ctx.Value(ulidKey{}).(string); ok {
		return v
	}
	return ulid.Make().String()
}

func WithFixedULID(ctx context.Context, id string) context.Context {
	return context.WithValue(ctx, ulidKey{}, id)
}

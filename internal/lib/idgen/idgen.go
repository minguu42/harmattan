package idgen

import (
	"context"
	"testing"

	"github.com/oklog/ulid/v2"
)

type ulidKey struct{}

func ULID(ctx context.Context) string {
	if testing.Testing() {
		if v, ok := ctx.Value(ulidKey{}).(string); ok {
			return v
		}
	}
	return ulid.Make().String()
}

func WithFixedULID(ctx context.Context, id string) context.Context {
	return context.WithValue(ctx, ulidKey{}, id)
}

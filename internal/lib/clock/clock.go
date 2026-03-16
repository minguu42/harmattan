package clock

import (
	"context"
	"testing"
	"time"
)

type nowKey struct{}

func Now(ctx context.Context) time.Time {
	if testing.Testing() {
		if v, ok := ctx.Value(nowKey{}).(time.Time); ok {
			return v
		}
	}
	return time.Now()
}

func WithFixedNow(ctx context.Context, tm time.Time) context.Context {
	return context.WithValue(ctx, nowKey{}, tm)
}

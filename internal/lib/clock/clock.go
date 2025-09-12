package clock

import (
	"context"
	"testing"
	"time"
)

var internalNow = func(_ context.Context) time.Time { return time.Now() }

func init() {
	if testing.Testing() {
		internalNow = nowForTest
	}
}

func Now(ctx context.Context) time.Time {
	return internalNow(ctx)
}

type nowKey struct{}

func nowForTest(ctx context.Context) time.Time {
	if v, ok := ctx.Value(nowKey{}).(time.Time); ok {
		return v
	}
	return time.Now()
}

func WithFixedNow(ctx context.Context, tm time.Time) context.Context {
	return context.WithValue(ctx, nowKey{}, tm)
}

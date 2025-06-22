package clocktest

import (
	"context"
	"testing"
	"time"

	"github.com/minguu42/harmattan/lib/clock/internal"
)

type nowKey struct{}

func init() {
	if testing.Testing() {
		internal.Now = nowForTest
	}
}

func nowForTest(ctx context.Context) time.Time {
	if now, ok := ctx.Value(nowKey{}).(time.Time); ok {
		return now
	}
	return internal.DefaultNow(ctx)
}

func WithFixedNow(ctx context.Context, tm time.Time) context.Context {
	return context.WithValue(ctx, nowKey{}, tm)
}

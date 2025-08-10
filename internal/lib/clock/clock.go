package clock

import (
	"context"
	"time"

	"github.com/minguu42/harmattan/internal/lib/clock/internal"
)

func Now(ctx context.Context) time.Time {
	return internal.Now(ctx)
}

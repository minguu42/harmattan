package clock

import (
	"context"
	"time"

	"github.com/minguu42/harmattan/lib/clock/internal"
)

func Now(ctx context.Context) time.Time {
	return internal.Now(ctx)
}

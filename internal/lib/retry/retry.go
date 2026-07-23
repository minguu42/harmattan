package retry

import (
	"context"
	"errors"
	"time"

	"github.com/minguu42/harmattan/internal/lib/errtrace"
)

func Fixed(ctx context.Context, attempts int, delay time.Duration, f func() error) error {
	if attempts <= 0 {
		return errtrace.Wrap(errors.New("attempts must be greater than 0"))
	}
	if err := context.Cause(ctx); err != nil {
		return errtrace.Wrap(err)
	}

	for attempt := 1; ; attempt++ {
		err := f()
		if err == nil {
			return nil
		}

		if attempt == attempts {
			return errtrace.Wrap(err)
		}

		select {
		case <-time.After(delay):
		case <-ctx.Done():
			return errtrace.Wrap(errors.Join(err, context.Cause(ctx)))
		}
	}
}

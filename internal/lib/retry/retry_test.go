package retry_test

import (
	"context"
	"errors"
	"testing"
	"testing/synctest"
	"time"

	"github.com/minguu42/harmattan/internal/lib/retry"
	"github.com/stretchr/testify/assert"
)

func TestFixed(t *testing.T) {
	t.Parallel()

	t.Run("first_attempt_succeeds", func(t *testing.T) {
		t.Parallel()

		var calls int
		err := retry.Fixed(t.Context(), 3, time.Millisecond, func() error {
			calls++
			return nil
		})
		assert.NoError(t, err)
		assert.Equal(t, 1, calls)
	})
	t.Run("succeeds_after_retries", func(t *testing.T) {
		t.Parallel()

		var calls int
		err := retry.Fixed(t.Context(), 3, time.Millisecond, func() error {
			calls++
			if calls < 3 {
				return errors.New("some error")
			}
			return nil
		})
		assert.NoError(t, err)
		assert.Equal(t, 3, calls)
	})
	t.Run("all_attempts_fail", func(t *testing.T) {
		t.Parallel()

		wantErr := errors.New("some error")
		var calls int
		err := retry.Fixed(t.Context(), 3, time.Millisecond, func() error {
			calls++
			return wantErr
		})
		assert.ErrorIs(t, err, wantErr)
		assert.Equal(t, 3, calls)
	})
	t.Run("waits_between_attempts", func(t *testing.T) {
		t.Parallel()

		synctest.Test(t, func(t *testing.T) {
			start := time.Now()
			err := retry.Fixed(t.Context(), 3, time.Second, func() error {
				return errors.New("some error")
			})
			assert.Error(t, err)
			assert.Equal(t, 2*time.Second, time.Since(start))
		})
	})
	t.Run("zero_attempts", func(t *testing.T) {
		t.Parallel()

		var calls int
		err := retry.Fixed(t.Context(), 0, time.Millisecond, func() error {
			calls++
			return nil
		})
		assert.Error(t, err)
		assert.Equal(t, 0, calls)
	})
	t.Run("context_already_canceled", func(t *testing.T) {
		t.Parallel()

		ctx, cancel := context.WithCancel(t.Context())
		cancel()

		var calls int
		err := retry.Fixed(ctx, 3, time.Millisecond, func() error {
			calls++
			return nil
		})
		assert.ErrorIs(t, err, context.Canceled)
		assert.Equal(t, 0, calls)
	})
	t.Run("canceled_while_waiting", func(t *testing.T) {
		t.Parallel()

		wantErr := errors.New("some error")
		ctx, cancel := context.WithCancel(t.Context())

		var calls int
		err := retry.Fixed(ctx, 3, time.Hour, func() error {
			calls++
			cancel()
			return wantErr
		})
		assert.ErrorIs(t, err, wantErr)
		assert.ErrorIs(t, err, context.Canceled)
		assert.Equal(t, 1, calls)
	})
}

package clock_test

import (
	"testing"
	"time"

	"github.com/minguu42/harmattan/internal/lib/clock"
	"github.com/stretchr/testify/assert"
)

func TestNow(t *testing.T) {
	want := time.Now()
	got := clock.Now(t.Context())

	diff := got.Sub(want)
	if diff < 0 {
		diff = -diff
	}
	if diff > 10*time.Millisecond {
		t.Errorf("clock.Now() returned a time too far from want: got %s, want: %s", got, want)
	}
}

func TestWithFixedNow(t *testing.T) {
	want := time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)
	ctx := clock.WithFixedNow(t.Context(), want)

	got := clock.Now(ctx)
	assert.Equal(t, want, got)
}

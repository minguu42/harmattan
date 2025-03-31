package clock_test

import (
	"testing"
	"time"

	"github.com/minguu42/harmattan/lib/clock"
	"github.com/minguu42/harmattan/lib/clock/clocktest"
	"github.com/stretchr/testify/assert"
)

func TestNow(t *testing.T) {
	want := time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)
	ctx := clocktest.WithFixedNow(t, t.Context(), want)

	got := clock.Now(ctx)
	assert.Equal(t, want, got)
}

package clock_test

import (
	"testing"
	"time"

	"github.com/minguu42/harmattan/internal/lib/clock"
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

package pointers_test

import (
	"testing"
	"time"

	"github.com/minguu42/harmattan/internal/lib/pointers"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRef(t *testing.T) {
	values := []any{
		"Hello, World!",
		1,
		time.Date(2024, 2, 29, 12, 34, 56, 0, time.UTC),
	}
	for _, v := range values {
		p := pointers.Ref(v)
		require.NotNil(t, p)
		assert.Equal(t, v, *p)
	}
}

func TestOrZero(t *testing.T) {
	tests := []struct {
		v    *int
		want int
	}{
		{v: pointers.Ref(1), want: 1},
		{v: nil, want: 0},
	}
	for _, tt := range tests {
		got := pointers.OrZero(tt.v)
		assert.Equal(t, tt.want, got)
	}
}

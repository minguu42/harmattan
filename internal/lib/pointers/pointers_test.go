package pointers_test

import (
	"testing"
	"time"

	"github.com/minguu42/harmattan/internal/lib/pointers"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRef(t *testing.T) {
	tests := []struct {
		name string
		v    any
	}{
		{name: "string", v: "Hello, World!"},
		{name: "int", v: 42},
		{name: "time.Time", v: time.Date(2024, 2, 29, 12, 34, 56, 0, time.UTC)},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := pointers.Ref(tt.v)
			require.NotNil(t, p)

			assert.Equal(t, tt.v, *p)
		})
	}
}

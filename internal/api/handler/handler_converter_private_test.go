package handler

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTernary(t *testing.T) {
	tests := []struct {
		condition bool
		trueVal   any
		falseVal  any
		want      any
	}{
		{condition: true, trueVal: 1, falseVal: -1, want: 1},
		{condition: false, trueVal: 1, falseVal: -1, want: -1},
	}
	for _, tt := range tests {
		got := ternary(tt.condition, tt.trueVal, tt.falseVal)
		assert.Equal(t, tt.want, got)
	}
}

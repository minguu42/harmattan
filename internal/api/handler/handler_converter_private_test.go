package handler

import (
	"testing"
	"time"

	"github.com/minguu42/harmattan/internal/api/openapi"
	"github.com/minguu42/harmattan/internal/lib/pointers"
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

func TestConvertOptDateTime(t *testing.T) {
	tests := []struct {
		t    *time.Time
		want openapi.OptDateTime
	}{
		{t: nil, want: openapi.OptDateTime{}},
		{
			t: pointers.Ref(time.Date(2025, 1, 2, 15, 4, 5, 0, time.Local)),
			want: openapi.OptDateTime{
				Value: time.Date(2025, 1, 2, 15, 4, 5, 0, time.Local),
				Set:   true,
			},
		},
	}
	for _, tt := range tests {
		got := convertOptDateTime(tt.t)
		assert.Equal(t, tt.want, got)
	}
}

func TestConvertOptDate(t *testing.T) {
	tests := []struct {
		t    *time.Time
		want openapi.OptDate
	}{
		{t: nil, want: openapi.OptDate{}},
		{
			t: pointers.Ref(time.Date(2025, 1, 2, 0, 0, 0, 0, time.Local)),
			want: openapi.OptDate{
				Value: time.Date(2025, 1, 2, 0, 0, 0, 0, time.Local),
				Set:   true,
			},
		},
	}
	for _, tt := range tests {
		got := convertOptDate(tt.t)
		assert.Equal(t, tt.want, got)
	}
}

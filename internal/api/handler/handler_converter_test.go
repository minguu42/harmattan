package handler_test

import (
	"testing"
	"time"

	"github.com/minguu42/harmattan/internal/api/handler"
	"github.com/minguu42/harmattan/internal/api/openapi"
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
		got := handler.Ternary(tt.condition, tt.trueVal, tt.falseVal)
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
			t: new(time.Date(2025, 1, 2, 15, 4, 5, 0, time.Local)),
			want: openapi.OptDateTime{
				Value: time.Date(2025, 1, 2, 15, 4, 5, 0, time.Local),
				Set:   true,
			},
		},
	}
	for _, tt := range tests {
		got := handler.ConvertOptDateTime(tt.t)
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
			t: new(time.Date(2025, 1, 2, 0, 0, 0, 0, time.Local)),
			want: openapi.OptDate{
				Value: time.Date(2025, 1, 2, 0, 0, 0, 0, time.Local),
				Set:   true,
			},
		},
	}
	for _, tt := range tests {
		got := handler.ConvertOptDate(tt.t)
		assert.Equal(t, tt.want, got)
	}
}

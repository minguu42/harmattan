package handler_test

import (
	"testing"
	"time"

	"github.com/minguu42/harmattan/internal/api/handler"
	"github.com/minguu42/harmattan/internal/api/openapi"
	"github.com/minguu42/harmattan/internal/lib/httpcheck"
	"github.com/stretchr/testify/assert"
)

func TestHandler_NotFound(t *testing.T) {
	httpcheck.New(th).Test(t, "GET", "/non-existent-path").
		Check().
		HasStatus(404).
		HasJSON(handler.ErrorResponse{Code: 404, Message: "指定したパスは見つかりません"})
}

func TestHandler_MethodNotFound(t *testing.T) {
	t.Run("method_not_allowed", func(t *testing.T) {
		httpcheck.New(th).Test(t, "POST", "/health").
			Check().
			HasStatus(405).
			HasJSON(handler.ErrorResponse{Code: 405, Message: "指定したメソッドは許可されていません"})
	})
	t.Run("options", func(t *testing.T) {
		httpcheck.New(th).Test(t, "OPTIONS", "/health").
			Check().
			HasStatus(204).
			HasHeader("Access-Control-Allow-Methods", "GET").
			HasHeader("Access-Control-Allow-Headers", "Content-Type")
	})
}

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

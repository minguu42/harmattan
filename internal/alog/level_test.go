package alog_test

import (
	"log/slog"
	"testing"

	"github.com/minguu42/harmattan/internal/alog"
	"github.com/stretchr/testify/assert"
)

func TestLevel_Level(t *testing.T) {
	t.Parallel()

	tests := []struct {
		level alog.Level
		want  slog.Level
	}{
		{level: alog.LevelDebug, want: slog.LevelDebug},
		{level: alog.LevelInfo, want: slog.LevelInfo},
		{level: alog.LevelWarn, want: slog.LevelWarn},
		{level: alog.LevelError, want: slog.LevelError},
		{level: alog.Level("unknown"), want: slog.LevelInfo},
		{level: alog.Level(""), want: slog.LevelInfo},
	}
	for _, tt := range tests {
		got := tt.level.Level()
		assert.Equal(t, tt.want, got)
	}
}

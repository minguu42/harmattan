package alog_test

import (
	"log/slog"
	"testing"

	"github.com/minguu42/harmattan/internal/alog"
	"github.com/stretchr/testify/assert"
)

func TestParseLevel(t *testing.T) {
	t.Parallel()

	tests := []struct {
		s    string
		want alog.Level
	}{
		{s: "debug", want: alog.LevelDebug},
		{s: "info", want: alog.LevelInfo},
		{s: "warn", want: alog.LevelWarn},
		{s: "error", want: alog.LevelError},
		{s: "unknown", want: alog.LevelUnknown},
		{s: "", want: alog.LevelUnknown},
		{s: "INFO", want: alog.LevelUnknown},
	}
	for _, tt := range tests {
		got := alog.ParseLevel(tt.s)
		assert.Equal(t, tt.want, got)
	}
}

func TestLevel_Level(t *testing.T) {
	t.Parallel()

	tests := []struct {
		level alog.Level
		want  slog.Level
	}{
		{level: alog.LevelUnknown, want: slog.LevelInfo},
		{level: alog.LevelDebug, want: slog.LevelDebug},
		{level: alog.LevelInfo, want: slog.LevelInfo},
		{level: alog.LevelWarn, want: slog.LevelWarn},
		{level: alog.LevelError, want: slog.LevelError},
	}
	for _, tt := range tests {
		got := tt.level.Level()
		assert.Equal(t, tt.want, got)
	}
}

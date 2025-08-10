package logutil_test

import (
	"bytes"
	"encoding/json"
	"log/slog"
	"testing"

	"github.com/minguu42/harmattan/internal/lib/logutil"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMaskAttr(t *testing.T) {
	t.Parallel()

	type Foo struct {
		F1 string `json:"f1"`
		F2 string `json:"f2" log:"mask"`
	}
	type Bar struct {
		Foo Foo `json:"foo"`
	}
	type Baz struct {
		Foo *Foo `json:"foo"`
	}
	tests := []struct {
		name  string
		value any
		want  map[string]any
	}{
		{
			name:  "mask",
			value: Foo{F1: "v1", F2: "v2"},
			want:  map[string]any{"f1": "v1", "f2": "<hidden>"},
		},
		{
			name:  "mask-pointer",
			value: &Foo{F1: "v1", F2: "v2"},
			want:  map[string]any{"f1": "v1", "f2": "<hidden>"},
		},
		{
			name:  "mask-deep",
			value: Bar{Foo: Foo{F1: "v1", F2: "v2"}},
			want:  map[string]any{"foo": map[string]any{"f1": "v1", "f2": "<hidden>"}},
		},
		{
			name:  "mask-deep-pointer",
			value: Baz{Foo: &Foo{F1: "v1", F2: "v2"}},
			want:  map[string]any{"foo": map[string]any{"f1": "v1", "f2": "<hidden>"}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			var buf bytes.Buffer
			l := slog.New(slog.NewJSONHandler(&buf, &slog.HandlerOptions{ReplaceAttr: func(_ []string, a slog.Attr) slog.Attr {
				return logutil.MaskAttr(a)
			}}))
			l.Info("", slog.Any("test", tt.value))

			var m map[string]any
			require.NoError(t, json.Unmarshal(buf.Bytes(), &m))
			got, ok := m["test"].(map[string]any)
			require.True(t, ok)

			assert.Equal(t, tt.want, got)
		})
	}
}

package atel_test

import (
	"bytes"
	"encoding/json"
	"log/slog"
	"testing"

	"github.com/minguu42/harmattan/internal/atel"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMaskAttr(t *testing.T) {
	t.Parallel()

	type AllFieldsExported struct {
		Unmask  string            `json:"f1"`
		Bool    bool              `json:"f2" log:"mask"`
		Int     int               `json:"f3" log:"mask"`
		Uint    uint              `json:"f4" log:"mask"`
		Float   float64           `json:"f5" log:"mask"`
		Map     map[string]string `json:"f6" log:"mask"`
		Slice   []string          `json:"f7" log:"mask"`
		String  string            `json:"f8" log:"mask"`
		Pointer *string           `json:"f9" log:"mask"`
	}
	type StructWrapper struct {
		F AllFieldsExported `json:"f"`
	}
	type StructPointerWrapper struct {
		F *AllFieldsExported `json:"f"`
	}
	type WithUnexported struct {
		Public  string `json:"public" log:"mask"`
		private string
	}
	type StructWrapperWithUnexported struct {
		F WithUnexported `json:"f"`
	}
	type StructPointerWrapperWithUnexported struct {
		F *WithUnexported `json:"f"`
	}
	tests := []struct {
		name  string
		value any
		want  map[string]any
	}{
		{
			name: "mask",
			value: AllFieldsExported{
				Unmask:  "foo",
				Bool:    true,
				Int:     1,
				Uint:    10,
				Float:   3.14,
				Map:     map[string]string{"one": "1", "two": "2"},
				Slice:   []string{"1", "2"},
				String:  "bar",
				Pointer: new("baz"),
			},
			want: map[string]any{
				"f1": "foo",
				"f2": false,
				"f3": 0.0,
				"f4": 0.0,
				"f5": 0.0,
				"f6": map[string]any{},
				"f7": []any{},
				"f8": "<hidden>",
				"f9": nil,
			},
		},
		{
			name: "mask-pointer",
			value: &AllFieldsExported{
				Unmask:  "foo",
				Bool:    true,
				Int:     1,
				Uint:    10,
				Float:   3.14,
				Map:     map[string]string{"one": "1"},
				Slice:   []string{"1", "2"},
				String:  "bar",
				Pointer: new("baz"),
			},
			want: map[string]any{
				"f1": "foo",
				"f2": false,
				"f3": 0.0,
				"f4": 0.0,
				"f5": 0.0,
				"f6": map[string]any{},
				"f7": []any{},
				"f8": "<hidden>",
				"f9": nil,
			},
		},
		{
			name: "mask-deep",
			value: StructWrapper{F: AllFieldsExported{
				Unmask:  "foo",
				Bool:    true,
				Int:     1,
				Uint:    10,
				Float:   3.14,
				Map:     map[string]string{"one": "1", "two": "2"},
				Slice:   []string{"1", "2"},
				String:  "bar",
				Pointer: new("baz"),
			}},
			want: map[string]any{"f": map[string]any{
				"f1": "foo",
				"f2": false,
				"f3": 0.0,
				"f4": 0.0,
				"f5": 0.0,
				"f6": map[string]any{},
				"f7": []any{},
				"f8": "<hidden>",
				"f9": nil,
			}},
		},
		{
			name: "mask-deep-pointer",
			value: StructPointerWrapper{F: &AllFieldsExported{
				Unmask:  "foo",
				Bool:    true,
				Int:     1,
				Uint:    10,
				Float:   3.14,
				Map:     map[string]string{"one": "1", "two": "2"},
				Slice:   []string{"1", "2"},
				String:  "bar",
				Pointer: new("baz"),
			}},
			want: map[string]any{"f": map[string]any{
				"f1": "foo",
				"f2": false,
				"f3": 0.0,
				"f4": 0.0,
				"f5": 0.0,
				"f6": map[string]any{},
				"f7": []any{},
				"f8": "<hidden>",
				"f9": nil,
			}},
		},
		{
			name:  "contains-unexported-field",
			value: WithUnexported{Public: "secret", private: "ignore"},
			want:  map[string]any{"public": "secret"},
		},
		{
			name:  "contains-unexported-field-pointer",
			value: &WithUnexported{Public: "secret", private: "ignore"},
			want:  map[string]any{"public": "secret"},
		},
		{
			name:  "contains-unexported-field-deep",
			value: StructWrapperWithUnexported{F: WithUnexported{Public: "secret", private: "ignore"}},
			want:  map[string]any{"f": map[string]any{"public": "secret"}},
		},
		{
			name:  "contains-unexported-field-deep-pointer",
			value: StructPointerWrapperWithUnexported{F: &WithUnexported{Public: "secret", private: "ignore"}},
			want:  map[string]any{"f": map[string]any{"public": "secret"}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			var buf bytes.Buffer
			l := slog.New(slog.NewJSONHandler(&buf, &slog.HandlerOptions{ReplaceAttr: func(_ []string, a slog.Attr) slog.Attr {
				return atel.MaskAttr(a)
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

package atel

import (
	"bytes"
	"encoding/json"
	"log/slog"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMaskAttr(t *testing.T) {
	t.Parallel()

	type Unmasked struct {
		Bool    bool              `json:"bool" log:"allow"`
		Int     int               `json:"int" log:"allow"`
		Uint    uint              `json:"uint" log:"allow"`
		Float   float64           `json:"float" log:"allow"`
		Map     map[string]string `json:"map" log:"allow"`
		Slice   []string          `json:"slice" log:"allow"`
		String  string            `json:"string" log:"allow"`
		Pointer *string           `json:"pointer" log:"allow"`
	}
	type Masked struct {
		Bool    bool              `json:"bool"`
		Int     int               `json:"int"`
		Uint    uint              `json:"uint"`
		Float   float64           `json:"float"`
		Map     map[string]string `json:"map"`
		Slice   []string          `json:"slice"`
		String  string            `json:"string"`
		Pointer *string           `json:"pointer"`
	}
	type WithUnexported struct {
		Public  string `json:"public"`
		private string
	}
	type UnmaskedWrapper[T any] struct {
		Struct T `json:"struct" log:"allow"`
	}
	type UnmaskedPointerWrapper[T any] struct {
		Struct *T `json:"struct" log:"allow"`
	}
	type MaskedWrapper[T any] struct {
		Struct T `json:"struct"`
	}
	tests := []struct {
		name  string
		value any
		want  map[string]any
	}{
		{
			name: "allow_all_fields",
			value: Unmasked{
				Bool:    true,
				Int:     1,
				Uint:    10,
				Float:   3.14,
				Map:     map[string]string{"one": "1", "two": "2"},
				Slice:   []string{"1", "2"},
				String:  "foo",
				Pointer: new("bar"),
			},
			want: map[string]any{
				"bool":    true,
				"int":     1.0,
				"uint":    10.0,
				"float":   3.14,
				"map":     map[string]any{"one": "1", "two": "2"},
				"slice":   []any{"1", "2"},
				"string":  "foo",
				"pointer": "bar",
			},
		},
		{
			name: "allow_all_fields_pointer",
			value: &Unmasked{
				Bool:    true,
				Int:     1,
				Uint:    10,
				Float:   3.14,
				Map:     map[string]string{"one": "1", "two": "2"},
				Slice:   []string{"1", "2"},
				String:  "foo",
				Pointer: new("bar"),
			},
			want: map[string]any{
				"bool":    true,
				"int":     1.0,
				"uint":    10.0,
				"float":   3.14,
				"map":     map[string]any{"one": "1", "two": "2"},
				"slice":   []any{"1", "2"},
				"string":  "foo",
				"pointer": "bar",
			},
		},
		{
			name: "allow_nested_struct",
			value: UnmaskedWrapper[Unmasked]{Struct: Unmasked{
				Bool:    true,
				Int:     1,
				Uint:    10,
				Float:   3.14,
				Map:     map[string]string{"one": "1", "two": "2"},
				Slice:   []string{"1", "2"},
				String:  "foo",
				Pointer: new("bar"),
			}},
			want: map[string]any{"struct": map[string]any{
				"bool":    true,
				"int":     1.0,
				"uint":    10.0,
				"float":   3.14,
				"map":     map[string]any{"one": "1", "two": "2"},
				"slice":   []any{"1", "2"},
				"string":  "foo",
				"pointer": "bar",
			}},
		},
		{
			name: "allow_nested_struct_pointer",
			value: UnmaskedPointerWrapper[Unmasked]{Struct: &Unmasked{
				Bool:    true,
				Int:     1,
				Uint:    10,
				Float:   3.14,
				Map:     map[string]string{"one": "1", "two": "2"},
				Slice:   []string{"1", "2"},
				String:  "foo",
				Pointer: new("bar"),
			}},
			want: map[string]any{"struct": map[string]any{
				"bool":    true,
				"int":     1.0,
				"uint":    10.0,
				"float":   3.14,
				"map":     map[string]any{"one": "1", "two": "2"},
				"slice":   []any{"1", "2"},
				"string":  "foo",
				"pointer": "bar",
			}},
		},
		{
			name:  "allow_nil_nested_struct_pointer",
			value: UnmaskedPointerWrapper[Unmasked]{Struct: nil},
			want:  map[string]any{"struct": nil},
		},
		{
			name: "mask_struct",
			value: Masked{
				Bool:    true,
				Int:     1,
				Uint:    10,
				Float:   3.14,
				Map:     map[string]string{"one": "1", "two": "2"},
				Slice:   []string{"1", "2"},
				String:  "foo",
				Pointer: new("bar"),
			},
			want: map[string]any{
				"bool":    false,
				"int":     0.0,
				"uint":    0.0,
				"float":   0.0,
				"map":     map[string]any{},
				"slice":   []any{},
				"string":  "<hidden>",
				"pointer": nil,
			},
		},
		{
			name: "mask_struct_pointer",
			value: &Masked{
				Bool:    true,
				Int:     1,
				Uint:    10,
				Float:   3.14,
				Map:     map[string]string{"one": "1", "two": "2"},
				Slice:   []string{"1", "2"},
				String:  "foo",
				Pointer: new("bar"),
			},
			want: map[string]any{
				"bool":    false,
				"int":     0.0,
				"uint":    0.0,
				"float":   0.0,
				"map":     map[string]any{},
				"slice":   []any{},
				"string":  "<hidden>",
				"pointer": nil,
			},
		},
		{
			name: "mask_deep_struct",
			value: UnmaskedWrapper[Masked]{Struct: Masked{
				Bool:    true,
				Int:     1,
				Uint:    10,
				Float:   3.14,
				Map:     map[string]string{"one": "1", "two": "2"},
				Slice:   []string{"1", "2"},
				String:  "foo",
				Pointer: new("bar"),
			}},
			want: map[string]any{"struct": map[string]any{
				"bool":    false,
				"int":     0.0,
				"uint":    0.0,
				"float":   0.0,
				"map":     map[string]any{},
				"slice":   []any{},
				"string":  "<hidden>",
				"pointer": nil,
			}},
		},
		{
			name: "mask_deep_struct_pointer",
			value: UnmaskedPointerWrapper[Masked]{Struct: &Masked{
				Bool:    true,
				Int:     1,
				Uint:    10,
				Float:   3.14,
				Map:     map[string]string{"one": "1", "two": "2"},
				Slice:   []string{"1", "2"},
				String:  "foo",
				Pointer: new("bar"),
			}},
			want: map[string]any{"struct": map[string]any{
				"bool":    false,
				"int":     0.0,
				"uint":    0.0,
				"float":   0.0,
				"map":     map[string]any{},
				"slice":   []any{},
				"string":  "<hidden>",
				"pointer": nil,
			}},
		},
		{
			name: "mask_nested_without_parent_allow",
			value: MaskedWrapper[Unmasked]{Struct: Unmasked{
				Bool:    true,
				Int:     1,
				Uint:    10,
				Float:   3.14,
				Map:     map[string]string{"one": "1", "two": "2"},
				Slice:   []string{"1", "2"},
				String:  "foo",
				Pointer: new("bar"),
			}},
			want: map[string]any{"struct": map[string]any{
				"bool":    false,
				"int":     0.0,
				"uint":    0.0,
				"float":   0.0,
				"map":     map[string]any{},
				"slice":   []any{},
				"string":  "<hidden>",
				"pointer": nil,
			}},
		},
		{
			name:  "skip_unexported_field",
			value: WithUnexported{Public: "secret", private: "ignore"},
			want:  map[string]any{"public": "<hidden>"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			var buf bytes.Buffer
			l := slog.New(slog.NewJSONHandler(&buf, &slog.HandlerOptions{ReplaceAttr: func(_ []string, a slog.Attr) slog.Attr {
				return maskAttr(a)
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

func TestMaskAttr_passthrough(t *testing.T) {
	t.Parallel()

	type Foo struct {
		Foo string
	}
	tests := []struct {
		name string
		attr slog.Attr
	}{
		{
			name: "not_kind_any",
			attr: slog.String("foo", "bar"),
		},
		{
			name: "nil_struct_pointer",
			attr: slog.Any("foo", (*Foo)(nil)),
		},
		{
			name: "non_struct_pointer",
			attr: slog.Any("foo", new(42)),
		},
		{
			name: "non_struct_value",
			attr: slog.Any("foo", "bar"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got := maskAttr(tt.attr)
			assert.Equal(t, tt.attr, got)
		})
	}
}

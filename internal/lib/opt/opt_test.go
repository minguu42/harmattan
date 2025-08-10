package opt_test

import (
	"testing"

	"github.com/minguu42/harmattan/internal/lib/opt"
	"github.com/minguu42/harmattan/internal/lib/pointers"
	"github.com/stretchr/testify/assert"
)

func TestFromPointer(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name string
		p    any
		want any
	}{
		{
			name: "nil string pointer",
			p:    (*string)(nil),
			want: opt.Option[string]{},
		},
		{
			name: "valid string pointer",
			p:    pointers.Ref("Hello, World!"),
			want: opt.Option[string]{V: "Hello, World!", Valid: true},
		},
		{
			name: "nil int pointer",
			p:    (*int)(nil),
			want: opt.Option[int]{},
		},
		{
			name: "valid int pointer",
			p:    pointers.Ref(42),
			want: opt.Option[int]{V: 42, Valid: true},
		},
		{
			name: "nil bool pointer",
			p:    (*bool)(nil),
			want: opt.Option[bool]{},
		},
		{
			name: "valid bool pointer",
			p:    pointers.Ref(true),
			want: opt.Option[bool]{V: true, Valid: true},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			switch ptr := tt.p.(type) {
			case *string:
				got := opt.FromPointer(ptr)
				assert.Equal(t, tt.want, got)
			case *int:
				got := opt.FromPointer(ptr)
				assert.Equal(t, tt.want, got)
			case *bool:
				got := opt.FromPointer(ptr)
				assert.Equal(t, tt.want, got)
			}
		})
	}
}

func TestOpt_ToPointer(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name string
		o    any
		want any
	}{
		{
			name: "valid string option",
			o:    opt.Option[string]{V: "Hello, World!", Valid: true},
			want: pointers.Ref("Hello, World!"),
		},
		{
			name: "invalid string option",
			o:    opt.Option[string]{},
			want: (*string)(nil),
		},
		{
			name: "valid int option",
			o:    opt.Option[int]{V: 42, Valid: true},
			want: pointers.Ref(42),
		},
		{
			name: "invalid int option",
			o:    opt.Option[int]{},
			want: (*int)(nil),
		},
		{
			name: "valid bool option",
			o:    opt.Option[bool]{V: true, Valid: true},
			want: pointers.Ref(true),
		},
		{
			name: "invalid bool option",
			o:    opt.Option[bool]{},
			want: (*bool)(nil),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			switch o := tt.o.(type) {
			case opt.Option[string]:
				got := o.ToPointer()
				assert.Equal(t, tt.want, got)
			case opt.Option[int]:
				got := o.ToPointer()
				assert.Equal(t, tt.want, got)
			case opt.Option[bool]:
				got := o.ToPointer()
				assert.Equal(t, tt.want, got)
			}
		})
	}
}

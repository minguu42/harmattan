package handler_test

import (
	"strings"
	"testing"

	"github.com/minguu42/harmattan/internal/api/handler"
	"github.com/stretchr/testify/assert"
)

func TestValidateTagName(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		tagName string
		want    []error
	}{
		{name: "empty", tagName: "", want: []error{handler.ErrTagNameLength}},
		{name: "min_length_boundary", tagName: "a"},
		{name: "max_length_boundary", tagName: strings.Repeat("a", 20)},
		{name: "max_length_boundary_multibyte", tagName: strings.Repeat("あ", 20)},
		{name: "above_max_length", tagName: strings.Repeat("a", 21), want: []error{handler.ErrTagNameLength}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			assert.ElementsMatch(t, tt.want, handler.ValidateTagName(tt.tagName))
		})
	}
}

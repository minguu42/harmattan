package handler_test

import (
	"strings"
	"testing"

	"github.com/minguu42/harmattan/internal/api/handler"
	"github.com/stretchr/testify/assert"
)

func TestValidateStepName(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		stepName string
		want     []error
	}{
		{name: "empty", stepName: "", want: []error{handler.ErrStepNameLength}},
		{name: "min_length_boundary", stepName: "a"},
		{name: "max_length_boundary", stepName: strings.Repeat("a", 100)},
		{name: "max_length_boundary_multibyte", stepName: strings.Repeat("あ", 100)},
		{name: "above_max_length", stepName: strings.Repeat("a", 101), want: []error{handler.ErrStepNameLength}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			assert.ElementsMatch(t, tt.want, handler.ValidateStepName(tt.stepName))
		})
	}
}

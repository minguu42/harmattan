package handler_test

import (
	"strings"
	"testing"

	"github.com/minguu42/harmattan/internal/api/handler"
	"github.com/stretchr/testify/assert"
)

func TestValidateProjectName(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name        string
		projectName string
		want        []error
	}{
		{name: "empty", projectName: "", want: []error{handler.ErrProjectNameLength}},
		{name: "min_length_boundary", projectName: "a"},
		{name: "max_length_boundary", projectName: strings.Repeat("a", 80)},
		{name: "max_length_boundary_multibyte", projectName: strings.Repeat("あ", 80)},
		{name: "above_max_length", projectName: strings.Repeat("a", 81), want: []error{handler.ErrProjectNameLength}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			assert.ElementsMatch(t, tt.want, handler.ValidateProjectName(tt.projectName))
		})
	}
}

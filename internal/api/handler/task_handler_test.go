package handler_test

import (
	"strings"
	"testing"

	"github.com/minguu42/harmattan/internal/api/handler"
	"github.com/stretchr/testify/assert"
)

func TestValidateTaskName(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		taskName string
		want     []error
	}{
		{name: "empty", taskName: "", want: []error{handler.ErrTaskNameLength}},
		{name: "min_length_boundary", taskName: "a"},
		{name: "max_length_boundary", taskName: strings.Repeat("a", 100)},
		{name: "max_length_boundary_multibyte", taskName: strings.Repeat("あ", 100)},
		{name: "above_max_length", taskName: strings.Repeat("a", 101), want: []error{handler.ErrTaskNameLength}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			assert.ElementsMatch(t, tt.want, handler.ValidateTaskName(tt.taskName))
		})
	}
}

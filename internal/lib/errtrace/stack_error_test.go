package errtrace_test

import (
	"errors"
	"fmt"
	"log/slog"
	"testing"

	"github.com/minguu42/harmattan/internal/lib/errtrace"
	"github.com/stretchr/testify/assert"
)

func TestStackError_Format(t *testing.T) {
	t.Parallel()

	err := errtrace.Wrap(errors.New("test error"),
		slog.String("k1", "v1"),
		slog.String("k2", "v2"),
	)

	assert.Equal(t, "test error", fmt.Sprintf("%s", err))
	assert.Equal(t, "test error", fmt.Sprintf("%v", err))
	output := fmt.Sprintf("%+v", err)
	assert.Contains(t, output, "test error")
	assert.Contains(t, output, "stack_error_test.go:")
	assert.Contains(t, output, "[k1=v1 k2=v2]")
}

func TestStackError_Unwrap(t *testing.T) {
	t.Parallel()

	err := errors.New("test error")
	wrapped := errtrace.Wrap(err)
	assert.ErrorIs(t, wrapped, err)
}

func TestStackError_Frames(t *testing.T) {
	t.Parallel()

	err := errtrace.Wrap(errors.New("test error"))

	var stackError *errtrace.StackError
	assert.ErrorAs(t, err, &stackError)
	assert.NotEqual(t, 0, len(stackError.Frames()))
}

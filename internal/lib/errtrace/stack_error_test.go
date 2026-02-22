package errtrace_test

import (
	"errors"
	"fmt"
	"testing"

	"github.com/minguu42/harmattan/internal/lib/errtrace"
	"github.com/stretchr/testify/assert"
)

func TestStackError_Unwrap(t *testing.T) {
	t.Parallel()

	err := errors.New("test error")
	wrapped := errtrace.Wrap(err)
	assert.ErrorIs(t, wrapped, err)
}

func TestStackError_Format(t *testing.T) {
	t.Parallel()

	err := errtrace.Wrap(errors.New("test error"))

	plusVOutput := fmt.Sprintf("%+v", err)
	assert.Contains(t, plusVOutput, "test error")
	assert.Contains(t, plusVOutput, "stack_error_test.go")

	vOutput := fmt.Sprintf("%v", err)
	assert.Equal(t, "test error", vOutput)

	sOutput := fmt.Sprintf("%s", err)
	assert.Equal(t, "test error", sOutput)
}

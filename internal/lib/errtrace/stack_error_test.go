package errtrace_test

import (
	"encoding/json"
	"errors"
	"testing"

	"github.com/minguu42/harmattan/internal/lib/errtrace"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestStackError_Unwrap(t *testing.T) {
	t.Parallel()

	err := errors.New("test error")
	wrapped := errtrace.Wrap(err)
	assert.ErrorIs(t, wrapped, err)
}

func TestStackError_MarshalJSON(t *testing.T) {
	t.Parallel()

	got, err := json.Marshal(errtrace.Wrap(errors.New("test error")))
	require.NoError(t, err)

	var result map[string]any
	require.NoError(t, json.Unmarshal(got, &result))
	assert.Equal(t, "test error", result["message"])
	assert.NotEmpty(t, result["frames"])
}

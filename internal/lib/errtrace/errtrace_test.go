package errtrace_test

import (
	"encoding/json"
	"errors"
	"testing"

	"github.com/minguu42/harmattan/internal/lib/errtrace"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestWrap(t *testing.T) {
	t.Parallel()

	t.Run("nil", func(t *testing.T) {
		t.Parallel()

		got := errtrace.Wrap(nil)
		assert.Nil(t, got)
	})
	t.Run("wrap", func(t *testing.T) {
		t.Parallel()

		got := errtrace.Wrap(errors.New("test error"))
		var serr *errtrace.StackError
		assert.ErrorAs(t, got, &serr)
		assert.Equal(t, "test error", got.Error())
	})
	t.Run("no double wrap", func(t *testing.T) {
		t.Parallel()

		wrapped1 := errtrace.Wrap(errors.New("test error"))
		wrapped2 := errtrace.Wrap(wrapped1)
		assert.Same(t, wrapped1, wrapped2)
	})
	t.Run("implement unwrap interface", func(t *testing.T) {
		t.Parallel()

		err := errors.New("original")
		wrapped := errtrace.Wrap(err)
		assert.ErrorIs(t, wrapped, err)
		assert.Equal(t, err, errors.Unwrap(wrapped))
	})
	t.Run("have stack trace", func(t *testing.T) {
		t.Parallel()

		got, err := json.Marshal(errtrace.Wrap(errors.New("test error")))
		require.NoError(t, err)

		var result map[string]any
		require.NoError(t, json.Unmarshal(got, &result))
		assert.Equal(t, "test error", result["message"])
		assert.NotEmpty(t, result["frames"])
	})
}

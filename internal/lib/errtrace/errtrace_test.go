package errtrace_test

import (
	"errors"
	"testing"

	"github.com/minguu42/harmattan/internal/lib/errtrace"
	"github.com/stretchr/testify/assert"
)

func TestFromStack(t *testing.T) {
	t.Parallel()

	t.Run("nil", func(t *testing.T) {
		t.Parallel()

		got := errtrace.FromStack(nil, []uintptr{1, 2, 3})
		assert.Nil(t, got)
	})
	t.Run("create from stack", func(t *testing.T) {
		t.Parallel()

		stack := []uintptr{1, 2, 3}
		got := errtrace.FromStack(errors.New("test error"), stack)
		var serr *errtrace.StackError
		assert.ErrorAs(t, got, &serr)
		assert.Equal(t, "test error", got.Error())
	})
	t.Run("no double wrap", func(t *testing.T) {
		t.Parallel()

		wrapped1 := errtrace.FromStack(errors.New("test error"), []uintptr{1, 2, 3})
		wrapped2 := errtrace.FromStack(wrapped1, []uintptr{4, 5, 6})
		assert.Same(t, wrapped1, wrapped2)
	})
}

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
}

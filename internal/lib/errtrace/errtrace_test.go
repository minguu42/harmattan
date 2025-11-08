package errtrace_test

import (
	"errors"
	"testing"

	"github.com/minguu42/harmattan/internal/lib/errtrace"
	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	t.Parallel()

	t.Run("nil", func(t *testing.T) {
		t.Parallel()

		got := errtrace.New(nil, []uintptr{1, 2, 3})
		assert.Nil(t, got)
	})
	t.Run("create with stack", func(t *testing.T) {
		t.Parallel()

		stack := []uintptr{1, 2, 3}
		got := errtrace.New(errors.New("test error"), stack)
		var serr *errtrace.StackError
		assert.ErrorAs(t, got, &serr)
		assert.Equal(t, "test error", got.Error())
	})
	t.Run("no double wrap", func(t *testing.T) {
		t.Parallel()

		wrapped1 := errtrace.New(errors.New("test error"), []uintptr{1, 2, 3})
		wrapped2 := errtrace.New(wrapped1, []uintptr{4, 5, 6})
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

package errtrace_test

import (
	"errors"
	"log/slog"
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
	t.Run("with attrs", func(t *testing.T) {
		t.Parallel()

		stack := []uintptr{1, 2, 3}
		got := errtrace.FromStack(errors.New("test error"), stack, slog.String("k1", "v1"))
		var serr *errtrace.StackError
		assert.ErrorAs(t, got, &serr)
		assert.Len(t, serr.Attrs(), 1)
		assert.Equal(t, "k1", serr.Attrs()[0].Key)
		assert.Equal(t, "v1", serr.Attrs()[0].Value.String())
	})
	t.Run("add attrs to existing stack error", func(t *testing.T) {
		t.Parallel()

		wrapped := errtrace.FromStack(errors.New("test error"), []uintptr{1, 2, 3}, slog.String("k1", "v1"))
		got := errtrace.FromStack(wrapped, []uintptr{4, 5, 6}, slog.String("k2", "v2"))
		assert.Same(t, wrapped, got)
		var serr *errtrace.StackError
		assert.ErrorAs(t, got, &serr)
		assert.Len(t, serr.Attrs(), 2)
		assert.Equal(t, "k1", serr.Attrs()[0].Key)
		assert.Equal(t, "k2", serr.Attrs()[1].Key)
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
	t.Run("with attrs", func(t *testing.T) {
		t.Parallel()

		got := errtrace.Wrap(errors.New("test error"), slog.String("k1", "v1"))
		var serr *errtrace.StackError
		assert.ErrorAs(t, got, &serr)
		assert.Len(t, serr.Attrs(), 1)
		assert.Equal(t, "k1", serr.Attrs()[0].Key)
		assert.Equal(t, "v1", serr.Attrs()[0].Value.String())
	})
	t.Run("add attrs to existing stack error", func(t *testing.T) {
		t.Parallel()

		wrapped := errtrace.Wrap(errors.New("test error"), slog.String("k1", "v1"))
		got := errtrace.Wrap(wrapped, slog.String("k2", "v2"))
		assert.Same(t, wrapped, got)
		var serr *errtrace.StackError
		assert.ErrorAs(t, got, &serr)
		assert.Len(t, serr.Attrs(), 2)
		assert.Equal(t, "k1", serr.Attrs()[0].Key)
		assert.Equal(t, "k2", serr.Attrs()[1].Key)
	})
}

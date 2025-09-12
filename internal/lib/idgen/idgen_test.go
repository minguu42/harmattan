package idgen_test

import (
	"testing"

	"github.com/minguu42/harmattan/internal/lib/idgen"
	"github.com/stretchr/testify/assert"
)

func TestULID(t *testing.T) {
	ctx := t.Context()

	id := idgen.ULID(ctx)
	assert.NotEmpty(t, id)
	assert.Len(t, id, 26)
	assert.NotEqual(t, id, idgen.ULID(ctx))
}

func TestWithFixedULID(t *testing.T) {
	want := "01G65Z755AFWAKHE12NY0CQ9FH"
	ctx := idgen.WithFixedULID(t.Context(), want)

	got := idgen.ULID(ctx)
	assert.Equal(t, want, got)
}

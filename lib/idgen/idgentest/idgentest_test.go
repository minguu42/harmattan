package idgentest_test

import (
	"testing"

	"github.com/minguu42/harmattan/lib/idgen"
	"github.com/minguu42/harmattan/lib/idgen/idgentest"
	"github.com/stretchr/testify/assert"
)

func TestWithFixedULID(t *testing.T) {
	want := "01JGFJJZ000000000000000000"
	ctx := idgentest.WithFixedULID(t.Context(), want)

	got := idgen.ULID(ctx)
	assert.Equal(t, want, got)
}

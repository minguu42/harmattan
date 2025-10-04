package ptr_test

import (
	"testing"
	"time"

	"github.com/minguu42/harmattan/internal/lib/ptr"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRef(t *testing.T) {
	values := []any{
		"Hello, World!",
		1,
		time.Date(2024, 2, 29, 12, 34, 56, 0, time.UTC),
	}
	for _, v := range values {
		p := ptr.Ref(v)
		require.NotNil(t, p)
		assert.Equal(t, v, *p)
	}
}

package idgen

import (
	"context"

	"github.com/minguu42/harmattan/lib/idgen/internal"
)

func ULID(ctx context.Context) string {
	return internal.ULID(ctx)
}

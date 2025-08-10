package internal

import (
	"context"

	"github.com/oklog/ulid/v2"
)

var ULID = DefaultULID

func DefaultULID(_ context.Context) string {
	return ulid.Make().String()
}

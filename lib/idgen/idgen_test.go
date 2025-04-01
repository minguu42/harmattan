package idgen_test

import (
	"testing"

	"github.com/minguu42/harmattan/lib/idgen"
)

func TestULID(t *testing.T) {
	ctx := t.Context()

	id1 := idgen.ULID(ctx)
	id2 := idgen.ULID(ctx)

	if id1 == "" || id2 == "" {
		t.Errorf("ULID() returned an empty string")
	}
	if len(id1) != 26 {
		t.Errorf("Length of ULID = %d, want 26", len(id1))
	}
	if id1 == id2 {
		t.Errorf("ULID() returned duplicate values: id1 = %s, id2 = %s", id1, id2)
	}
}

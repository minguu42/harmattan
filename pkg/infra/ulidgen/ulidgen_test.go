package ulidgen

import "testing"

func TestGenerator_Generate(t *testing.T) {
	t.Run("ULIDの長さは26文字である", func(t *testing.T) {
		g := Generator{}
		if got := g.Generate(); len(got) != 26 {
			t.Errorf("ULID is a 26-character string, but got %d-character string", len(got))
		}
	})
}

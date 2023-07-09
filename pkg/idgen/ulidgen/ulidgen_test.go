package ulidgen

import "testing"

func TestGenerator_Generate(t *testing.T) {
	t.Run("ULIDは26文字の文字列である", func(t *testing.T) {
		g := Generator{}
		if got := g.Generate(); len(got) != 26 {
			t.Errorf("ULID is a 26-character string, but got %d-character string", len(got))
		}
	})
}

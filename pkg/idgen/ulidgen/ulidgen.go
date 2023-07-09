// Package ulidgen は ULID を生成するに関するパッケージ
package ulidgen

import "github.com/oklog/ulid/v2"

type Generator struct{}

// Generate はULIDを生成する
func (g *Generator) Generate() string {
	return ulid.Make().String()
}

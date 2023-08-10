// Package ulidgen は ULID を生成するに関するパッケージ
package ulidgen

import "github.com/oklog/ulid/v2"

type Generator struct{}

// Generate は ULID を生成する
func (g *Generator) Generate() string {
	return ulid.Make().String()
}

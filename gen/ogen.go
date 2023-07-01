// Package gen は生成コマンドを含むパッケージ
package gen

//go:generate go run github.com/ogen-go/ogen/cmd/ogen@latest -clean -no-client -no-webhook-client -no-webhook-server -package ogen -target ./ogen ../doc/openapi.yaml

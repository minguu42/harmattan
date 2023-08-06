// Package gen は生成コマンドを含むパッケージ
package gen

//go:generate ogen -clean -no-client -no-webhook-client -no-webhook-server -package ogen -target ./ogen ../doc/openapi.yaml
//go:generate sqlc generate -f ../sqlc.yaml

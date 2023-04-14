VERSION  := 0.1.0
REVISION := $(shell git rev-parse --short HEAD)
LDFLAGS  := "-X main.revision=$(REVISION)"

.DEFAULT_GOAL := help
.PHONY: setup build run fmt lint test help

setup: ## 開発に必要なツールをインストールする
	go install golang.org/x/tools/cmd/goimports@latest
	go install honnef.co/go/tools/cmd/staticcheck@latest

build: ## serverプログラムをビルドする
	go build -ldflags $(LDFLAGS) -o server ./cmd/server

run: ## serverプログラムを実行する
	@docker compose up -d db-local
	@docker compose up api

fmt: ## フォーマットを実行する
	@goimports -w .

lint: ## 静的解析を実行する
	@go vet ./...
	@staticcheck ./...

test: ## テストを実行する
	@go test ./...

help: ## ヘルプを表示する
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-10s\033[0m %s\n", $$1, $$2}'

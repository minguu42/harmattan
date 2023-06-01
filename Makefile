VERSION  := 0.1.0
REVISION := $(shell git rev-parse --short HEAD)

.DEFAULT_GOAL := help
.PHONY: setup build run gen fmt lint test help

setup: ## 開発に必要なツールをインストールする
	go install golang.org/x/tools/cmd/goimports@latest
	go install honnef.co/go/tools/cmd/staticcheck@latest
	go install github.com/ogen-go/ogen/cmd/ogen@latest

build: ## serverプログラムを含むDockerイメージをビルドする
	@docker build \
            --build-arg="API_VERSION=v$(VERSION)" \
            --build-arg="API_REVISION=$(REVISION)" \
            --tag=mtasks-api --target=prod .

run: ## serverプログラムを実行する
	@docker compose up -d db-local
	@docker compose up api

gen: ## コードを生成する
	@go generate ./gen

fmt: ## フォーマットを実行する
	@goimports -w .

lint: ## 静的解析を実行する
	@go vet $$(go list ./... | grep -v /pkg/ogen)
	@staticcheck $$(go list ./... | grep -v /pkg/ogen)

test: ## テストを実行する
	@go test ./...

help: ## ヘルプを表示する
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) \
      | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-10s\033[0m %s\n", $$1, $$2}'

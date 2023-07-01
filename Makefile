VERSION  := v0.1.2
REVISION := $(shell git rev-parse --short HEAD)

.DEFAULT_GOAL := help
.PHONY: setup build run gen fmt lint test help

setup: ## 開発に必要なツールをインストールする
	go install golang.org/x/tools/cmd/goimports@latest
	go install honnef.co/go/tools/cmd/staticcheck@latest
	go install github.com/ogen-go/ogen/cmd/ogen@latest

build: ## APIサーバのコンテナイメージをビルドする
	@docker build \
            --build-arg="API_VERSION=$(VERSION)" \
            --build-arg="API_REVISION=$(REVISION)" \
            --tag=mtasks-api --target=prod .

run: ## APIサーバを実行する
	@docker compose --env-file .env.local up api

gen: ## コードを生成する
	@go generate ./gen

fmt: ## フォーマットを実行する
	@goimports -w .

lint: ## 静的解析を実行する
	@go vet $$(go list ./... | grep -v /gen)
	@staticcheck $$(go list ./... | grep -v /gen)

test: ## テストを実行する
	@go test ./...

help: ## ヘルプを表示する
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) \
      | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-10s\033[0m %s\n", $$1, $$2}'

VERSION  := $(shell yq '.info.version' ./doc/openapi.yaml)
REVISION := $(shell git rev-parse --short HEAD)

.DEFAULT_GOAL := help
.PHONY: setup gen build run dev fmt lint test help

setup: ## 開発に必要なツールをインストールする
	go install github.com/ogen-go/ogen/cmd/ogen@latest
	go install github.com/golang/mock/mockgen@latest
	go install golang.org/x/tools/cmd/goimports@latest
	go install honnef.co/go/tools/cmd/staticcheck@latest

gen: ## コードを生成する
	@go generate ./...

build: ## APIサーバのコンテナイメージをビルドする
	@docker build \
            --build-arg="API_VERSION=$(VERSION)" \
            --build-arg="API_REVISION=$(REVISION)" \
            --tag=opepe-api:latest --tag=opepe-api:$(VERSION) \
            --target=prod .

run: ## APIサーバを実行する
	@docker container run \
            --env-file .env.local \
            --name opepe-api \
            --network=opepe_default \
            -p 8080:8080 \
            --rm \
            opepe-api

dev: ## 開発用のAPIサーバを実行する
	@docker compose up api

fmt: ## コードを整形する
	@goimports -w .

lint: ## 静的解析を実行する
	@go vet $$(go list ./... | grep -v /gen)
	@staticcheck $$(go list ./... | grep -v /gen)

test: ## テストを実行する
	@go test $$(go list ./... | grep -v /gen)

help: ## ヘルプを表示する
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) \
      | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-10s\033[0m %s\n", $$1, $$2}'

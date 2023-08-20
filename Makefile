VERSION  := $(shell yq '.info.version' ./doc/openapi.yaml)

.DEFAULT_GOAL := help
.PHONY: setup gen build run dev fmt lint test help

setup: ## 開発に必要なツールをインストールする
	go install github.com/ogen-go/ogen/cmd/ogen@latest
	go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest
	go install go.uber.org/mock/mockgen@latest
	go install golang.org/x/tools/cmd/goimports@latest
	go install honnef.co/go/tools/cmd/staticcheck@latest

gen: ## コードを生成する
	@go generate ./...

build: ## 本番用APIサーバのコンテナイメージをビルドする
	@docker build \
            --build-arg="API_VERSION=v$(VERSION)" \
            --tag=opepe-api:latest \
            --target=prod .

run: ## 本番用APIサーバを実行する
	@docker compose up -d db
	@docker container run \
            --env-file .env.local \
            --name opepe-api \
            --network=opepe_default \
            -p 8080:8080 \
            --rm \
            opepe-api

dev: ## 開発用APIサーバを実行する
	@docker compose run \
            --name opepe-api \
            -p 8080:8080 \
            --rm \
            api

fmt: ## コードを整形する
	@goimports -l -w .

lint: ## 静的解析を実行する
	@go vet $$(go list ./... | grep -v /gen)
	@staticcheck $$(go list ./... | grep -v /gen)

test: ## テストを実行する
	@go test $$(go list ./... | grep -v /gen)

help: ## ヘルプを表示する
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) \
      | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-10s\033[0m %s\n", $$1, $$2}'

.DEFAULT_GOAL := help

.PHONY: up
up: ## 開発サーバを立ち上げる
	@docker compose up -d

.PHONY: down
down: ## 開発サーバを停止し、削除する
	@docker compose down

.PHONY: gen
gen: ## コードを生成する
	@go generate ./...

.PHONY: fmt
fmt: ## コードを整形する
	@go tool goimports -w .

.PHONY: lint
lint: ## 静的解析を実行する
	@go vet ./...
	@go tool staticcheck ./...

.PHONY: test
test: ## テストを実行する
	@go test -shuffle=on ./...

.PHONY: build
build: ## 本番環境向けのコンテナイメージをビルドする
	@docker image build --file=./cmd/api/Dockerfile --tag=harmattan-api:$$(git rev-parse main) --target=prod .

.PHONY: help
help: ## ヘルプを表示する
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(firstword $(MAKEFILE_LIST)) \
      | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-14s\033[0m %s\n", $$1, $$2}'

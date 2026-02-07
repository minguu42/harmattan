# CLAUDE.md

## 概要

Harmattanはタスク管理アプリである。
フロントエンドはReact、バックエンドはGoで作成されている。

## ディレクトリ構成

- `cmd`: バックエンドアプリケーションのエントリポイント
- `doc`: OpenAPIドキュメント
- `infra`: インフラストラクチャ関係のコード（詳細は`infra/CLAUDE.md`を参照）
- `internal`: バックエンドコード
- `web`: フロントエンドコード（詳細は`web/CLAUDE.md`を参照）

## バックエンド(Go)

### アーキテクチャ

- **レイヤー構造**: Handler → Usecase → Database/Domain
  - `internal/api/handler`: OpenAPI仕様に基づくHTTPハンドラー
  - `internal/api/usecase`: ビジネスロジック層
  - `internal/database`: GORMを使用したデータアクセス層
  - `internal/domain`: Harmattanのメンタルモデルを素直に表現するドメインモデル
- **OpenAPI駆動**: `doc/openapi.yaml`から`ogen`により`internal/api/openapi`配下のコードを生成
- **認証**: JWT (internal/auth)
- **ログ**: 構造化ログ (internal/atel)

### コマンド

- `go generate ./...`: OpenAPIやstringerでコードを生成する
- `gofmt -s -w . && go tool goimports -w .`: コードを整形する
- `go build ./... && go vet ./... && go tool staticcheck ./...`: 静的解析を実行する
- `go test -shuffle=on ./...`: テストを実行する
- `go test -shuffle=on ./internal/api/handler`: 特定のパッケージのテストを実行する
- `go test -shuffle=on ./internal/api/handler -run TestHandlerName`: 特定のテストを実行する

### 開発ノート

#### Code Generation
- `internal/api/openapi`のコードは自動生成のため直接編集不可
- OpenAPI仕様を変更した場合は`go generate ./internal/api/...`を実行

#### Testing
- データベース関連のテストは`testcontainers-go`を使用したコンテナベーステストを実行
- テストヘルパーは`internal/database/databasetest`にある
- テストでの列挙型の扱い:
  - 値が意味を持つ列挙型は文字列リテラルで代入（例: `Color: "blue"` not `Color: domain.ProjectColorBlue`）
  - iotaなど値が意味を持たない列挙型は定数を使用

#### Environment Variables
- APIサーバーの環境変数は`internal/api/config.go`で`env.Load[api.Config]()`により読み込み
- ローカル開発時は`cmd/api/.env`ファイルを使用

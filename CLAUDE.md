# CLAUDE.md

## 概要

Harmattanはタスク管理アプリである。
フロントエンドはReact、バックエンドAPIはGoで作成されている。
インフラストラクチャはAWSで、そのリソースはTerraformで管理されているが、開発中のため、すべてのリソースがTerraform管理対象外のものもある。

## ディレクトリ構成

- `cmd`: バックエンドコードのエントリポイント
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
- **ログ**: 構造化ログ (internal/alog)

### コマンド

```bash
# コード生成 (OpenAPI、stringer等)
go generate ./cmd/api/... ./internal/...

# フォーマット
gofmt -l -s -w ./cmd/api ./internal
go tool goimports -l -w ./cmd/api ./internal

# リント
go vet ./cmd/api/... ./internal/...
go tool staticcheck ./cmd/api/... ./internal/...

# ビルドチェック
go build -o /dev/null ./cmd/api

# テスト
go test -shuffle=on ./cmd/api/... ./internal/...

# 特定のパッケージのテスト
go test -shuffle=on ./internal/api/handler

# 特定のテスト関数を実行
go test -shuffle=on ./internal/api/handler -run TestHandlerName
```

### 開発ノート

#### Code Generation
- `internal/api/openapi`のコードは自動生成のため直接編集不可
- OpenAPI仕様を変更した場合は`go generate ./internal/api/...`を実行

#### Testing
- データベース関連のテストは`testcontainers-go`を使用したコンテナベーステストを実行
- テストヘルパーは`internal/database/databasetest`にある

#### Environment Variables
- APIサーバーの環境変数は`internal/api/config.go`で`env.Load[api.Config]()`により読み込み
- ローカル開発時は`cmd/api/.env`ファイルを使用

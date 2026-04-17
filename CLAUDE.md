# CLAUDE.md

## 概要

Harmattanはタスク管理アプリである。
フロントエンドはReact、バックエンドはGoで作成されている。

## ディレクトリ構成

- `cmd`: バックエンドアプリケーションのエントリポイント
- `doc`: OpenAPIドキュメント
- `infra`: インフラストラクチャ関係のコード
- `internal`: バックエンドコード
- `web`: フロントエンドコード
- ルート直下のフロントエンド設定（`package.json`、`vite.config.ts`、`tsconfig.*.json`、`eslint.config.js`、`pnpm-*.yaml`）はリポジトリルートで`pnpm`コマンドを実行する都合で配置している

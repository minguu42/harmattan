---
paths:
  - "web/**/*.{ts,tsx,css,html}"
  - "*.config.{ts,js}"
  - "tsconfig*.json"
  - "package.json"
  - "pnpm-*.yaml"
  - "eslint.config.js"
---

# フロントエンド規約

## 技術スタック

- React 19
- TypeScript
- Vite
- TanStack Router
- TanStack Query
- Tailwind CSS 4

## ディレクトリ構成

- `web/src/routes`: TanStack Routerによるファイルベースルーティング
- `web/src/components`: 共通コンポーネント
- `web/src/api`: バックエンドAPIクライアント

## コマンド

- `pnpm dev`: 開発サーバー起動
- `pnpm build`: ビルド
- `pnpm lint`: Lint
- `pnpm preview`: プレビュー

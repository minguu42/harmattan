---
paths:
  - "infra/**"
---

# インフラストラクチャ規約

## 構成

- **AWS Lambda**: `infra/lambdas`配下に関数を配置
- **Terraform**: `infra/terraform`でインフラ管理
  - リモート状態管理用S3バケットとGitHub Actions用IAMリソースはTerraform管理対象外
- **MySQL**: `infra/mysql`にスキーマと設定（`schema.sql`、`my.cnf`）
- **OpenTelemetry Collector**: `infra/otelcol`に設定（`config.yaml`、`config.local.yaml`）

## コマンド

- `cd infra/terraform && terraform init -backend-config=stg.tfbackend`: Terraform初期化
- `cd infra/terraform && terraform plan`: プラン確認
- `cd infra/terraform && terraform apply`: 適用

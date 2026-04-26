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

- `terraform -chdir=./infra/terraform init -backend-config=stg.tfbackend -reconfigure`: モジュールを初期化する
- `terraform fmt -recursive ./infra/terraform`: コードを整形する
- `terraform -chdir=./infra/terraform plan -var "env=stg"`: プランを確認する
- `terraform -chdir=./infra/terraform apply -var "env=stg"`: プランを適用する

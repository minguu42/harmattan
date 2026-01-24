# インフラストラクチャ

## 構成

- **AWS Lambda**: `infra/lambdas`配下に関数を配置
- **Terraform**: `infra/terraform`でインフラ管理
  - リモート状態管理用S3バケットとGitHub Actions用IAMリソースはTerraform管理対象外

## コマンド

```bash
# Terraform初期化
cd infra/terraform
terraform init -backend-config=stg.tfbackend

# プラン確認
terraform plan

# 適用
terraform apply
```

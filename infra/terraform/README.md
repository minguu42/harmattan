# Terraform管理対象外リソース一覧

- Terraformステートファイルを管理するS3バケット: ${product}-${env}-remote-tfstate
- GitHub Actions用のIAM IDプロバイダ: token.actions.githubusercontent.com
- GitHub Actions用のIAMロール: ${product}-${env}-github-actions

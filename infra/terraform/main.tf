terraform {
  required_version = "~> 1.14.0"
  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "6.13.0"
    }
  }
  backend "s3" {}
}

provider "aws" {
  region = "ap-northeast-1"
  default_tags {
    tags = {
      Product   = local.product
      Env       = var.env
      Owner     = "me"
      ManagedBy = "terraform"
    }
  }
}

variable "env" {
  type = string
  validation {
    condition     = contains(["prod", "stg"], var.env)
    error_message = "Only the following environment is allowed: prod | stg"
  }
}

locals {
  product      = "harmattan"
  isProduction = var.env == "prod"
}

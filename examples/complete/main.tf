terraform {
  required_version = ">= 1.5.0"
  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = ">= 5.40.0"
    }
    random = {
      source  = "hashicorp/random"
      version = ">= 3.6.0"
    }
    archive = {
      source  = "hashicorp/archive"
      version = ">= 2.0.0"
    }
  }
}

provider "aws" {
  region = "us-east-2"
}

locals {
  git           = "terraform-aws-github-data-lake"
  shared_secret = "testing123"
  tags = {
    git     = local.git
    cost    = "shared"
    creator = "terraform"
  }
}

module "kms" {
  source                  = "github.com/champ-oss/terraform-aws-kms.git?ref=v1.0.33-cb3be31"
  git                     = local.git
  name                    = "alias/${local.git}"
  deletion_window_in_days = 7
  account_actions         = []
}

resource "aws_kms_ciphertext" "this" {
  key_id    = module.kms.key_id
  plaintext = local.shared_secret
}

module "this" {
  source             = "../../"
  git                = local.git
  protect            = false
  buffering_interval = 60
  shared_secret      = aws_kms_ciphertext.this.ciphertext_blob
  tags               = local.tags
  prefix             = "test/"
}

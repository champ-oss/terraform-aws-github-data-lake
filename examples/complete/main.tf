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
  source                  = "github.com/champ-oss/terraform-aws-kms.git?ref=v1.0.32-a3f9aff"
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
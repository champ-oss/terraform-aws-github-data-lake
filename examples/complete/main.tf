provider "aws" {
  region = "us-east-2"
}

locals {
  git = "terraform-aws-github-data-lake"
  tags = {
    git     = local.git
    cost    = "shared"
    creator = "terraform"
  }
}

module "kms" {
  source                  = "github.com/champ-oss/terraform-aws-kms.git?ref=v1.0.28-8a5df9c"
  git                     = local.git
  name                    = "alias/${local.git}"
  deletion_window_in_days = 7
  account_actions         = []
}

resource "aws_kms_ciphertext" "this" {
  key_id    = module.kms.key_id
  plaintext = var.shared_secret
}

module "this" {
  source          = "../../"
  git             = local.git
  protect         = false
  buffer_interval = 60
  shared_secret   = aws_kms_ciphertext.this.ciphertext_blob
  tags            = local.tags
}
provider "aws" {
  region = "us-east-2"
}

module "this" {
  source          = "../../"
  git             = "terraform-aws-github-data-lake"
  protect         = false
  buffer_interval = 60
  shared_secret   = "testing123"
}
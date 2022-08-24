data "aws_region" "this" {}

locals {
  tags = {
    cost    = "shared"
    creator = "terraform"
    git     = var.git
  }
}
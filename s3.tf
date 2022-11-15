module "s3" {
  source  = "github.com/champ-oss/terraform-aws-s3.git?ref=v1.0.35-1940fd8"
  git     = var.git
  protect = var.protect
  tags    = merge(local.tags, var.tags)
}
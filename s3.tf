module "s3" {
  source  = "github.com/champ-oss/terraform-aws-s3.git?ref=v1.0.38-de0af1f"
  git     = var.git
  protect = var.protect
  tags    = merge(local.tags, var.tags)
}
resource "aws_sns_topic" "this" {
  name_prefix = "${var.git}-"
  tags        = merge(local.tags, var.tags)
}
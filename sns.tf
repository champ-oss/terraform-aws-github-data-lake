resource "aws_sns_topic" "this" {
  name_prefix = "${var.git}-"
  tags        = merge(local.tags, var.tags)
}

resource "aws_sns_topic_subscription" "this" {
  topic_arn             = aws_sns_topic.this.arn
  protocol              = "firehose"
  endpoint              = aws_kinesis_firehose_delivery_stream.this.arn
  subscription_role_arn = aws_iam_role.this.arn
}
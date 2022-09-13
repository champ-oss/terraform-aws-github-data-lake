resource "aws_kinesis_firehose_delivery_stream" "this" {
  name        = var.git
  destination = "extended_s3"
  tags        = merge(local.tags, var.tags)

  extended_s3_configuration {
    role_arn           = aws_iam_role.this.arn
    bucket_arn         = module.s3.arn
    prefix             = var.prefix
    buffer_size        = var.buffer_size
    buffer_interval    = var.buffer_interval
    compression_format = var.compression_format
  }
}
data "archive_file" "this" {
  type             = "zip"
  output_file_mode = "0666"
  source_file      = "${path.module}/github_event_receiver_lambda/github_event_receiver_lambda.py"
  output_path      = "${path.module}/github_event_receiver_lambda.zip"
}

module "lambda" {
  source                          = "github.com/champ-oss/terraform-aws-lambda.git?ref=v1.0.92-3e98cfe"
  git                             = var.git
  name                            = "github_event_receiver_lambda"
  tags                            = merge(local.tags, var.tags)
  runtime                         = var.runtime
  handler                         = "github_event_receiver_lambda.handler"
  filename                        = data.archive_file.this.output_path
  source_code_hash                = data.archive_file.this.output_base64sha256
  enable_function_url             = true
  function_url_authorization_type = "NONE"
  environment = {
    SNS_TOPIC_ARN = aws_sns_topic.this.arn
  }
}
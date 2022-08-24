data "aws_iam_policy_document" "this" {
  statement {
    actions = [
      "SNS:Publish"
    ]
    resources = [aws_sns_topic.this.arn]
  }
}

resource "aws_iam_policy" "this" {
  name_prefix = var.git
  policy      = data.aws_iam_policy_document.this.json
  tags        = merge(local.tags, var.tags)
}

resource "aws_iam_role_policy_attachment" "this" {
  policy_arn = aws_iam_policy.this.arn
  role       = module.lambda.role_name
}
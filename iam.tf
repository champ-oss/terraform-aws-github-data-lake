# Used for Firehose
resource "aws_iam_role" "this" {
  name_prefix        = "${var.git}-"
  tags               = merge(local.tags, var.tags)
  assume_role_policy = data.aws_iam_policy_document.assume.json
}

data "aws_iam_policy_document" "assume" {
  statement {
    actions = ["sts:AssumeRole"]
    principals {
      identifiers = [
        "firehose.amazonaws.com",
        "sns.amazonaws.com"
      ]
      type = "Service"
    }
  }
}

resource "aws_iam_role_policy_attachment" "firehose" {
  role       = aws_iam_role.this.name
  policy_arn = "arn:aws:iam::aws:policy/AmazonKinesisFirehoseFullAccess"
}

resource "aws_iam_role_policy_attachment" "s3" {
  role       = aws_iam_role.this.name
  policy_arn = "arn:aws:iam::aws:policy/AmazonS3FullAccess"
}

# Used for the Lambda
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
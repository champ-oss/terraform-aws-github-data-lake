variable "git" {
  description = "Identifier to be used on all resources"
  type        = string
}

variable "shared_secret" {
  description = "KMS encrypted ciphertext (https://docs.github.com/en/developers/webhooks-and-events/webhooks/securing-your-webhooks)"
  type        = string
}

variable "tags" {
  description = "https://docs.aws.amazon.com/general/latest/gr/aws_tagging.html"
  type        = map(string)
  default     = {}
}

variable "protect" {
  description = "Enables deletion protection on eligible resources"
  type        = bool
  default     = true
}

variable "runtime" {
  description = "https://docs.aws.amazon.com/lambda/latest/dg/lambda-runtimes.html"
  type        = string
  default     = "python3.8"
}

variable "buffer_size" {
  description = "https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/kinesis_firehose_delivery_stream#buffer_size"
  type        = number
  default     = 5
}

variable "buffer_interval" {
  description = "https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/kinesis_firehose_delivery_stream#buffer_interval"
  type        = number
  default     = 300
}

variable "compression_format" {
  description = "https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/kinesis_firehose_delivery_stream#compression_format"
  type        = string
  default     = "UNCOMPRESSED"
}

variable "prefix" {
  description = "https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/kinesis_firehose_delivery_stream#prefix"
  type        = string
  default     = "firehose/"
}

variable "table_string_columns" {
  description = "https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/glue_catalog_table#name"
  type        = list(string)
  default = [
    "action",
    "after",
    "before",
    "changes",
    "check_suite",
    "check_run",
    "comment",
    "issue",
    "number",
    "organization",
    "pull_request",
    "repository",
    "sender",
    "workflow",
    "workflow_job",
    "workflow_run",
  ]
}
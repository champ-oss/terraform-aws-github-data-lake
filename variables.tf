variable "git" {
  description = "Identifier to be used on all resources"
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
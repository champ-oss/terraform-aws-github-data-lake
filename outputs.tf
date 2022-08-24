output "function_url" {
  description = "https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/lambda_function_url#function_url"
  value       = module.lambda.function_url
}

output "function_arn" {
  description = "https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/lambda_function#arn"
  value       = module.lambda.arn
}

output "function_name" {
  description = "https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/lambda_function#function_name"
  value       = module.lambda.function_name
}

output "bucket" {
  description = "S3 bucket name"
  value       = module.s3.bucket
}

output "region" {
  description = "AWS Region"
  value       = data.aws_region.this.name
}
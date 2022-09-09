output "function_url" {
  description = "https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/lambda_function_url#function_url"
  value       = module.this.function_url
}

output "bucket" {
  description = "S3 bucket name"
  value       = module.this.bucket
}

output "region" {
  description = "AWS Region"
  value       = module.this.region
}

output "signature_header_key" {
  description = "Signature Header Key"
  value       = module.this.signature_header_key
}
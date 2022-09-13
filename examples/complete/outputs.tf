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

output "database" {
  description = "https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/glue_catalog_database"
  value       = module.this.database
}

output "table" {
  description = "https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/glue_catalog_table"
  value       = module.this.table
}
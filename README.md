# terraform-aws-github-data-lake

A Terraform module for ingesting GitHub event data

[![.github/workflows/lint.yml](https://github.com/champ-oss/terraform-aws-github-data-lake/actions/workflows/lint.yml/badge.svg)](https://github.com/champ-oss/terraform-aws-github-data-lake/actions/workflows/lint.yml)
[![.github/workflows/module.yml](https://github.com/champ-oss/terraform-aws-github-data-lake/actions/workflows/module.yml/badge.svg)](https://github.com/champ-oss/terraform-aws-github-data-lake/actions/workflows/module.yml)
[![.github/workflows/pytest.yml](https://github.com/champ-oss/terraform-aws-github-data-lake/actions/workflows/pytest.yml/badge.svg)](https://github.com/champ-oss/terraform-aws-github-data-lake/actions/workflows/pytest.yml)
[![.github/workflows/sonar.yml](https://github.com/champ-oss/terraform-aws-github-data-lake/actions/workflows/sonar.yml/badge.svg)](https://github.com/champ-oss/terraform-aws-github-data-lake/actions/workflows/sonar.yml)

[![SonarCloud](https://sonarcloud.io/images/project_badges/sonarcloud-black.svg)](https://sonarcloud.io/summary/new_code?id=terraform-aws-github-data-lake_champ-oss)

[![Quality Gate Status](https://sonarcloud.io/api/project_badges/measure?project=terraform-aws-github-data-lake_champ-oss&metric=alert_status)](https://sonarcloud.io/summary/new_code?id=terraform-aws-github-data-lake_champ-oss)
[![Vulnerabilities](https://sonarcloud.io/api/project_badges/measure?project=terraform-aws-github-data-lake_champ-oss&metric=vulnerabilities)](https://sonarcloud.io/summary/new_code?id=terraform-aws-github-data-lake_champ-oss)
[![Reliability Rating](https://sonarcloud.io/api/project_badges/measure?project=terraform-aws-github-data-lake_champ-oss&metric=reliability_rating)](https://sonarcloud.io/summary/new_code?id=terraform-aws-github-data-lake_champ-oss)

## Features

- AWS Lambda function to act as a receiver for HTTP webhook events from
  GitHub ([about GitHub webhooks](https://docs.github.com/en/developers/webhooks-and-events/webhooks/about-webhooks))
- Supports GitHub shared secret to secure the
  endpoint ([more information](https://docs.github.com/en/developers/webhooks-and-events/webhooks/securing-your-webhooks))
- AWS Kinesis Data Firehose receives the event data and writes to S3 in JSON format
- AWS Athena table is created and configured to query data

## Example Usage

Look at `examples/complete/main.tf` for an example of how to deploy this Terraform module

## Querying Data

AWS Athena can be used to query the data in S3. This Terraform module sets up the Athena table so that it is possible to
immediately begin running queries.

Below is an example query for extracting nested JSON fields:

```sql
select json_extract_scalar(repository, '$.name')        as name,
       json_extract_scalar(repository, '$.owner.login') as login
from "my-athena-table"
```

## Testing

Several integration tests are run in `test/src/examples_complete_test.go`, which are executed on each commit to this
repository.

- An HTTP test event is sent to the Lambda function with a secret set to validate the event is received successfully
- An HTTP test event is sent to the Lambda function without a secret set to validate the event is rejected
- The S3 bucket is inspected to ensure the test event data is written successfully
- An AWS Athena query is executed and the result is checked to ensure data is returned successfully

<!-- BEGIN_TF_DOCS -->
## Requirements

| Name | Version |
|------|---------|
| <a name="requirement_terraform"></a> [terraform](#requirement\_terraform) | >= 0.15.1 |
| <a name="requirement_aws"></a> [aws](#requirement\_aws) | >= 4.17.1 |

## Providers

| Name | Version |
|------|---------|
| <a name="provider_archive"></a> [archive](#provider\_archive) | n/a |
| <a name="provider_aws"></a> [aws](#provider\_aws) | >= 4.17.1 |

## Modules

| Name | Source | Version |
|------|--------|---------|
| <a name="module_lambda"></a> [lambda](#module\_lambda) | github.com/champ-oss/terraform-aws-lambda.git | v1.0.97-948bb8b |
| <a name="module_s3"></a> [s3](#module\_s3) | github.com/champ-oss/terraform-aws-s3.git | v1.0.29-4a98121 |

## Resources

| Name | Type |
|------|------|
| [aws_glue_catalog_database.this](https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/glue_catalog_database) | resource |
| [aws_glue_catalog_table.this](https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/glue_catalog_table) | resource |
| [aws_iam_policy.this](https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/iam_policy) | resource |
| [aws_iam_role.this](https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/iam_role) | resource |
| [aws_iam_role_policy_attachment.firehose](https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/iam_role_policy_attachment) | resource |
| [aws_iam_role_policy_attachment.s3](https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/iam_role_policy_attachment) | resource |
| [aws_iam_role_policy_attachment.this](https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/iam_role_policy_attachment) | resource |
| [aws_kinesis_firehose_delivery_stream.this](https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/kinesis_firehose_delivery_stream) | resource |
| [aws_sns_topic.this](https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/sns_topic) | resource |
| [aws_sns_topic_subscription.this](https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/sns_topic_subscription) | resource |
| [archive_file.this](https://registry.terraform.io/providers/hashicorp/archive/latest/docs/data-sources/file) | data source |
| [aws_iam_policy_document.assume](https://registry.terraform.io/providers/hashicorp/aws/latest/docs/data-sources/iam_policy_document) | data source |
| [aws_iam_policy_document.this](https://registry.terraform.io/providers/hashicorp/aws/latest/docs/data-sources/iam_policy_document) | data source |
| [aws_region.this](https://registry.terraform.io/providers/hashicorp/aws/latest/docs/data-sources/region) | data source |

## Inputs

| Name | Description | Type | Default | Required |
|------|-------------|------|---------|:--------:|
| <a name="input_buffer_interval"></a> [buffer\_interval](#input\_buffer\_interval) | https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/kinesis_firehose_delivery_stream#buffer_interval | `number` | `300` | no |
| <a name="input_buffer_size"></a> [buffer\_size](#input\_buffer\_size) | https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/kinesis_firehose_delivery_stream#buffer_size | `number` | `5` | no |
| <a name="input_compression_format"></a> [compression\_format](#input\_compression\_format) | https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/kinesis_firehose_delivery_stream#compression_format | `string` | `"UNCOMPRESSED"` | no |
| <a name="input_git"></a> [git](#input\_git) | Identifier to be used on all resources | `string` | n/a | yes |
| <a name="input_prefix"></a> [prefix](#input\_prefix) | https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/kinesis_firehose_delivery_stream#prefix | `string` | `"firehose/"` | no |
| <a name="input_protect"></a> [protect](#input\_protect) | Enables deletion protection on eligible resources | `bool` | `true` | no |
| <a name="input_runtime"></a> [runtime](#input\_runtime) | https://docs.aws.amazon.com/lambda/latest/dg/lambda-runtimes.html | `string` | `"python3.8"` | no |
| <a name="input_shared_secret"></a> [shared\_secret](#input\_shared\_secret) | https://docs.github.com/en/developers/webhooks-and-events/webhooks/securing-your-webhooks | `string` | n/a | yes |
| <a name="input_table_string_columns"></a> [table\_string\_columns](#input\_table\_string\_columns) | https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/glue_catalog_table#name | `list(string)` | <pre>[<br>  "action",<br>  "after",<br>  "before",<br>  "changes",<br>  "check_suite",<br>  "check_run",<br>  "comment",<br>  "issue",<br>  "number",<br>  "organization",<br>  "pull_request",<br>  "repository",<br>  "sender",<br>  "workflow",<br>  "workflow_job",<br>  "workflow_run"<br>]</pre> | no |
| <a name="input_tags"></a> [tags](#input\_tags) | https://docs.aws.amazon.com/general/latest/gr/aws_tagging.html | `map(string)` | `{}` | no |

## Outputs

| Name | Description |
|------|-------------|
| <a name="output_bucket"></a> [bucket](#output\_bucket) | S3 bucket name |
| <a name="output_database"></a> [database](#output\_database) | https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/glue_catalog_database |
| <a name="output_function_arn"></a> [function\_arn](#output\_function\_arn) | https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/lambda_function#arn |
| <a name="output_function_name"></a> [function\_name](#output\_function\_name) | https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/lambda_function#function_name |
| <a name="output_function_url"></a> [function\_url](#output\_function\_url) | https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/lambda_function_url#function_url |
| <a name="output_region"></a> [region](#output\_region) | AWS Region |
| <a name="output_table"></a> [table](#output\_table) | https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/glue_catalog_table |
<!-- END_TF_DOCS -->
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

resource "aws_glue_catalog_database" "this" {
  name = module.s3.bucket
}

resource "aws_glue_catalog_table" "this" {
  name          = module.s3.bucket
  database_name = aws_glue_catalog_database.this.name
  table_type    = "EXTERNAL_TABLE"

  parameters = {
    "classification"                    = "json"
    "projection.datehour.format"        = "yyyy/MM/dd/HH"
    "projection.datehour.interval"      = "1"
    "projection.datehour.interval.unit" = "HOURS"
    "projection.datehour.range"         = "2020/01/01/00,NOW"
    "projection.datehour.type"          = "date"
    "projection.enabled"                = "true"
    "storage.location.template"         = "s3://${module.s3.bucket}/${datehour}/"
  }

  partition_keys {
    name = "datehour"
    type = "string"
  }

  storage_descriptor {
    input_format  = "org.apache.hadoop.mapred.TextInputFormat"
    output_format = "org.apache.hadoop.hive.ql.io.HiveIgnoreKeyTextOutputFormat"
    location      = "s3://${module.s3.bucket}/"

    ser_de_info {
      name                  = "JsonSerDe"
      serialization_library = "org.openx.data.jsonserde.JsonSerDe"
    }

    dynamic "columns" {
      for_each = var.table_string_columns
      content {
        name = columns.value
        type = "string"
      }
    }
  }
}
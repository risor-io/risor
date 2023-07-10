
resource "aws_lambda_function" "this" {
  function_name     = var.name
  description       = var.description
  role              = var.role
  handler           = var.handler
  memory_size       = var.memory_size
  runtime           = var.runtime
  timeout           = var.timeout
  kms_key_arn       = var.kms_key
  architectures     = ["x86_64"]
  # filename          = local.filename
  source_code_hash  = filebase64sha256(var.package_path)
  s3_bucket         = var.s3_bucket
  s3_key            = var.s3_key
  tags              = var.tags
  dynamic "vpc_config" {
    for_each = var.vpc_subnet_ids != null && var.vpc_security_group_ids != null ? [true] : []
    content {
      security_group_ids = var.vpc_security_group_ids
      subnet_ids         = var.vpc_subnet_ids
    }
  }
  depends_on = [
    aws_s3_object.lambda_package,
    aws_cloudwatch_log_group.lambda,
  ]
}

resource "aws_s3_object" "lambda_package" {
  bucket                 = var.s3_bucket
  key                    = var.s3_key
  source                 = var.package_path
  server_side_encryption = "AES256"
  etag                   = filemd5(var.package_path)
}

resource "aws_cloudwatch_log_group" "lambda" {
  name              = "/aws/lambda/${var.name}"
  retention_in_days = 365
  kms_key_id        = var.kms_key
  tags              = var.tags
}

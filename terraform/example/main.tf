
provider "aws" {
  region                      = "us-east-1"
  alias                       = "us-east-1"
  skip_metadata_api_check     = true
  skip_region_validation      = true
  skip_credentials_validation = true
  skip_requesting_account_id  = true
}

module "lambda" {
  providers    = { aws = aws.us-east-1 }
  source       = "../modules/aws/lambda"
  name         = "example1"
  description  = "My awesome lambda function"
  handler      = "risor-lambda"
  runtime      = "go1.x"
  package_path = "../../dist/risor-lambda.zip"
  s3_bucket    = "test-506282801638"
  s3_key       = "dist/risor-lambda.zip"
  role         = aws_iam_role.lambda.arn
}

module "vpc" {
  providers   = { aws = aws.us-east-1 }
  source      = "../modules/aws/vpc"
  name        = "example1"
}

data "aws_iam_policy_document" "assume_role" {
  statement {
    effect  = "Allow"
    actions = ["sts:AssumeRole"]
    principals {
      type        = "Service"
      identifiers = ["lambda.amazonaws.com"]
    }
  }
}

resource "aws_iam_role" "lambda" {
  name                  = "lambda-ex"
  description           = "lambda example role"
  assume_role_policy    = data.aws_iam_policy_document.assume_role.json
}

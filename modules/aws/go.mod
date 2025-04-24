module github.com/risor-io/risor/modules/aws

go 1.23.0

replace github.com/risor-io/risor => ../..

require (
	github.com/aws/aws-sdk-go-v2 v1.36.3
	github.com/aws/aws-sdk-go-v2/config v1.29.14
	github.com/aws/aws-sdk-go-v2/credentials v1.17.67
	github.com/aws/aws-sdk-go-v2/service/apigatewayv2 v1.27.1
	github.com/aws/aws-sdk-go-v2/service/athena v1.50.4
	github.com/aws/aws-sdk-go-v2/service/backup v1.41.2
	github.com/aws/aws-sdk-go-v2/service/cloudformation v1.59.2
	github.com/aws/aws-sdk-go-v2/service/cloudfront v1.45.3
	github.com/aws/aws-sdk-go-v2/service/cloudtrail v1.48.4
	github.com/aws/aws-sdk-go-v2/service/cloudwatch v1.44.3
	github.com/aws/aws-sdk-go-v2/service/cloudwatchlogs v1.47.3
	github.com/aws/aws-sdk-go-v2/service/dynamodb v1.43.0
	github.com/aws/aws-sdk-go-v2/service/ebs v1.28.3
	github.com/aws/aws-sdk-go-v2/service/ec2 v1.212.0
	github.com/aws/aws-sdk-go-v2/service/ecr v1.43.3
	github.com/aws/aws-sdk-go-v2/service/ecs v1.56.1
	github.com/aws/aws-sdk-go-v2/service/eks v1.64.0
	github.com/aws/aws-sdk-go-v2/service/elasticache v1.46.0
	github.com/aws/aws-sdk-go-v2/service/elasticsearchservice v1.33.3
	github.com/aws/aws-sdk-go-v2/service/eventbridge v1.39.0
	github.com/aws/aws-sdk-go-v2/service/firehose v1.37.4
	github.com/aws/aws-sdk-go-v2/service/glue v1.109.2
	github.com/aws/aws-sdk-go-v2/service/iam v1.41.1
	github.com/aws/aws-sdk-go-v2/service/kinesis v1.33.3
	github.com/aws/aws-sdk-go-v2/service/kms v1.38.3
	github.com/aws/aws-sdk-go-v2/service/lambda v1.71.2
	github.com/aws/aws-sdk-go-v2/service/ram v1.30.3
	github.com/aws/aws-sdk-go-v2/service/rds v1.95.0
	github.com/aws/aws-sdk-go-v2/service/redshift v1.54.3
	github.com/aws/aws-sdk-go-v2/service/route53 v1.51.1
	github.com/aws/aws-sdk-go-v2/service/s3 v1.79.2
	github.com/aws/aws-sdk-go-v2/service/secretsmanager v1.35.4
	github.com/aws/aws-sdk-go-v2/service/sesv2 v1.45.0
	github.com/aws/aws-sdk-go-v2/service/sfn v1.35.4
	github.com/aws/aws-sdk-go-v2/service/sns v1.34.4
	github.com/aws/aws-sdk-go-v2/service/sqs v1.38.5
	github.com/aws/aws-sdk-go-v2/service/sts v1.33.19
	github.com/aws/aws-sdk-go-v2/service/wafv2 v1.60.1
	github.com/aws/aws-sdk-go-v2/service/xray v1.31.3
	github.com/risor-io/risor v1.7.0
)

require (
	github.com/aws/aws-sdk-go-v2/aws/protocol/eventstream v1.6.10 // indirect
	github.com/aws/aws-sdk-go-v2/feature/ec2/imds v1.16.30 // indirect
	github.com/aws/aws-sdk-go-v2/internal/configsources v1.3.34 // indirect
	github.com/aws/aws-sdk-go-v2/internal/endpoints/v2 v2.6.34 // indirect
	github.com/aws/aws-sdk-go-v2/internal/ini v1.8.3 // indirect
	github.com/aws/aws-sdk-go-v2/internal/v4a v1.3.34 // indirect
	github.com/aws/aws-sdk-go-v2/service/internal/accept-encoding v1.12.3 // indirect
	github.com/aws/aws-sdk-go-v2/service/internal/checksum v1.7.0 // indirect
	github.com/aws/aws-sdk-go-v2/service/internal/endpoint-discovery v1.10.15 // indirect
	github.com/aws/aws-sdk-go-v2/service/internal/presigned-url v1.12.15 // indirect
	github.com/aws/aws-sdk-go-v2/service/internal/s3shared v1.18.15 // indirect
	github.com/aws/aws-sdk-go-v2/service/sso v1.25.3 // indirect
	github.com/aws/aws-sdk-go-v2/service/ssooidc v1.30.1 // indirect
	github.com/aws/smithy-go v1.22.3 // indirect
)

module github.com/risor-io/risor/modules/aws

go 1.20

replace github.com/risor-io/risor => ../..

require (
	github.com/aws/aws-sdk-go-v2 v1.21.0
	github.com/aws/aws-sdk-go-v2/config v1.18.37
	github.com/aws/aws-sdk-go-v2/credentials v1.13.35
	github.com/aws/aws-sdk-go-v2/service/apigatewayv2 v1.14.5
	github.com/aws/aws-sdk-go-v2/service/athena v1.31.6
	github.com/aws/aws-sdk-go-v2/service/backup v1.24.4
	github.com/aws/aws-sdk-go-v2/service/cloudformation v1.34.5
	github.com/aws/aws-sdk-go-v2/service/cloudfront v1.28.5
	github.com/aws/aws-sdk-go-v2/service/cloudtrail v1.28.6
	github.com/aws/aws-sdk-go-v2/service/cloudwatch v1.27.6
	github.com/aws/aws-sdk-go-v2/service/cloudwatchlogs v1.23.5
	github.com/aws/aws-sdk-go-v2/service/dynamodb v1.21.5
	github.com/aws/aws-sdk-go-v2/service/ebs v1.18.5
	github.com/aws/aws-sdk-go-v2/service/ec2 v1.115.0
	github.com/aws/aws-sdk-go-v2/service/ecr v1.19.5
	github.com/aws/aws-sdk-go-v2/service/ecs v1.29.6
	github.com/aws/aws-sdk-go-v2/service/eks v1.29.5
	github.com/aws/aws-sdk-go-v2/service/elasticache v1.29.3
	github.com/aws/aws-sdk-go-v2/service/elasticsearchservice v1.20.5
	github.com/aws/aws-sdk-go-v2/service/eventbridge v1.20.5
	github.com/aws/aws-sdk-go-v2/service/firehose v1.17.5
	github.com/aws/aws-sdk-go-v2/service/glue v1.62.0
	github.com/aws/aws-sdk-go-v2/service/iam v1.22.5
	github.com/aws/aws-sdk-go-v2/service/kinesis v1.18.5
	github.com/aws/aws-sdk-go-v2/service/kms v1.24.5
	github.com/aws/aws-sdk-go-v2/service/lambda v1.39.5
	github.com/aws/aws-sdk-go-v2/service/ram v1.20.5
	github.com/aws/aws-sdk-go-v2/service/rds v1.53.0
	github.com/aws/aws-sdk-go-v2/service/redshift v1.29.5
	github.com/aws/aws-sdk-go-v2/service/route53 v1.29.5
	github.com/aws/aws-sdk-go-v2/service/s3 v1.38.5
	github.com/aws/aws-sdk-go-v2/service/secretsmanager v1.21.3
	github.com/aws/aws-sdk-go-v2/service/sesv2 v1.19.5
	github.com/aws/aws-sdk-go-v2/service/sfn v1.19.5
	github.com/aws/aws-sdk-go-v2/service/sns v1.21.5
	github.com/aws/aws-sdk-go-v2/service/sqs v1.24.5
	github.com/aws/aws-sdk-go-v2/service/sts v1.21.5
	github.com/aws/aws-sdk-go-v2/service/wafv2 v1.37.5
	github.com/aws/aws-sdk-go-v2/service/xray v1.17.5
	github.com/risor-io/risor v0.14.1-0.20230825185206-8956c356a975
)

require (
	github.com/aws/aws-sdk-go-v2/aws/protocol/eventstream v1.4.13 // indirect
	github.com/aws/aws-sdk-go-v2/feature/ec2/imds v1.13.11 // indirect
	github.com/aws/aws-sdk-go-v2/internal/configsources v1.1.41 // indirect
	github.com/aws/aws-sdk-go-v2/internal/endpoints/v2 v2.4.35 // indirect
	github.com/aws/aws-sdk-go-v2/internal/ini v1.3.42 // indirect
	github.com/aws/aws-sdk-go-v2/internal/v4a v1.1.4 // indirect
	github.com/aws/aws-sdk-go-v2/service/internal/accept-encoding v1.9.14 // indirect
	github.com/aws/aws-sdk-go-v2/service/internal/checksum v1.1.36 // indirect
	github.com/aws/aws-sdk-go-v2/service/internal/endpoint-discovery v1.7.35 // indirect
	github.com/aws/aws-sdk-go-v2/service/internal/presigned-url v1.9.35 // indirect
	github.com/aws/aws-sdk-go-v2/service/internal/s3shared v1.15.4 // indirect
	github.com/aws/aws-sdk-go-v2/service/sso v1.13.5 // indirect
	github.com/aws/aws-sdk-go-v2/service/ssooidc v1.15.5 // indirect
	github.com/aws/smithy-go v1.14.2 // indirect
	github.com/google/go-cmp v0.5.9 // indirect
	github.com/jmespath/go-jmespath v0.4.0 // indirect
)

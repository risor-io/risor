module github.com/risor-io/risor/modules/aws

go 1.22

toolchain go1.22.0

replace github.com/risor-io/risor => ../..

require (
	github.com/aws/aws-sdk-go-v2 v1.25.3
	github.com/aws/aws-sdk-go-v2/config v1.27.7
	github.com/aws/aws-sdk-go-v2/credentials v1.17.7
	github.com/aws/aws-sdk-go-v2/service/apigatewayv2 v1.20.2
	github.com/aws/aws-sdk-go-v2/service/athena v1.40.2
	github.com/aws/aws-sdk-go-v2/service/backup v1.33.2
	github.com/aws/aws-sdk-go-v2/service/cloudformation v1.47.2
	github.com/aws/aws-sdk-go-v2/service/cloudfront v1.35.2
	github.com/aws/aws-sdk-go-v2/service/cloudtrail v1.39.0
	github.com/aws/aws-sdk-go-v2/service/cloudwatch v1.36.2
	github.com/aws/aws-sdk-go-v2/service/cloudwatchlogs v1.34.3
	github.com/aws/aws-sdk-go-v2/service/dynamodb v1.30.4
	github.com/aws/aws-sdk-go-v2/service/ebs v1.23.2
	github.com/aws/aws-sdk-go-v2/service/ec2 v1.150.1
	github.com/aws/aws-sdk-go-v2/service/ecr v1.27.2
	github.com/aws/aws-sdk-go-v2/service/ecs v1.41.2
	github.com/aws/aws-sdk-go-v2/service/eks v1.41.1
	github.com/aws/aws-sdk-go-v2/service/elasticache v1.37.3
	github.com/aws/aws-sdk-go-v2/service/elasticsearchservice v1.28.2
	github.com/aws/aws-sdk-go-v2/service/eventbridge v1.30.2
	github.com/aws/aws-sdk-go-v2/service/firehose v1.28.2
	github.com/aws/aws-sdk-go-v2/service/glue v1.77.3
	github.com/aws/aws-sdk-go-v2/service/iam v1.31.2
	github.com/aws/aws-sdk-go-v2/service/kinesis v1.27.2
	github.com/aws/aws-sdk-go-v2/service/kms v1.29.2
	github.com/aws/aws-sdk-go-v2/service/lambda v1.53.2
	github.com/aws/aws-sdk-go-v2/service/ram v1.25.2
	github.com/aws/aws-sdk-go-v2/service/rds v1.75.2
	github.com/aws/aws-sdk-go-v2/service/redshift v1.43.3
	github.com/aws/aws-sdk-go-v2/service/route53 v1.40.2
	github.com/aws/aws-sdk-go-v2/service/s3 v1.52.0
	github.com/aws/aws-sdk-go-v2/service/secretsmanager v1.28.3
	github.com/aws/aws-sdk-go-v2/service/sesv2 v1.27.1
	github.com/aws/aws-sdk-go-v2/service/sfn v1.26.2
	github.com/aws/aws-sdk-go-v2/service/sns v1.29.2
	github.com/aws/aws-sdk-go-v2/service/sqs v1.31.2
	github.com/aws/aws-sdk-go-v2/service/sts v1.28.4
	github.com/aws/aws-sdk-go-v2/service/wafv2 v1.48.0
	github.com/aws/aws-sdk-go-v2/service/xray v1.25.2
	github.com/risor-io/risor v1.5.0
)

require (
	github.com/aws/aws-sdk-go-v2/aws/protocol/eventstream v1.6.1 // indirect
	github.com/aws/aws-sdk-go-v2/feature/ec2/imds v1.15.3 // indirect
	github.com/aws/aws-sdk-go-v2/internal/configsources v1.3.3 // indirect
	github.com/aws/aws-sdk-go-v2/internal/endpoints/v2 v2.6.3 // indirect
	github.com/aws/aws-sdk-go-v2/internal/ini v1.8.0 // indirect
	github.com/aws/aws-sdk-go-v2/internal/v4a v1.3.3 // indirect
	github.com/aws/aws-sdk-go-v2/service/internal/accept-encoding v1.11.1 // indirect
	github.com/aws/aws-sdk-go-v2/service/internal/checksum v1.3.5 // indirect
	github.com/aws/aws-sdk-go-v2/service/internal/endpoint-discovery v1.9.4 // indirect
	github.com/aws/aws-sdk-go-v2/service/internal/presigned-url v1.11.5 // indirect
	github.com/aws/aws-sdk-go-v2/service/internal/s3shared v1.17.3 // indirect
	github.com/aws/aws-sdk-go-v2/service/sso v1.20.2 // indirect
	github.com/aws/aws-sdk-go-v2/service/ssooidc v1.23.2 // indirect
	github.com/aws/smithy-go v1.20.1 // indirect
	github.com/jmespath/go-jmespath v0.4.0 // indirect
	gopkg.in/yaml.v2 v2.4.0 // indirect
)

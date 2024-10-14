module github.com/risor-io/risor/modules/aws

go 1.22

toolchain go1.22.0

replace github.com/risor-io/risor => ../..

require (
	github.com/aws/aws-sdk-go-v2 v1.30.4
	github.com/aws/aws-sdk-go-v2/config v1.27.31
	github.com/aws/aws-sdk-go-v2/credentials v1.17.30
	github.com/aws/aws-sdk-go-v2/service/apigatewayv2 v1.22.6
	github.com/aws/aws-sdk-go-v2/service/athena v1.44.5
	github.com/aws/aws-sdk-go-v2/service/backup v1.36.4
	github.com/aws/aws-sdk-go-v2/service/cloudformation v1.53.5
	github.com/aws/aws-sdk-go-v2/service/cloudfront v1.38.5
	github.com/aws/aws-sdk-go-v2/service/cloudtrail v1.42.5
	github.com/aws/aws-sdk-go-v2/service/cloudwatch v1.40.5
	github.com/aws/aws-sdk-go-v2/service/cloudwatchlogs v1.37.5
	github.com/aws/aws-sdk-go-v2/service/dynamodb v1.34.6
	github.com/aws/aws-sdk-go-v2/service/ebs v1.25.5
	github.com/aws/aws-sdk-go-v2/service/ec2 v1.177.0
	github.com/aws/aws-sdk-go-v2/service/ecr v1.32.2
	github.com/aws/aws-sdk-go-v2/service/ecs v1.45.2
	github.com/aws/aws-sdk-go-v2/service/eks v1.48.2
	github.com/aws/aws-sdk-go-v2/service/elasticache v1.40.7
	github.com/aws/aws-sdk-go-v2/service/elasticsearchservice v1.30.5
	github.com/aws/aws-sdk-go-v2/service/eventbridge v1.33.5
	github.com/aws/aws-sdk-go-v2/service/firehose v1.32.2
	github.com/aws/aws-sdk-go-v2/service/glue v1.95.0
	github.com/aws/aws-sdk-go-v2/service/iam v1.35.0
	github.com/aws/aws-sdk-go-v2/service/kinesis v1.29.5
	github.com/aws/aws-sdk-go-v2/service/kms v1.35.5
	github.com/aws/aws-sdk-go-v2/service/lambda v1.58.1
	github.com/aws/aws-sdk-go-v2/service/ram v1.27.5
	github.com/aws/aws-sdk-go-v2/service/rds v1.82.2
	github.com/aws/aws-sdk-go-v2/service/redshift v1.46.6
	github.com/aws/aws-sdk-go-v2/service/route53 v1.43.0
	github.com/aws/aws-sdk-go-v2/service/s3 v1.61.0
	github.com/aws/aws-sdk-go-v2/service/secretsmanager v1.32.6
	github.com/aws/aws-sdk-go-v2/service/sesv2 v1.33.0
	github.com/aws/aws-sdk-go-v2/service/sfn v1.31.0
	github.com/aws/aws-sdk-go-v2/service/sns v1.31.5
	github.com/aws/aws-sdk-go-v2/service/sqs v1.34.5
	github.com/aws/aws-sdk-go-v2/service/sts v1.30.5
	github.com/aws/aws-sdk-go-v2/service/wafv2 v1.52.0
	github.com/aws/aws-sdk-go-v2/service/xray v1.27.5
	github.com/risor-io/risor v1.7.0
)

require (
	github.com/aws/aws-sdk-go-v2/aws/protocol/eventstream v1.6.4 // indirect
	github.com/aws/aws-sdk-go-v2/feature/ec2/imds v1.16.12 // indirect
	github.com/aws/aws-sdk-go-v2/internal/configsources v1.3.16 // indirect
	github.com/aws/aws-sdk-go-v2/internal/endpoints/v2 v2.6.16 // indirect
	github.com/aws/aws-sdk-go-v2/internal/ini v1.8.1 // indirect
	github.com/aws/aws-sdk-go-v2/internal/v4a v1.3.16 // indirect
	github.com/aws/aws-sdk-go-v2/service/internal/accept-encoding v1.11.4 // indirect
	github.com/aws/aws-sdk-go-v2/service/internal/checksum v1.3.18 // indirect
	github.com/aws/aws-sdk-go-v2/service/internal/endpoint-discovery v1.9.17 // indirect
	github.com/aws/aws-sdk-go-v2/service/internal/presigned-url v1.11.18 // indirect
	github.com/aws/aws-sdk-go-v2/service/internal/s3shared v1.17.16 // indirect
	github.com/aws/aws-sdk-go-v2/service/sso v1.22.5 // indirect
	github.com/aws/aws-sdk-go-v2/service/ssooidc v1.26.5 // indirect
	github.com/aws/smithy-go v1.20.4 // indirect
	github.com/jmespath/go-jmespath v0.4.0 // indirect
	gopkg.in/yaml.v2 v2.4.0 // indirect
)

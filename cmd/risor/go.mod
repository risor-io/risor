module github.com/risor-io/risor/cmd/risor

go 1.21

toolchain go1.21.0

replace (
	github.com/risor-io/risor => ../..
	github.com/risor-io/risor/modules/aws => ../../modules/aws
	github.com/risor-io/risor/modules/image => ../../modules/image
	github.com/risor-io/risor/modules/pgx => ../../modules/pgx
	github.com/risor-io/risor/modules/uuid => ../../modules/uuid
	github.com/risor-io/risor/os/s3fs => ../../os/s3fs
)

require (
	atomicgo.dev/keyboard v0.2.9
	github.com/aws/aws-sdk-go-v2/config v1.18.37
	github.com/aws/aws-sdk-go-v2/service/s3 v1.38.5
	github.com/fatih/color v1.15.0
	github.com/hokaccha/go-prettyjson v0.0.0-20211117102719-0474bc63780f
	github.com/mitchellh/go-homedir v1.1.0
	github.com/risor-io/risor v0.14.1-0.20230825185206-8956c356a975
	github.com/risor-io/risor/modules/aws v0.0.0-00010101000000-000000000000
	github.com/risor-io/risor/modules/image v0.0.0-00010101000000-000000000000
	github.com/risor-io/risor/modules/pgx v0.0.0-00010101000000-000000000000
	github.com/risor-io/risor/modules/uuid v0.0.0-00010101000000-000000000000
	github.com/risor-io/risor/os/s3fs v0.0.0-00010101000000-000000000000
	github.com/spf13/cobra v1.7.0
	github.com/spf13/viper v1.16.0
)

require (
	github.com/anthonynsimon/bild v0.13.0 // indirect
	github.com/aws/aws-sdk-go-v2 v1.21.0 // indirect
	github.com/aws/aws-sdk-go-v2/aws/protocol/eventstream v1.4.13 // indirect
	github.com/aws/aws-sdk-go-v2/credentials v1.13.35 // indirect
	github.com/aws/aws-sdk-go-v2/feature/ec2/imds v1.13.11 // indirect
	github.com/aws/aws-sdk-go-v2/internal/configsources v1.1.41 // indirect
	github.com/aws/aws-sdk-go-v2/internal/endpoints/v2 v2.4.35 // indirect
	github.com/aws/aws-sdk-go-v2/internal/ini v1.3.42 // indirect
	github.com/aws/aws-sdk-go-v2/internal/v4a v1.1.4 // indirect
	github.com/aws/aws-sdk-go-v2/service/apigatewayv2 v1.14.5 // indirect
	github.com/aws/aws-sdk-go-v2/service/athena v1.31.6 // indirect
	github.com/aws/aws-sdk-go-v2/service/backup v1.24.4 // indirect
	github.com/aws/aws-sdk-go-v2/service/cloudformation v1.34.5 // indirect
	github.com/aws/aws-sdk-go-v2/service/cloudfront v1.28.5 // indirect
	github.com/aws/aws-sdk-go-v2/service/cloudtrail v1.28.6 // indirect
	github.com/aws/aws-sdk-go-v2/service/cloudwatch v1.27.6 // indirect
	github.com/aws/aws-sdk-go-v2/service/cloudwatchlogs v1.23.5 // indirect
	github.com/aws/aws-sdk-go-v2/service/dynamodb v1.21.5 // indirect
	github.com/aws/aws-sdk-go-v2/service/ebs v1.18.5 // indirect
	github.com/aws/aws-sdk-go-v2/service/ec2 v1.115.0 // indirect
	github.com/aws/aws-sdk-go-v2/service/ecr v1.19.5 // indirect
	github.com/aws/aws-sdk-go-v2/service/ecs v1.29.6 // indirect
	github.com/aws/aws-sdk-go-v2/service/eks v1.29.5 // indirect
	github.com/aws/aws-sdk-go-v2/service/elasticache v1.29.3 // indirect
	github.com/aws/aws-sdk-go-v2/service/elasticsearchservice v1.20.5 // indirect
	github.com/aws/aws-sdk-go-v2/service/eventbridge v1.20.5 // indirect
	github.com/aws/aws-sdk-go-v2/service/firehose v1.17.5 // indirect
	github.com/aws/aws-sdk-go-v2/service/glue v1.62.0 // indirect
	github.com/aws/aws-sdk-go-v2/service/iam v1.22.5 // indirect
	github.com/aws/aws-sdk-go-v2/service/internal/accept-encoding v1.9.14 // indirect
	github.com/aws/aws-sdk-go-v2/service/internal/checksum v1.1.36 // indirect
	github.com/aws/aws-sdk-go-v2/service/internal/endpoint-discovery v1.7.35 // indirect
	github.com/aws/aws-sdk-go-v2/service/internal/presigned-url v1.9.35 // indirect
	github.com/aws/aws-sdk-go-v2/service/internal/s3shared v1.15.4 // indirect
	github.com/aws/aws-sdk-go-v2/service/kinesis v1.18.5 // indirect
	github.com/aws/aws-sdk-go-v2/service/kms v1.24.5 // indirect
	github.com/aws/aws-sdk-go-v2/service/lambda v1.39.5 // indirect
	github.com/aws/aws-sdk-go-v2/service/ram v1.20.5 // indirect
	github.com/aws/aws-sdk-go-v2/service/rds v1.53.0 // indirect
	github.com/aws/aws-sdk-go-v2/service/redshift v1.29.5 // indirect
	github.com/aws/aws-sdk-go-v2/service/route53 v1.29.5 // indirect
	github.com/aws/aws-sdk-go-v2/service/secretsmanager v1.21.3 // indirect
	github.com/aws/aws-sdk-go-v2/service/sesv2 v1.19.5 // indirect
	github.com/aws/aws-sdk-go-v2/service/sfn v1.19.5 // indirect
	github.com/aws/aws-sdk-go-v2/service/sns v1.21.5 // indirect
	github.com/aws/aws-sdk-go-v2/service/sqs v1.24.5 // indirect
	github.com/aws/aws-sdk-go-v2/service/sso v1.13.5 // indirect
	github.com/aws/aws-sdk-go-v2/service/ssooidc v1.15.5 // indirect
	github.com/aws/aws-sdk-go-v2/service/sts v1.21.5 // indirect
	github.com/aws/aws-sdk-go-v2/service/wafv2 v1.37.5 // indirect
	github.com/aws/aws-sdk-go-v2/service/xray v1.17.5 // indirect
	github.com/aws/smithy-go v1.14.2 // indirect
	github.com/containerd/console v1.0.3 // indirect
	github.com/fsnotify/fsnotify v1.6.0 // indirect
	github.com/gofrs/uuid v4.4.0+incompatible // indirect
	github.com/hashicorp/hcl v1.0.0 // indirect
	github.com/inconshreveable/mousetrap v1.1.0 // indirect
	github.com/jackc/pgpassfile v1.0.0 // indirect
	github.com/jackc/pgservicefile v0.0.0-20221227161230-091c0ba34f0a // indirect
	github.com/jackc/pgx/v5 v5.4.1 // indirect
	github.com/jmespath/go-jmespath v0.4.0 // indirect
	github.com/magiconair/properties v1.8.7 // indirect
	github.com/mattn/go-colorable v0.1.13 // indirect
	github.com/mattn/go-isatty v0.0.17 // indirect
	github.com/mitchellh/mapstructure v1.5.0 // indirect
	github.com/pelletier/go-toml/v2 v2.0.8 // indirect
	github.com/spf13/afero v1.9.5 // indirect
	github.com/spf13/cast v1.5.1 // indirect
	github.com/spf13/jwalterweatherman v1.1.0 // indirect
	github.com/spf13/pflag v1.0.5 // indirect
	github.com/subosito/gotenv v1.4.2 // indirect
	golang.org/x/crypto v0.9.0 // indirect
	golang.org/x/image v0.5.0 // indirect
	golang.org/x/sys v0.8.0 // indirect
	golang.org/x/text v0.9.0 // indirect
	gopkg.in/ini.v1 v1.67.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)

module github.com/risor-io/risor

go 1.20

require (
	atomicgo.dev/keyboard v0.2.9
	github.com/anthonynsimon/bild v0.13.0
	github.com/aws/aws-lambda-go v1.41.0
	github.com/aws/aws-sdk-go-v2 v1.20.0
	github.com/aws/aws-sdk-go-v2/config v1.18.32
	github.com/aws/aws-sdk-go-v2/credentials v1.13.31
	github.com/aws/aws-sdk-go-v2/service/apigatewayv2 v1.14.1
	github.com/aws/aws-sdk-go-v2/service/athena v1.31.1
	github.com/aws/aws-sdk-go-v2/service/backup v1.23.1
	github.com/aws/aws-sdk-go-v2/service/cloudformation v1.30.0
	github.com/aws/aws-sdk-go-v2/service/cloudfront v1.26.8
	github.com/aws/aws-sdk-go-v2/service/cloudtrail v1.27.1
	github.com/aws/aws-sdk-go-v2/service/cloudwatch v1.26.2
	github.com/aws/aws-sdk-go-v2/service/cloudwatchlogs v1.23.1
	github.com/aws/aws-sdk-go-v2/service/dynamodb v1.19.11
	github.com/aws/aws-sdk-go-v2/service/ebs v1.16.14
	github.com/aws/aws-sdk-go-v2/service/ec2 v1.102.0
	github.com/aws/aws-sdk-go-v2/service/ecr v1.18.13
	github.com/aws/aws-sdk-go-v2/service/ecs v1.27.4
	github.com/aws/aws-sdk-go-v2/service/eks v1.27.14
	github.com/aws/aws-sdk-go-v2/service/elasticache v1.27.2
	github.com/aws/aws-sdk-go-v2/service/elasticsearchservice v1.19.2
	github.com/aws/aws-sdk-go-v2/service/eventbridge v1.20.1
	github.com/aws/aws-sdk-go-v2/service/firehose v1.17.1
	github.com/aws/aws-sdk-go-v2/service/glue v1.52.0
	github.com/aws/aws-sdk-go-v2/service/iam v1.21.0
	github.com/aws/aws-sdk-go-v2/service/kinesis v1.17.14
	github.com/aws/aws-sdk-go-v2/service/kms v1.22.2
	github.com/aws/aws-sdk-go-v2/service/lambda v1.37.0
	github.com/aws/aws-sdk-go-v2/service/ram v1.20.1
	github.com/aws/aws-sdk-go-v2/service/rds v1.46.0
	github.com/aws/aws-sdk-go-v2/service/redshift v1.29.1
	github.com/aws/aws-sdk-go-v2/service/route53 v1.28.3
	github.com/aws/aws-sdk-go-v2/service/s3 v1.36.0
	github.com/aws/aws-sdk-go-v2/service/secretsmanager v1.20.1
	github.com/aws/aws-sdk-go-v2/service/sesv2 v1.19.1
	github.com/aws/aws-sdk-go-v2/service/sfn v1.18.0
	github.com/aws/aws-sdk-go-v2/service/sns v1.20.13
	github.com/aws/aws-sdk-go-v2/service/sqs v1.23.2
	github.com/aws/aws-sdk-go-v2/service/wafv2 v1.35.1
	github.com/aws/aws-sdk-go-v2/service/xray v1.17.1
	github.com/fatih/color v1.15.0
	github.com/gofrs/uuid v4.4.0+incompatible
	github.com/hashicorp/go-multierror v1.1.1
	github.com/hokaccha/go-prettyjson v0.0.0-20211117102719-0474bc63780f
	github.com/jackc/pgx/v5 v5.4.1
	github.com/jdbaldry/go-language-server-protocol v0.0.0-20211013214444-3022da0884b2
	github.com/mitchellh/go-homedir v1.1.0
	github.com/rs/zerolog v1.29.1
	github.com/spf13/cobra v1.7.0
	github.com/spf13/viper v1.16.0
	github.com/stretchr/testify v1.8.3
)

require (
	cloud.google.com/go/compute v1.23.0 // indirect
	cloud.google.com/go/compute/metadata v0.2.3 // indirect
	github.com/fsnotify/fsnotify v1.6.0 // indirect
	github.com/golang/groupcache v0.0.0-20210331224755-41bb18bfe9da // indirect
	github.com/golang/protobuf v1.5.3 // indirect
	github.com/google/s2a-go v0.1.4 // indirect
	github.com/googleapis/enterprise-certificate-proxy v0.2.3 // indirect
	github.com/googleapis/gax-go/v2 v2.11.0 // indirect
	github.com/hashicorp/hcl v1.0.0 // indirect
	github.com/inconshreveable/mousetrap v1.1.0 // indirect
	github.com/magiconair/properties v1.8.7 // indirect
	github.com/mitchellh/mapstructure v1.5.0 // indirect
	github.com/pelletier/go-toml/v2 v2.0.8 // indirect
	github.com/spf13/afero v1.9.5 // indirect
	github.com/spf13/cast v1.5.1 // indirect
	github.com/spf13/jwalterweatherman v1.1.0 // indirect
	github.com/spf13/pflag v1.0.5 // indirect
	github.com/subosito/gotenv v1.4.2 // indirect
	go.opencensus.io v0.24.0 // indirect
	golang.org/x/net v0.10.0 // indirect
	golang.org/x/oauth2 v0.8.0 // indirect
	google.golang.org/api v0.126.0 // indirect
	google.golang.org/appengine v1.6.7 // indirect
	google.golang.org/genproto v0.0.0-20230530153820-e85fd2cbaebc // indirect
	google.golang.org/genproto/googleapis/api v0.0.0-20230530153820-e85fd2cbaebc // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20230530153820-e85fd2cbaebc // indirect
	google.golang.org/grpc v1.55.0 // indirect
	google.golang.org/protobuf v1.30.0 // indirect
	gopkg.in/ini.v1 v1.67.0 // indirect
)

require (
	github.com/aws/aws-sdk-go-v2/aws/protocol/eventstream v1.4.10 // indirect
	github.com/aws/aws-sdk-go-v2/feature/ec2/imds v1.13.7 // indirect
	github.com/aws/aws-sdk-go-v2/internal/configsources v1.1.37 // indirect
	github.com/aws/aws-sdk-go-v2/internal/endpoints/v2 v2.4.31 // indirect
	github.com/aws/aws-sdk-go-v2/internal/ini v1.3.38 // indirect
	github.com/aws/aws-sdk-go-v2/internal/v4a v1.1.0 // indirect
	github.com/aws/aws-sdk-go-v2/service/internal/accept-encoding v1.9.11 // indirect
	github.com/aws/aws-sdk-go-v2/service/internal/checksum v1.1.29 // indirect
	github.com/aws/aws-sdk-go-v2/service/internal/endpoint-discovery v1.7.28 // indirect
	github.com/aws/aws-sdk-go-v2/service/internal/presigned-url v1.9.31 // indirect
	github.com/aws/aws-sdk-go-v2/service/internal/s3shared v1.14.3 // indirect
	github.com/aws/aws-sdk-go-v2/service/sso v1.13.1 // indirect
	github.com/aws/aws-sdk-go-v2/service/ssooidc v1.15.1 // indirect
	github.com/aws/aws-sdk-go-v2/service/sts v1.21.1
	github.com/aws/smithy-go v1.14.0 // indirect
	github.com/containerd/console v1.0.3 // indirect
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/go-chi/chi/v5 v5.0.8
	github.com/hashicorp/errwrap v1.0.0 // indirect
	github.com/jackc/pgpassfile v1.0.0 // indirect
	github.com/jackc/pgservicefile v0.0.0-20221227161230-091c0ba34f0a // indirect
	github.com/jmespath/go-jmespath v0.4.0 // indirect
	github.com/mattn/go-colorable v0.1.13 // indirect
	github.com/mattn/go-isatty v0.0.17 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	github.com/rogpeppe/go-internal v1.10.0 // indirect
	golang.org/x/crypto v0.9.0 // indirect
	golang.org/x/image v0.5.0 // indirect
	golang.org/x/sys v0.8.0 // indirect
	golang.org/x/text v0.9.0 // indirect
	golang.org/x/xerrors v0.0.0-20220907171357-04be3eba64a2 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)

retract (
	v1.0.1 // ignores Tamarin release
	v1.0.0 // ignores Tamarin release
)

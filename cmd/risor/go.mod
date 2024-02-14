module github.com/risor-io/risor/cmd/risor

go 1.22

replace (
	github.com/risor-io/risor => ../..
	github.com/risor-io/risor/modules/aws => ../../modules/aws
	github.com/risor-io/risor/modules/cli => ../../modules/cli
	github.com/risor-io/risor/modules/gha => ../../modules/gha
	github.com/risor-io/risor/modules/image => ../../modules/image
	github.com/risor-io/risor/modules/jmespath => ../../modules/jmespath
	github.com/risor-io/risor/modules/kubernetes => ../../modules/kubernetes
	github.com/risor-io/risor/modules/pgx => ../../modules/pgx
	github.com/risor-io/risor/modules/sql => ../../modules/sql
	github.com/risor-io/risor/modules/template => ../../modules/template
	github.com/risor-io/risor/modules/uuid => ../../modules/uuid
	github.com/risor-io/risor/modules/vault => ../../modules/vault
	github.com/risor-io/risor/os/s3fs => ../../os/s3fs
)

require (
	atomicgo.dev/keyboard v0.2.9
	github.com/aws/aws-sdk-go-v2/config v1.18.39
	github.com/aws/aws-sdk-go-v2/service/s3 v1.38.5
	github.com/fatih/color v1.15.0
	github.com/hokaccha/go-prettyjson v0.0.0-20211117102719-0474bc63780f
	github.com/mattn/go-isatty v0.0.20
	github.com/mitchellh/go-homedir v1.1.0
	github.com/risor-io/risor v1.3.2
	github.com/risor-io/risor/modules/aws v1.1.1
	github.com/risor-io/risor/modules/cli v0.0.0-00010101000000-000000000000
	github.com/risor-io/risor/modules/gha v0.0.0-20240213105055-b1d3a53935e5
	github.com/risor-io/risor/modules/image v1.1.1
	github.com/risor-io/risor/modules/jmespath v0.0.0-00010101000000-000000000000
	github.com/risor-io/risor/modules/kubernetes v0.0.0-00010101000000-000000000000
	github.com/risor-io/risor/modules/pgx v1.1.1
	github.com/risor-io/risor/modules/sql v0.0.0-00010101000000-000000000000
	github.com/risor-io/risor/modules/template v0.0.0-00010101000000-000000000000
	github.com/risor-io/risor/modules/uuid v1.1.1
	github.com/risor-io/risor/modules/vault v0.0.0-00010101000000-000000000000
	github.com/risor-io/risor/os/s3fs v1.1.1
	github.com/spf13/cobra v1.7.0
	github.com/spf13/viper v1.16.0
)

require (
	github.com/Masterminds/goutils v1.1.1 // indirect
	github.com/Masterminds/semver/v3 v3.2.0 // indirect
	github.com/Masterminds/sprig/v3 v3.2.3 // indirect
	github.com/anthonynsimon/bild v0.13.0 // indirect
	github.com/aws/aws-sdk-go-v2 v1.21.0 // indirect
	github.com/aws/aws-sdk-go-v2/aws/protocol/eventstream v1.4.13 // indirect
	github.com/aws/aws-sdk-go-v2/credentials v1.13.37 // indirect
	github.com/aws/aws-sdk-go-v2/feature/ec2/imds v1.13.11 // indirect
	github.com/aws/aws-sdk-go-v2/internal/configsources v1.1.41 // indirect
	github.com/aws/aws-sdk-go-v2/internal/endpoints/v2 v2.4.35 // indirect
	github.com/aws/aws-sdk-go-v2/internal/ini v1.3.42 // indirect
	github.com/aws/aws-sdk-go-v2/internal/v4a v1.1.4 // indirect
	github.com/aws/aws-sdk-go-v2/service/apigatewayv2 v1.14.5 // indirect
	github.com/aws/aws-sdk-go-v2/service/athena v1.31.6 // indirect
	github.com/aws/aws-sdk-go-v2/service/backup v1.25.0 // indirect
	github.com/aws/aws-sdk-go-v2/service/cloudformation v1.34.6 // indirect
	github.com/aws/aws-sdk-go-v2/service/cloudfront v1.28.5 // indirect
	github.com/aws/aws-sdk-go-v2/service/cloudtrail v1.29.0 // indirect
	github.com/aws/aws-sdk-go-v2/service/cloudwatch v1.27.7 // indirect
	github.com/aws/aws-sdk-go-v2/service/cloudwatchlogs v1.23.5 // indirect
	github.com/aws/aws-sdk-go-v2/service/dynamodb v1.21.5 // indirect
	github.com/aws/aws-sdk-go-v2/service/ebs v1.18.5 // indirect
	github.com/aws/aws-sdk-go-v2/service/ec2 v1.118.0 // indirect
	github.com/aws/aws-sdk-go-v2/service/ecr v1.20.0 // indirect
	github.com/aws/aws-sdk-go-v2/service/ecs v1.30.1 // indirect
	github.com/aws/aws-sdk-go-v2/service/eks v1.29.5 // indirect
	github.com/aws/aws-sdk-go-v2/service/elasticache v1.29.3 // indirect
	github.com/aws/aws-sdk-go-v2/service/elasticsearchservice v1.20.6 // indirect
	github.com/aws/aws-sdk-go-v2/service/eventbridge v1.22.0 // indirect
	github.com/aws/aws-sdk-go-v2/service/firehose v1.18.0 // indirect
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
	github.com/aws/aws-sdk-go-v2/service/rds v1.54.0 // indirect
	github.com/aws/aws-sdk-go-v2/service/redshift v1.29.5 // indirect
	github.com/aws/aws-sdk-go-v2/service/route53 v1.29.5 // indirect
	github.com/aws/aws-sdk-go-v2/service/secretsmanager v1.21.3 // indirect
	github.com/aws/aws-sdk-go-v2/service/sesv2 v1.20.0 // indirect
	github.com/aws/aws-sdk-go-v2/service/sfn v1.19.5 // indirect
	github.com/aws/aws-sdk-go-v2/service/sns v1.21.5 // indirect
	github.com/aws/aws-sdk-go-v2/service/sqs v1.24.5 // indirect
	github.com/aws/aws-sdk-go-v2/service/sso v1.13.6 // indirect
	github.com/aws/aws-sdk-go-v2/service/ssooidc v1.15.6 // indirect
	github.com/aws/aws-sdk-go-v2/service/sts v1.21.5 // indirect
	github.com/aws/aws-sdk-go-v2/service/wafv2 v1.38.0 // indirect
	github.com/aws/aws-sdk-go-v2/service/xray v1.18.0 // indirect
	github.com/aws/smithy-go v1.14.2 // indirect
	github.com/containerd/console v1.0.3 // indirect
	github.com/cpuguy83/go-md2man/v2 v2.0.3 // indirect
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/emicklei/go-restful/v3 v3.9.0 // indirect
	github.com/evanphx/json-patch/v5 v5.6.0 // indirect
	github.com/fsnotify/fsnotify v1.6.0 // indirect
	github.com/go-logr/logr v1.2.3 // indirect
	github.com/go-openapi/jsonpointer v0.19.5 // indirect
	github.com/go-openapi/jsonreference v0.20.0 // indirect
	github.com/go-openapi/swag v0.19.14 // indirect
	github.com/go-sql-driver/mysql v1.7.0 // indirect
	github.com/gofrs/uuid v4.4.0+incompatible // indirect
	github.com/gogo/protobuf v1.3.2 // indirect
	github.com/golang-sql/civil v0.0.0-20220223132316-b832511892a9 // indirect
	github.com/golang-sql/sqlexp v0.1.0 // indirect
	github.com/golang/protobuf v1.5.3 // indirect
	github.com/google/gnostic v0.5.7-v3refs // indirect
	github.com/google/go-cmp v0.6.0 // indirect
	github.com/google/gofuzz v1.1.0 // indirect
	github.com/google/uuid v1.3.0 // indirect
	github.com/hashicorp/go-cleanhttp v0.5.2 // indirect
	github.com/hashicorp/go-retryablehttp v0.7.1 // indirect
	github.com/hashicorp/go-rootcerts v1.0.2 // indirect
	github.com/hashicorp/go-secure-stdlib/strutil v0.1.2 // indirect
	github.com/hashicorp/hcl v1.0.0 // indirect
	github.com/hashicorp/vault-client-go v0.4.3 // indirect
	github.com/huandu/xstrings v1.3.3 // indirect
	github.com/imdario/mergo v0.3.11 // indirect
	github.com/inconshreveable/mousetrap v1.1.0 // indirect
	github.com/jackc/pgpassfile v1.0.0 // indirect
	github.com/jackc/pgservicefile v0.0.0-20221227161230-091c0ba34f0a // indirect
	github.com/jackc/pgx/v5 v5.5.0 // indirect
	github.com/jmespath-community/go-jmespath v1.1.1 // indirect
	github.com/jmespath/go-jmespath v0.4.0 // indirect
	github.com/josharian/intern v1.0.0 // indirect
	github.com/json-iterator/go v1.1.12 // indirect
	github.com/lib/pq v1.10.7 // indirect
	github.com/magiconair/properties v1.8.7 // indirect
	github.com/mailru/easyjson v0.7.6 // indirect
	github.com/mattn/go-colorable v0.1.13 // indirect
	github.com/microsoft/go-mssqldb v1.6.0 // indirect
	github.com/mitchellh/copystructure v1.0.0 // indirect
	github.com/mitchellh/mapstructure v1.5.0 // indirect
	github.com/mitchellh/reflectwalk v1.0.0 // indirect
	github.com/modern-go/concurrent v0.0.0-20180306012644-bacd9c7ef1dd // indirect
	github.com/modern-go/reflect2 v1.0.2 // indirect
	github.com/munnerz/goautoneg v0.0.0-20191010083416-a7dc8b61c822 // indirect
	github.com/pelletier/go-toml/v2 v2.1.0 // indirect
	github.com/pkg/errors v0.9.1 // indirect
	github.com/russross/blackfriday/v2 v2.1.0 // indirect
	github.com/ryanuber/go-glob v1.0.0 // indirect
	github.com/shopspring/decimal v1.2.0 // indirect
	github.com/spf13/afero v1.9.5 // indirect
	github.com/spf13/cast v1.5.1 // indirect
	github.com/spf13/jwalterweatherman v1.1.0 // indirect
	github.com/spf13/pflag v1.0.5 // indirect
	github.com/subosito/gotenv v1.6.0 // indirect
	github.com/urfave/cli/v2 v2.27.1 // indirect
	github.com/xo/dburl v0.20.0 // indirect
	github.com/xrash/smetrics v0.0.0-20231213231151-1d8dd44e695e // indirect
	golang.org/x/crypto v0.18.0 // indirect
	golang.org/x/exp v0.0.0-20230314191032-db074128a8ec // indirect
	golang.org/x/image v0.14.0 // indirect
	golang.org/x/net v0.20.0 // indirect
	golang.org/x/oauth2 v0.7.0 // indirect
	golang.org/x/sys v0.16.0 // indirect
	golang.org/x/term v0.16.0 // indirect
	golang.org/x/text v0.14.0 // indirect
	golang.org/x/time v0.3.0 // indirect
	google.golang.org/appengine v1.6.7 // indirect
	google.golang.org/protobuf v1.30.0 // indirect
	gopkg.in/inf.v0 v0.9.1 // indirect
	gopkg.in/ini.v1 v1.67.0 // indirect
	gopkg.in/yaml.v2 v2.4.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
	k8s.io/api v0.26.1 // indirect
	k8s.io/apimachinery v0.26.1 // indirect
	k8s.io/client-go v0.26.1 // indirect
	k8s.io/klog/v2 v2.80.1 // indirect
	k8s.io/kube-openapi v0.0.0-20221012153701-172d655c2280 // indirect
	k8s.io/utils v0.0.0-20221128185143-99ec85e7a448 // indirect
	sigs.k8s.io/controller-runtime v0.14.2 // indirect
	sigs.k8s.io/json v0.0.0-20220713155537-f223a00ba0e2 // indirect
	sigs.k8s.io/structured-merge-diff/v4 v4.2.3 // indirect
	sigs.k8s.io/yaml v1.3.0 // indirect
)

module github.com/risor-io/risor/cmd/risor

go 1.22.1

replace (
	github.com/risor-io/risor => ../..
	github.com/risor-io/risor/modules/aws => ../../modules/aws
	github.com/risor-io/risor/modules/bcrypt => ../../modules/bcrypt
	github.com/risor-io/risor/modules/carbon => ../../modules/carbon
	github.com/risor-io/risor/modules/cli => ../../modules/cli
	github.com/risor-io/risor/modules/gha => ../../modules/gha
	github.com/risor-io/risor/modules/image => ../../modules/image
	github.com/risor-io/risor/modules/jmespath => ../../modules/jmespath
	github.com/risor-io/risor/modules/kubernetes => ../../modules/kubernetes
	github.com/risor-io/risor/modules/pgx => ../../modules/pgx
	github.com/risor-io/risor/modules/semver => ../../modules/semver
	github.com/risor-io/risor/modules/sql => ../../modules/sql
	github.com/risor-io/risor/modules/template => ../../modules/template
	github.com/risor-io/risor/modules/uuid => ../../modules/uuid
	github.com/risor-io/risor/modules/vault => ../../modules/vault
	github.com/risor-io/risor/os/s3fs => ../../os/s3fs
)

require (
	atomicgo.dev/keyboard v0.2.9
	github.com/aws/aws-sdk-go-v2/config v1.27.31
	github.com/aws/aws-sdk-go-v2/service/s3 v1.61.0
	github.com/fatih/color v1.17.0
	github.com/hokaccha/go-prettyjson v0.0.0-20211117102719-0474bc63780f
	github.com/mattn/go-isatty v0.0.20
	github.com/mitchellh/go-homedir v1.1.0
	github.com/risor-io/risor v1.6.0
	github.com/risor-io/risor/modules/aws v1.6.0
	github.com/risor-io/risor/modules/bcrypt v1.6.0
	github.com/risor-io/risor/modules/carbon v1.6.0
	github.com/risor-io/risor/modules/cli v1.6.0
	github.com/risor-io/risor/modules/gha v1.6.1-0.20240927135333-245e7b83abf4
	github.com/risor-io/risor/modules/image v1.6.0
	github.com/risor-io/risor/modules/jmespath v1.6.0
	github.com/risor-io/risor/modules/kubernetes v1.6.0
	github.com/risor-io/risor/modules/pgx v1.6.0
	github.com/risor-io/risor/modules/semver v0.0.0-00010101000000-000000000000
	github.com/risor-io/risor/modules/sql v1.6.0
	github.com/risor-io/risor/modules/template v1.6.0
	github.com/risor-io/risor/modules/uuid v1.6.0
	github.com/risor-io/risor/modules/vault v1.6.0
	github.com/risor-io/risor/os/s3fs v1.6.0
	github.com/spf13/cobra v1.8.1
	github.com/spf13/viper v1.19.0
	github.com/stretchr/testify v1.9.0
)

require (
	dario.cat/mergo v1.0.1 // indirect
	filippo.io/edwards25519 v1.1.0 // indirect
	github.com/Masterminds/goutils v1.1.1 // indirect
	github.com/Masterminds/semver/v3 v3.3.0 // indirect
	github.com/Masterminds/sprig/v3 v3.3.0 // indirect
	github.com/anthonynsimon/bild v0.14.0 // indirect
	github.com/aws/aws-sdk-go-v2 v1.30.4 // indirect
	github.com/aws/aws-sdk-go-v2/aws/protocol/eventstream v1.6.4 // indirect
	github.com/aws/aws-sdk-go-v2/credentials v1.17.30 // indirect
	github.com/aws/aws-sdk-go-v2/feature/ec2/imds v1.16.12 // indirect
	github.com/aws/aws-sdk-go-v2/internal/configsources v1.3.16 // indirect
	github.com/aws/aws-sdk-go-v2/internal/endpoints/v2 v2.6.16 // indirect
	github.com/aws/aws-sdk-go-v2/internal/ini v1.8.1 // indirect
	github.com/aws/aws-sdk-go-v2/internal/v4a v1.3.16 // indirect
	github.com/aws/aws-sdk-go-v2/service/apigatewayv2 v1.22.6 // indirect
	github.com/aws/aws-sdk-go-v2/service/athena v1.44.5 // indirect
	github.com/aws/aws-sdk-go-v2/service/backup v1.36.4 // indirect
	github.com/aws/aws-sdk-go-v2/service/cloudformation v1.53.5 // indirect
	github.com/aws/aws-sdk-go-v2/service/cloudfront v1.38.5 // indirect
	github.com/aws/aws-sdk-go-v2/service/cloudtrail v1.42.5 // indirect
	github.com/aws/aws-sdk-go-v2/service/cloudwatch v1.40.5 // indirect
	github.com/aws/aws-sdk-go-v2/service/cloudwatchlogs v1.37.5 // indirect
	github.com/aws/aws-sdk-go-v2/service/dynamodb v1.34.6 // indirect
	github.com/aws/aws-sdk-go-v2/service/ebs v1.25.5 // indirect
	github.com/aws/aws-sdk-go-v2/service/ec2 v1.177.0 // indirect
	github.com/aws/aws-sdk-go-v2/service/ecr v1.32.2 // indirect
	github.com/aws/aws-sdk-go-v2/service/ecs v1.45.2 // indirect
	github.com/aws/aws-sdk-go-v2/service/eks v1.48.2 // indirect
	github.com/aws/aws-sdk-go-v2/service/elasticache v1.40.7 // indirect
	github.com/aws/aws-sdk-go-v2/service/elasticsearchservice v1.30.5 // indirect
	github.com/aws/aws-sdk-go-v2/service/eventbridge v1.33.5 // indirect
	github.com/aws/aws-sdk-go-v2/service/firehose v1.32.2 // indirect
	github.com/aws/aws-sdk-go-v2/service/glue v1.95.0 // indirect
	github.com/aws/aws-sdk-go-v2/service/iam v1.35.0 // indirect
	github.com/aws/aws-sdk-go-v2/service/internal/accept-encoding v1.11.4 // indirect
	github.com/aws/aws-sdk-go-v2/service/internal/checksum v1.3.18 // indirect
	github.com/aws/aws-sdk-go-v2/service/internal/endpoint-discovery v1.9.17 // indirect
	github.com/aws/aws-sdk-go-v2/service/internal/presigned-url v1.11.18 // indirect
	github.com/aws/aws-sdk-go-v2/service/internal/s3shared v1.17.16 // indirect
	github.com/aws/aws-sdk-go-v2/service/kinesis v1.29.5 // indirect
	github.com/aws/aws-sdk-go-v2/service/kms v1.35.5 // indirect
	github.com/aws/aws-sdk-go-v2/service/lambda v1.58.1 // indirect
	github.com/aws/aws-sdk-go-v2/service/ram v1.27.5 // indirect
	github.com/aws/aws-sdk-go-v2/service/rds v1.82.2 // indirect
	github.com/aws/aws-sdk-go-v2/service/redshift v1.46.6 // indirect
	github.com/aws/aws-sdk-go-v2/service/route53 v1.43.0 // indirect
	github.com/aws/aws-sdk-go-v2/service/secretsmanager v1.32.6 // indirect
	github.com/aws/aws-sdk-go-v2/service/sesv2 v1.33.0 // indirect
	github.com/aws/aws-sdk-go-v2/service/sfn v1.31.0 // indirect
	github.com/aws/aws-sdk-go-v2/service/sns v1.31.5 // indirect
	github.com/aws/aws-sdk-go-v2/service/sqs v1.34.5 // indirect
	github.com/aws/aws-sdk-go-v2/service/sso v1.22.5 // indirect
	github.com/aws/aws-sdk-go-v2/service/ssooidc v1.26.5 // indirect
	github.com/aws/aws-sdk-go-v2/service/sts v1.30.5 // indirect
	github.com/aws/aws-sdk-go-v2/service/wafv2 v1.52.0 // indirect
	github.com/aws/aws-sdk-go-v2/service/xray v1.27.5 // indirect
	github.com/aws/smithy-go v1.20.4 // indirect
	github.com/blang/semver/v4 v4.0.0 // indirect
	github.com/containerd/console v1.0.4 // indirect
	github.com/cpuguy83/go-md2man/v2 v2.0.4 // indirect
	github.com/davecgh/go-spew v1.1.2-0.20180830191138-d8f796af33cc // indirect
	github.com/emicklei/go-restful/v3 v3.12.1 // indirect
	github.com/evanphx/json-patch/v5 v5.9.0 // indirect
	github.com/fsnotify/fsnotify v1.7.0 // indirect
	github.com/fxamacker/cbor/v2 v2.7.0 // indirect
	github.com/go-logr/logr v1.4.2 // indirect
	github.com/go-openapi/jsonpointer v0.21.0 // indirect
	github.com/go-openapi/jsonreference v0.21.0 // indirect
	github.com/go-openapi/swag v0.23.0 // indirect
	github.com/go-sql-driver/mysql v1.8.1 // indirect
	github.com/gofrs/uuid/v5 v5.3.0 // indirect
	github.com/gogo/protobuf v1.3.2 // indirect
	github.com/golang-module/carbon/v2 v2.3.12 // indirect
	github.com/golang-sql/civil v0.0.0-20220223132316-b832511892a9 // indirect
	github.com/golang-sql/sqlexp v0.1.0 // indirect
	github.com/golang/protobuf v1.5.4 // indirect
	github.com/google/gnostic-models v0.6.9-0.20230804172637-c7be7c783f49 // indirect
	github.com/google/go-cmp v0.6.0 // indirect
	github.com/google/gofuzz v1.2.0 // indirect
	github.com/google/uuid v1.6.0 // indirect
	github.com/hashicorp/go-cleanhttp v0.5.2 // indirect
	github.com/hashicorp/go-retryablehttp v0.7.7 // indirect
	github.com/hashicorp/go-rootcerts v1.0.2 // indirect
	github.com/hashicorp/go-secure-stdlib/strutil v0.1.2 // indirect
	github.com/hashicorp/hcl v1.0.0 // indirect
	github.com/hashicorp/vault-client-go v0.4.3 // indirect
	github.com/huandu/xstrings v1.5.0 // indirect
	github.com/imdario/mergo v0.3.16 // indirect
	github.com/inconshreveable/mousetrap v1.1.0 // indirect
	github.com/jackc/pgpassfile v1.0.0 // indirect
	github.com/jackc/pgservicefile v0.0.0-20240606120523-5a60cdf6a761 // indirect
	github.com/jackc/pgx/v5 v5.6.0 // indirect
	github.com/jmespath-community/go-jmespath v1.1.1 // indirect
	github.com/jmespath/go-jmespath v0.4.0 // indirect
	github.com/josharian/intern v1.0.0 // indirect
	github.com/json-iterator/go v1.1.12 // indirect
	github.com/lib/pq v1.10.9 // indirect
	github.com/magiconair/properties v1.8.7 // indirect
	github.com/mailru/easyjson v0.7.7 // indirect
	github.com/mattn/go-colorable v0.1.13 // indirect
	github.com/mattn/go-runewidth v0.0.16 // indirect
	github.com/microsoft/go-mssqldb v1.7.2 // indirect
	github.com/mitchellh/copystructure v1.2.0 // indirect
	github.com/mitchellh/mapstructure v1.5.0 // indirect
	github.com/mitchellh/reflectwalk v1.0.2 // indirect
	github.com/modern-go/concurrent v0.0.0-20180306012644-bacd9c7ef1dd // indirect
	github.com/modern-go/reflect2 v1.0.2 // indirect
	github.com/munnerz/goautoneg v0.0.0-20191010083416-a7dc8b61c822 // indirect
	github.com/olekukonko/tablewriter v0.0.5 // indirect
	github.com/pelletier/go-toml/v2 v2.2.3 // indirect
	github.com/pkg/errors v0.9.1 // indirect
	github.com/pmezard/go-difflib v1.0.1-0.20181226105442-5d4384ee4fb2 // indirect
	github.com/rivo/uniseg v0.4.7 // indirect
	github.com/russross/blackfriday/v2 v2.1.0 // indirect
	github.com/ryanuber/go-glob v1.0.0 // indirect
	github.com/sagikazarmark/locafero v0.6.0 // indirect
	github.com/sagikazarmark/slog-shim v0.1.0 // indirect
	github.com/shopspring/decimal v1.4.0 // indirect
	github.com/sourcegraph/conc v0.3.0 // indirect
	github.com/spf13/afero v1.11.0 // indirect
	github.com/spf13/cast v1.7.0 // indirect
	github.com/spf13/pflag v1.0.5 // indirect
	github.com/subosito/gotenv v1.6.0 // indirect
	github.com/urfave/cli/v2 v2.27.4 // indirect
	github.com/x448/float16 v0.8.4 // indirect
	github.com/xo/dburl v0.23.2 // indirect
	github.com/xrash/smetrics v0.0.0-20240521201337-686a1a2994c1 // indirect
	go.uber.org/multierr v1.11.0 // indirect
	golang.org/x/crypto v0.26.0 // indirect
	golang.org/x/exp v0.0.0-20240823005443-9b4947da3948 // indirect
	golang.org/x/image v0.19.0 // indirect
	golang.org/x/net v0.28.0 // indirect
	golang.org/x/oauth2 v0.22.0 // indirect
	golang.org/x/sys v0.24.0 // indirect
	golang.org/x/term v0.23.0 // indirect
	golang.org/x/text v0.17.0 // indirect
	golang.org/x/time v0.6.0 // indirect
	google.golang.org/protobuf v1.34.2 // indirect
	gopkg.in/inf.v0 v0.9.1 // indirect
	gopkg.in/ini.v1 v1.67.0 // indirect
	gopkg.in/yaml.v2 v2.4.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
	k8s.io/api v0.31.0 // indirect
	k8s.io/apimachinery v0.31.0 // indirect
	k8s.io/client-go v0.31.0 // indirect
	k8s.io/klog/v2 v2.130.1 // indirect
	k8s.io/kube-openapi v0.0.0-20240827152857-f7e401e7b4c2 // indirect
	k8s.io/utils v0.0.0-20240821151609-f90d01438635 // indirect
	sigs.k8s.io/controller-runtime v0.19.0 // indirect
	sigs.k8s.io/json v0.0.0-20221116044647-bc3834ca7abd // indirect
	sigs.k8s.io/structured-merge-diff/v4 v4.4.1 // indirect
	sigs.k8s.io/yaml v1.4.0 // indirect
)

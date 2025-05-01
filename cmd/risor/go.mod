module github.com/risor-io/risor/cmd/risor

go 1.24.0

replace (
	github.com/risor-io/risor => ../..
	github.com/risor-io/risor/modules/aws => ../../modules/aws
	github.com/risor-io/risor/modules/bcrypt => ../../modules/bcrypt
	github.com/risor-io/risor/modules/carbon => ../../modules/carbon
	github.com/risor-io/risor/modules/cli => ../../modules/cli
	github.com/risor-io/risor/modules/gha => ../../modules/gha
	github.com/risor-io/risor/modules/goquery => ../../modules/goquery
	github.com/risor-io/risor/modules/htmltomarkdown => ../../modules/htmltomarkdown
	github.com/risor-io/risor/modules/image => ../../modules/image
	github.com/risor-io/risor/modules/jmespath => ../../modules/jmespath
	github.com/risor-io/risor/modules/kubernetes => ../../modules/kubernetes
	github.com/risor-io/risor/modules/pgx => ../../modules/pgx
	github.com/risor-io/risor/modules/playwright => ../../modules/playwright
	github.com/risor-io/risor/modules/qrcode => ../../modules/qrcode
	github.com/risor-io/risor/modules/sched => ../../modules/sched
	github.com/risor-io/risor/modules/semver => ../../modules/semver
	github.com/risor-io/risor/modules/shlex => ../../modules/shlex
	github.com/risor-io/risor/modules/slack => ../../modules/slack
	github.com/risor-io/risor/modules/sql => ../../modules/sql
	github.com/risor-io/risor/modules/template => ../../modules/template
	github.com/risor-io/risor/modules/uuid => ../../modules/uuid
	github.com/risor-io/risor/modules/vault => ../../modules/vault
	github.com/risor-io/risor/os/s3fs => ../../os/s3fs
)

require (
	atomicgo.dev/keyboard v0.2.9
	github.com/aws/aws-sdk-go-v2/config v1.29.14
	github.com/aws/aws-sdk-go-v2/service/s3 v1.79.2
	github.com/fatih/color v1.18.0
	github.com/hokaccha/go-prettyjson v0.0.0-20211117102719-0474bc63780f
	github.com/mattn/go-isatty v0.0.20
	github.com/mitchellh/go-homedir v1.1.0
	github.com/risor-io/risor v1.7.0
	github.com/risor-io/risor/modules/aws v1.7.0
	github.com/risor-io/risor/modules/bcrypt v1.7.0
	github.com/risor-io/risor/modules/carbon v1.7.0
	github.com/risor-io/risor/modules/cli v1.7.0
	github.com/risor-io/risor/modules/gha v1.7.0
	github.com/risor-io/risor/modules/goquery v0.0.0-00010101000000-000000000000
	github.com/risor-io/risor/modules/htmltomarkdown v0.0.0-00010101000000-000000000000
	github.com/risor-io/risor/modules/image v1.7.0
	github.com/risor-io/risor/modules/jmespath v1.7.0
	github.com/risor-io/risor/modules/kubernetes v1.7.0
	github.com/risor-io/risor/modules/pgx v1.7.0
	github.com/risor-io/risor/modules/playwright v0.0.0-00010101000000-000000000000
	github.com/risor-io/risor/modules/qrcode v0.0.0-00010101000000-000000000000
	github.com/risor-io/risor/modules/sched v1.7.0
	github.com/risor-io/risor/modules/semver v1.7.0
	github.com/risor-io/risor/modules/shlex v1.7.0
	github.com/risor-io/risor/modules/slack v0.0.0-00010101000000-000000000000
	github.com/risor-io/risor/modules/sql v1.7.0
	github.com/risor-io/risor/modules/template v1.7.0
	github.com/risor-io/risor/modules/uuid v1.7.0
	github.com/risor-io/risor/modules/vault v1.7.0
	github.com/risor-io/risor/os/s3fs v1.7.0
	github.com/spf13/cobra v1.9.1
	github.com/spf13/viper v1.20.1
	github.com/stretchr/testify v1.10.0
)

require (
	codnect.io/chrono v1.1.3 // indirect
	dario.cat/mergo v1.0.1 // indirect
	filippo.io/edwards25519 v1.1.0 // indirect
	github.com/JohannesKaufmann/dom v0.2.0 // indirect
	github.com/JohannesKaufmann/html-to-markdown/v2 v2.3.1 // indirect
	github.com/Masterminds/goutils v1.1.1 // indirect
	github.com/Masterminds/semver/v3 v3.3.1 // indirect
	github.com/Masterminds/sprig/v3 v3.3.0 // indirect
	github.com/PuerkitoBio/goquery v1.10.3 // indirect
	github.com/andybalholm/cascadia v1.3.3 // indirect
	github.com/anthonynsimon/bild v0.14.0 // indirect
	github.com/aws/aws-sdk-go-v2 v1.36.3 // indirect
	github.com/aws/aws-sdk-go-v2/aws/protocol/eventstream v1.6.10 // indirect
	github.com/aws/aws-sdk-go-v2/credentials v1.17.67 // indirect
	github.com/aws/aws-sdk-go-v2/feature/ec2/imds v1.16.30 // indirect
	github.com/aws/aws-sdk-go-v2/internal/configsources v1.3.34 // indirect
	github.com/aws/aws-sdk-go-v2/internal/endpoints/v2 v2.6.34 // indirect
	github.com/aws/aws-sdk-go-v2/internal/ini v1.8.3 // indirect
	github.com/aws/aws-sdk-go-v2/internal/v4a v1.3.34 // indirect
	github.com/aws/aws-sdk-go-v2/service/apigatewayv2 v1.27.1 // indirect
	github.com/aws/aws-sdk-go-v2/service/athena v1.50.4 // indirect
	github.com/aws/aws-sdk-go-v2/service/backup v1.41.2 // indirect
	github.com/aws/aws-sdk-go-v2/service/cloudformation v1.59.2 // indirect
	github.com/aws/aws-sdk-go-v2/service/cloudfront v1.45.3 // indirect
	github.com/aws/aws-sdk-go-v2/service/cloudtrail v1.48.4 // indirect
	github.com/aws/aws-sdk-go-v2/service/cloudwatch v1.44.3 // indirect
	github.com/aws/aws-sdk-go-v2/service/cloudwatchlogs v1.47.3 // indirect
	github.com/aws/aws-sdk-go-v2/service/dynamodb v1.43.0 // indirect
	github.com/aws/aws-sdk-go-v2/service/ebs v1.28.3 // indirect
	github.com/aws/aws-sdk-go-v2/service/ec2 v1.212.0 // indirect
	github.com/aws/aws-sdk-go-v2/service/ecr v1.43.3 // indirect
	github.com/aws/aws-sdk-go-v2/service/ecs v1.56.1 // indirect
	github.com/aws/aws-sdk-go-v2/service/eks v1.64.0 // indirect
	github.com/aws/aws-sdk-go-v2/service/elasticache v1.46.0 // indirect
	github.com/aws/aws-sdk-go-v2/service/elasticsearchservice v1.33.3 // indirect
	github.com/aws/aws-sdk-go-v2/service/eventbridge v1.39.0 // indirect
	github.com/aws/aws-sdk-go-v2/service/firehose v1.37.4 // indirect
	github.com/aws/aws-sdk-go-v2/service/glue v1.109.2 // indirect
	github.com/aws/aws-sdk-go-v2/service/iam v1.41.1 // indirect
	github.com/aws/aws-sdk-go-v2/service/internal/accept-encoding v1.12.3 // indirect
	github.com/aws/aws-sdk-go-v2/service/internal/checksum v1.7.0 // indirect
	github.com/aws/aws-sdk-go-v2/service/internal/endpoint-discovery v1.10.15 // indirect
	github.com/aws/aws-sdk-go-v2/service/internal/presigned-url v1.12.15 // indirect
	github.com/aws/aws-sdk-go-v2/service/internal/s3shared v1.18.15 // indirect
	github.com/aws/aws-sdk-go-v2/service/kinesis v1.33.3 // indirect
	github.com/aws/aws-sdk-go-v2/service/kms v1.38.3 // indirect
	github.com/aws/aws-sdk-go-v2/service/lambda v1.71.2 // indirect
	github.com/aws/aws-sdk-go-v2/service/ram v1.30.3 // indirect
	github.com/aws/aws-sdk-go-v2/service/rds v1.95.0 // indirect
	github.com/aws/aws-sdk-go-v2/service/redshift v1.54.3 // indirect
	github.com/aws/aws-sdk-go-v2/service/route53 v1.51.1 // indirect
	github.com/aws/aws-sdk-go-v2/service/secretsmanager v1.35.4 // indirect
	github.com/aws/aws-sdk-go-v2/service/sesv2 v1.45.0 // indirect
	github.com/aws/aws-sdk-go-v2/service/sfn v1.35.4 // indirect
	github.com/aws/aws-sdk-go-v2/service/sns v1.34.4 // indirect
	github.com/aws/aws-sdk-go-v2/service/sqs v1.38.5 // indirect
	github.com/aws/aws-sdk-go-v2/service/sso v1.25.3 // indirect
	github.com/aws/aws-sdk-go-v2/service/ssooidc v1.30.1 // indirect
	github.com/aws/aws-sdk-go-v2/service/sts v1.33.19 // indirect
	github.com/aws/aws-sdk-go-v2/service/wafv2 v1.60.1 // indirect
	github.com/aws/aws-sdk-go-v2/service/xray v1.31.3 // indirect
	github.com/aws/smithy-go v1.22.3 // indirect
	github.com/blang/semver/v4 v4.0.0 // indirect
	github.com/containerd/console v1.0.4 // indirect
	github.com/cpuguy83/go-md2man/v2 v2.0.6 // indirect
	github.com/davecgh/go-spew v1.1.2-0.20180830191138-d8f796af33cc // indirect
	github.com/deckarep/golang-set/v2 v2.8.0 // indirect
	github.com/emicklei/go-restful/v3 v3.12.2 // indirect
	github.com/evanphx/json-patch/v5 v5.9.11 // indirect
	github.com/fogleman/gg v1.3.0 // indirect
	github.com/fsnotify/fsnotify v1.9.0 // indirect
	github.com/fxamacker/cbor/v2 v2.8.0 // indirect
	github.com/go-jose/go-jose/v3 v3.0.4 // indirect
	github.com/go-logr/logr v1.4.2 // indirect
	github.com/go-openapi/jsonpointer v0.21.1 // indirect
	github.com/go-openapi/jsonreference v0.21.0 // indirect
	github.com/go-openapi/swag v0.23.1 // indirect
	github.com/go-sql-driver/mysql v1.9.2 // indirect
	github.com/go-stack/stack v1.8.1 // indirect
	github.com/go-viper/mapstructure/v2 v2.2.1 // indirect
	github.com/gofrs/uuid/v5 v5.3.2 // indirect
	github.com/gogo/protobuf v1.3.2 // indirect
	github.com/golang-module/carbon/v2 v2.3.12 // indirect
	github.com/golang-sql/civil v0.0.0-20220223132316-b832511892a9 // indirect
	github.com/golang-sql/sqlexp v0.1.0 // indirect
	github.com/golang/freetype v0.0.0-20170609003504-e2365dfdc4a0 // indirect
	github.com/google/gnostic-models v0.6.9 // indirect
	github.com/google/go-cmp v0.7.0 // indirect
	github.com/google/uuid v1.6.0 // indirect
	github.com/gorilla/websocket v1.5.4-0.20250319132907-e064f32e3674 // indirect
	github.com/hashicorp/go-cleanhttp v0.5.2 // indirect
	github.com/hashicorp/go-retryablehttp v0.7.7 // indirect
	github.com/hashicorp/go-rootcerts v1.0.2 // indirect
	github.com/hashicorp/go-secure-stdlib/strutil v0.1.2 // indirect
	github.com/hashicorp/vault-client-go v0.4.3 // indirect
	github.com/huandu/xstrings v1.5.0 // indirect
	github.com/inconshreveable/mousetrap v1.1.0 // indirect
	github.com/jackc/pgpassfile v1.0.0 // indirect
	github.com/jackc/pgservicefile v0.0.0-20240606120523-5a60cdf6a761 // indirect
	github.com/jackc/pgx/v5 v5.7.4 // indirect
	github.com/jmespath-community/go-jmespath v1.1.1 // indirect
	github.com/josharian/intern v1.0.0 // indirect
	github.com/json-iterator/go v1.1.12 // indirect
	github.com/lib/pq v1.10.9 // indirect
	github.com/mailru/easyjson v0.9.0 // indirect
	github.com/mattn/go-colorable v0.1.14 // indirect
	github.com/mattn/go-runewidth v0.0.16 // indirect
	github.com/mattn/go-sqlite3 v1.14.28 // indirect
	github.com/microsoft/go-mssqldb v1.8.0 // indirect
	github.com/mitchellh/copystructure v1.2.0 // indirect
	github.com/mitchellh/reflectwalk v1.0.2 // indirect
	github.com/modern-go/concurrent v0.0.0-20180306012644-bacd9c7ef1dd // indirect
	github.com/modern-go/reflect2 v1.0.2 // indirect
	github.com/munnerz/goautoneg v0.0.0-20191010083416-a7dc8b61c822 // indirect
	github.com/olekukonko/tablewriter v0.0.5 // indirect
	github.com/pelletier/go-toml/v2 v2.2.4 // indirect
	github.com/pkg/errors v0.9.1 // indirect
	github.com/playwright-community/playwright-go v0.5101.0 // indirect
	github.com/pmezard/go-difflib v1.0.1-0.20181226105442-5d4384ee4fb2 // indirect
	github.com/rivo/uniseg v0.4.7 // indirect
	github.com/russross/blackfriday/v2 v2.1.0 // indirect
	github.com/ryanuber/go-glob v1.0.0 // indirect
	github.com/sagikazarmark/locafero v0.9.0 // indirect
	github.com/shopspring/decimal v1.4.0 // indirect
	github.com/slack-go/slack v0.16.0 // indirect
	github.com/sourcegraph/conc v0.3.0 // indirect
	github.com/spf13/afero v1.14.0 // indirect
	github.com/spf13/cast v1.7.1 // indirect
	github.com/spf13/pflag v1.0.6 // indirect
	github.com/subosito/gotenv v1.6.0 // indirect
	github.com/u-root/u-root v0.14.0 // indirect
	github.com/urfave/cli/v2 v2.27.6 // indirect
	github.com/x448/float16 v0.8.4 // indirect
	github.com/xo/dburl v0.23.7 // indirect
	github.com/xrash/smetrics v0.0.0-20240521201337-686a1a2994c1 // indirect
	github.com/yeqown/go-qrcode/v2 v2.2.5 // indirect
	github.com/yeqown/go-qrcode/writer/standard v1.2.1 // indirect
	github.com/yeqown/reedsolomon v1.0.0 // indirect
	go.uber.org/multierr v1.11.0 // indirect
	golang.org/x/crypto v0.37.0 // indirect
	golang.org/x/exp v0.0.0-20250408133849-7e4ce0ab07d0 // indirect
	golang.org/x/image v0.26.0 // indirect
	golang.org/x/net v0.39.0 // indirect
	golang.org/x/oauth2 v0.29.0 // indirect
	golang.org/x/sys v0.32.0 // indirect
	golang.org/x/term v0.31.0 // indirect
	golang.org/x/text v0.24.0 // indirect
	golang.org/x/time v0.11.0 // indirect
	google.golang.org/protobuf v1.36.6 // indirect
	gopkg.in/inf.v0 v0.9.1 // indirect
	gopkg.in/yaml.v2 v2.4.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
	k8s.io/api v0.33.0 // indirect
	k8s.io/apimachinery v0.33.0 // indirect
	k8s.io/client-go v0.33.0 // indirect
	k8s.io/klog/v2 v2.130.1 // indirect
	k8s.io/kube-openapi v0.0.0-20250318190949-c8a335a9a2ff // indirect
	k8s.io/utils v0.0.0-20250321185631-1f6e0b77f77e // indirect
	sigs.k8s.io/controller-runtime v0.20.4 // indirect
	sigs.k8s.io/json v0.0.0-20241014173422-cfa47c3a1cc8 // indirect
	sigs.k8s.io/randfill v1.0.0 // indirect
	sigs.k8s.io/structured-merge-diff/v4 v4.7.0 // indirect
	sigs.k8s.io/yaml v1.4.0 // indirect
)

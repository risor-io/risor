//go:build aws
// +build aws

package aws

import (
	"fmt"
	"reflect"
	"unicode"

	"github.com/aws/aws-sdk-go-v2/service/apigatewayv2"
	"github.com/aws/aws-sdk-go-v2/service/athena"
	"github.com/aws/aws-sdk-go-v2/service/backup"
	"github.com/aws/aws-sdk-go-v2/service/cloudformation"
	"github.com/aws/aws-sdk-go-v2/service/cloudfront"
	"github.com/aws/aws-sdk-go-v2/service/cloudtrail"
	"github.com/aws/aws-sdk-go-v2/service/cloudwatch"
	"github.com/aws/aws-sdk-go-v2/service/cloudwatchlogs"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/ebs"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/ecr"
	"github.com/aws/aws-sdk-go-v2/service/ecs"
	"github.com/aws/aws-sdk-go-v2/service/eks"
	"github.com/aws/aws-sdk-go-v2/service/elasticache"
	"github.com/aws/aws-sdk-go-v2/service/elasticsearchservice"
	"github.com/aws/aws-sdk-go-v2/service/eventbridge"
	"github.com/aws/aws-sdk-go-v2/service/firehose"
	"github.com/aws/aws-sdk-go-v2/service/glue"
	"github.com/aws/aws-sdk-go-v2/service/iam"
	"github.com/aws/aws-sdk-go-v2/service/kinesis"
	"github.com/aws/aws-sdk-go-v2/service/kms"
	"github.com/aws/aws-sdk-go-v2/service/lambda"
	"github.com/aws/aws-sdk-go-v2/service/ram"
	"github.com/aws/aws-sdk-go-v2/service/rds"
	"github.com/aws/aws-sdk-go-v2/service/redshift"
	"github.com/aws/aws-sdk-go-v2/service/route53"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/secretsmanager"
	"github.com/aws/aws-sdk-go-v2/service/sesv2"
	"github.com/aws/aws-sdk-go-v2/service/sfn"
	"github.com/aws/aws-sdk-go-v2/service/sns"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/aws/aws-sdk-go-v2/service/sts"
	"github.com/aws/aws-sdk-go-v2/service/wafv2"
	"github.com/aws/aws-sdk-go-v2/service/xray"
	"github.com/risor-io/risor/errz"
	"github.com/risor-io/risor/object"
	"github.com/risor-io/risor/op"
)

type Client struct {
	client  interface{}
	service string
	methods map[string]*GoMethod
	config  *Config
}

func (c *Client) Inspect() string {
	return fmt.Sprintf("aws.client(service=%s, region=%s)",
		c.service, c.config.Region())
}

func (c *Client) Type() object.Type {
	return "aws.client"
}

func (c *Client) Value() interface{} {
	return c.client
}

func (c *Client) GetAttr(name string) (object.Object, bool) {
	switch name {
	case "__config__":
		return c.config, true
	case "__service__":
		return object.NewString(c.service), true
	case "__region__":
		return object.NewString(c.config.Region()), true
	case "__api_methods__":
		var names []object.Object
		for name := range c.methods {
			names = append(names, object.NewString(name))
		}
		return object.NewList(names), true
	}
	method, ok := c.methods[name]
	if !ok {
		return nil, false
	}
	methodName := fmt.Sprintf("aws.%s.%s", c.service, method.Name)
	return NewMethod(methodName, c.client, method), true
}

func (c *Client) SetAttr(name string, value object.Object) error {
	return fmt.Errorf("attribute error: aws.client object has no attribute %q", name)
}

func (c *Client) Interface() interface{} {
	return c.client
}

func (c *Client) String() string {
	return c.Inspect()
}

func (c *Client) Compare(other object.Object) (int, error) {
	return 0, errz.TypeErrorf("type error: unable to compare aws.client")
}

func (c *Client) Equals(other object.Object) object.Object {
	if c == other {
		return object.True
	}
	return object.False
}

func (c *Client) IsTruthy() bool {
	return true
}

func (c *Client) RunOperation(opType op.BinaryOpType, right object.Object) object.Object {
	return object.EvalErrorf("eval error: unsupported operation for aws.client: %v ", opType)
}

func (c *Client) Cost() int {
	return 0
}

func (c *Client) MarshalJSON() ([]byte, error) {
	return nil, errz.TypeErrorf("type error: unable to marshal aws.client")
}

func NewClient(service string, client interface{}, config *Config) *Client {
	return &Client{
		service: service,
		client:  client,
		methods: loadMethods(client),
		config:  config,
	}
}

func getClient(service string, cfg *Config) object.Object {
	switch service {
	case "ec2":
		return NewClient("ec2", ec2.NewFromConfig(cfg.value), cfg)
	case "rds":
		return NewClient("rds", rds.NewFromConfig(cfg.value), cfg)
	case "s3":
		return NewClient("s3", s3.NewFromConfig(cfg.value), cfg)
	case "iam":
		return NewClient("iam", iam.NewFromConfig(cfg.value), cfg)
	case "ebs":
		return NewClient("ebs", ebs.NewFromConfig(cfg.value), cfg)
	case "lambda":
		return NewClient("lambda", lambda.NewFromConfig(cfg.value), cfg)
	case "cloudfront":
		return NewClient("cloudfront", cloudfront.NewFromConfig(cfg.value), cfg)
	case "sns":
		return NewClient("sns", sns.NewFromConfig(cfg.value), cfg)
	case "sqs":
		return NewClient("sqs", sqs.NewFromConfig(cfg.value), cfg)
	case "ddb":
		return NewClient("ddb", dynamodb.NewFromConfig(cfg.value), cfg)
	case "elasticache":
		return NewClient("elasticache", elasticache.NewFromConfig(cfg.value), cfg)
	case "elasticsearchservice":
		return NewClient("elasticsearchservice", elasticsearchservice.NewFromConfig(cfg.value), cfg)
	case "cloudwatch":
		return NewClient("cloudwatch", cloudwatch.NewFromConfig(cfg.value), cfg)
	case "kms":
		return NewClient("kms", kms.NewFromConfig(cfg.value), cfg)
	case "cloudformation":
		return NewClient("cloudformation", cloudformation.NewFromConfig(cfg.value), cfg)
	case "kinesis":
		return NewClient("kinesis", kinesis.NewFromConfig(cfg.value), cfg)
	case "route53":
		return NewClient("route53", route53.NewFromConfig(cfg.value), cfg)
	case "redshift":
		return NewClient("redshift", redshift.NewFromConfig(cfg.value), cfg)
	case "glue":
		return NewClient("glue", glue.NewFromConfig(cfg.value), cfg)
	case "cloudtrail":
		return NewClient("cloudtrail", cloudtrail.NewFromConfig(cfg.value), cfg)
	case "wafv2":
		return NewClient("wafv2", wafv2.NewFromConfig(cfg.value), cfg)
	case "sfn":
		return NewClient("sfn", sfn.NewFromConfig(cfg.value), cfg)
	case "eks":
		return NewClient("eks", eks.NewFromConfig(cfg.value), cfg)
	case "ecr":
		return NewClient("ecr", ecr.NewFromConfig(cfg.value), cfg)
	case "ecs":
		return NewClient("ecs", ecs.NewFromConfig(cfg.value), cfg)
	case "sts":
		return NewClient("sts", sts.NewFromConfig(cfg.value), cfg)
	case "apigatewayv2":
		return NewClient("apigatewayv2", apigatewayv2.NewFromConfig(cfg.value), cfg)
	case "athena":
		return NewClient("athena", athena.NewFromConfig(cfg.value), cfg)
	case "backup":
		return NewClient("backup", backup.NewFromConfig(cfg.value), cfg)
	case "cloudwatchlogs":
		return NewClient("cloudwatchlogs", cloudwatchlogs.NewFromConfig(cfg.value), cfg)
	case "eventbridge":
		return NewClient("eventbridge", eventbridge.NewFromConfig(cfg.value), cfg)
	case "ram":
		return NewClient("ram", ram.NewFromConfig(cfg.value), cfg)
	case "secretsmanager":
		return NewClient("secretsmanager", secretsmanager.NewFromConfig(cfg.value), cfg)
	case "sesv2":
		return NewClient("sesv2", sesv2.NewFromConfig(cfg.value), cfg)
	case "xray":
		return NewClient("xray", xray.NewFromConfig(cfg.value), cfg)
	case "firehose":
		return NewClient("firehose", firehose.NewFromConfig(cfg.value), cfg)
	default:
		return object.Errorf("unknown aws service: %s", service)
	}
}

func loadMethods(obj interface{}) map[string]*GoMethod {
	typ := reflect.TypeOf(obj)
	methods := make(map[string]*GoMethod, typ.NumMethod())
	for i := 0; i < typ.NumMethod(); i++ {
		m := typ.Method(i)
		name := toSnakeCase(m.Name)
		goMethod := &GoMethod{
			Method:     m,
			Name:       name,
			NumIn:      m.Type.NumIn(),
			NumOut:     m.Type.NumOut(),
			IsVariadic: m.Type.IsVariadic(),
		}
		for i := 0; i < goMethod.NumIn; i++ {
			goMethod.InTypes = append(goMethod.InTypes, m.Type.In(i))
		}
		for i := 0; i < goMethod.NumOut; i++ {
			goMethod.OutTypes = append(goMethod.OutTypes, m.Type.Out(i))
		}
		methods[name] = goMethod
	}
	return methods
}

type GoMethod struct {
	Name       string
	Method     reflect.Method
	NumIn      int
	NumOut     int
	InTypes    []reflect.Type
	OutTypes   []reflect.Type
	IsVariadic bool
}

func toSnakeCase(str string) string {
	var lastUpper bool
	var result string
	for i, v := range str {
		if unicode.IsUpper(v) {
			if i != 0 && !lastUpper {
				result += "_"
			}
			result += string(unicode.ToLower(v))
			lastUpper = true
		} else {
			result += string(v)
			lastUpper = false
		}
	}
	return result
}

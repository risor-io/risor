//go:build aws
// +build aws

package aws

import (
	"context"
	"errors"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/cloudformation"
	"github.com/aws/aws-sdk-go-v2/service/cloudfront"
	"github.com/aws/aws-sdk-go-v2/service/cloudtrail"
	"github.com/aws/aws-sdk-go-v2/service/cloudwatch"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/ebs"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/ecr"
	"github.com/aws/aws-sdk-go-v2/service/ecs"
	"github.com/aws/aws-sdk-go-v2/service/eks"
	"github.com/aws/aws-sdk-go-v2/service/elasticache"
	"github.com/aws/aws-sdk-go-v2/service/elasticsearchservice"
	"github.com/aws/aws-sdk-go-v2/service/glue"
	"github.com/aws/aws-sdk-go-v2/service/iam"
	"github.com/aws/aws-sdk-go-v2/service/kinesis"
	"github.com/aws/aws-sdk-go-v2/service/kms"
	"github.com/aws/aws-sdk-go-v2/service/lambda"
	"github.com/aws/aws-sdk-go-v2/service/rds"
	"github.com/aws/aws-sdk-go-v2/service/redshift"
	"github.com/aws/aws-sdk-go-v2/service/route53"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/sfn"
	"github.com/aws/aws-sdk-go-v2/service/sns"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/aws/aws-sdk-go-v2/service/wafv2"
	"github.com/cloudcmds/tamarin/v2/internal/arg"
	"github.com/cloudcmds/tamarin/v2/object"
	"github.com/cloudcmds/tamarin/v2/op"
)

type Config struct {
	cfg aws.Config
}

func (c *Config) Inspect() string {
	return "aws.config()"
}

func (c *Config) Type() object.Type {
	return "aws.config"
}

func (c *Config) Value() config.Config {
	return c.cfg
}

func (c *Config) GetAttr(name string) (object.Object, bool) {
	switch name {
	case "region":
		return object.NewString(c.cfg.Region), true
	case "with_region":
		return object.NewBuiltin("with_region", func(ctx context.Context, args ...object.Object) object.Object {
			if err := arg.Require("aws.config.with_region", 1, args); err != nil {
				return err
			}
			region, err := object.AsString(args[0])
			if err != nil {
				return err
			}
			new := c.cfg.Copy()
			new.Region = region
			return NewConfig(new)
		}), true
	case "copy":
		return object.NewBuiltin("copy", func(ctx context.Context, args ...object.Object) object.Object {
			return NewConfig(c.cfg.Copy())
		}), true
	case "credentials":
		return object.NewBuiltin("credentials", func(ctx context.Context, args ...object.Object) object.Object {
			creds, err := c.cfg.Credentials.Retrieve(ctx)
			if err != nil {
				return object.NewError(err)
			}
			credsMap := map[string]interface{}{
				"access_key_id":     creds.AccessKeyID,
				"secret_access_key": creds.SecretAccessKey,
				"session_token":     creds.SessionToken,
				"can_expire":        creds.CanExpire,
				"source":            creds.Source,
			}
			if creds.CanExpire {
				credsMap["expires"] = creds.Expires
			}
			return object.FromGoType(credsMap)
		}), true
	case "service":
		return object.NewBuiltin("service", func(ctx context.Context, args ...object.Object) object.Object {
			if err := arg.Require("aws.config.service", 1, args); err != nil {
				return err
			}
			service, err := object.AsString(args[0])
			if err != nil {
				return err
			}
			switch service {
			case "ec2":
				return NewClient("ec2", ec2.NewFromConfig(c.cfg))
			case "rds":
				return NewClient("rds", rds.NewFromConfig(c.cfg))
			case "s3":
				return NewClient("s3", s3.NewFromConfig(c.cfg))
			case "iam":
				return NewClient("iam", iam.NewFromConfig(c.cfg))
			case "ebs":
				return NewClient("ebs", ebs.NewFromConfig(c.cfg))
			case "lambda":
				return NewClient("lambda", lambda.NewFromConfig(c.cfg))
			case "cloudfront":
				return NewClient("cloudfront", cloudfront.NewFromConfig(c.cfg))
			case "sns":
				return NewClient("sns", sns.NewFromConfig(c.cfg))
			case "sqs":
				return NewClient("sqs", sqs.NewFromConfig(c.cfg))
			case "ddb":
				return NewClient("ddb", dynamodb.NewFromConfig(c.cfg))
			case "elasticache":
				return NewClient("elasticache", elasticache.NewFromConfig(c.cfg))
			case "elasticsearchservice":
				return NewClient("elasticsearchservice", elasticsearchservice.NewFromConfig(c.cfg))
			case "cloudwatch":
				return NewClient("cloudwatch", cloudwatch.NewFromConfig(c.cfg))
			case "kms":
				return NewClient("kms", kms.NewFromConfig(c.cfg))
			case "cloudformation":
				return NewClient("cloudformation", cloudformation.NewFromConfig(c.cfg))
			case "kinesis":
				return NewClient("kinesis", kinesis.NewFromConfig(c.cfg))
			case "route53":
				return NewClient("route53", route53.NewFromConfig(c.cfg))
			case "redshift":
				return NewClient("redshift", redshift.NewFromConfig(c.cfg))
			case "glue":
				return NewClient("glue", glue.NewFromConfig(c.cfg))
			case "cloudtrail":
				return NewClient("cloudtrail", cloudtrail.NewFromConfig(c.cfg))
			case "wafv2":
				return NewClient("wafv2", wafv2.NewFromConfig(c.cfg))
			case "sfn":
				return NewClient("sfn", sfn.NewFromConfig(c.cfg))
			case "eks":
				return NewClient("eks", eks.NewFromConfig(c.cfg))
			case "ecr":
				return NewClient("ecr", ecr.NewFromConfig(c.cfg))
			case "ecs":
				return NewClient("ecs", ecs.NewFromConfig(c.cfg))
			default:
				return object.Errorf("unknown service: %s", service)
			}
		}), true
	}
	return nil, false
}

func (c *Config) SetAttr(name string, value object.Object) error {
	return fmt.Errorf("attribute error: aws.config object has no attribute %q", name)
}

func (c *Config) Interface() interface{} {
	return c.cfg
}

func (c *Config) String() string {
	return c.Inspect()
}

func (c *Config) Compare(other object.Object) (int, error) {
	return 0, errors.New("type error: unable to compare aws.config")
}

func (c *Config) Equals(other object.Object) object.Object {
	if c == other {
		return object.True
	}
	return object.False
}

func (c *Config) IsTruthy() bool {
	return true
}

func (c *Config) RunOperation(opType op.BinaryOpType, right object.Object) object.Object {
	return object.NewError(fmt.Errorf("eval error: unsupported operation for aws.config: %v ", opType))
}

func (c *Config) Cost() int {
	return 0
}

func NewConfig(cfg aws.Config) *Config {
	return &Config{cfg: cfg}
}

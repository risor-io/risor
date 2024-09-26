//go:build aws
// +build aws

package aws

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/risor-io/risor/arg"
	"github.com/risor-io/risor/errz"
	"github.com/risor-io/risor/object"
	"github.com/risor-io/risor/op"
)

type Config struct {
	value aws.Config
}

func (c *Config) Inspect() string {
	return fmt.Sprintf("aws.config(region=%s)", c.value.Region)
}

func (c *Config) Type() object.Type {
	return "aws.config"
}

func (c *Config) Value() config.Config {
	return c.value
}

func (c *Config) GetAttr(name string) (object.Object, bool) {
	switch name {
	case "region":
		return object.NewString(c.value.Region), true
	case "with_region":
		return object.NewBuiltin("with_region", func(ctx context.Context, args ...object.Object) object.Object {
			if err := arg.Require("aws.config.with_region", 1, args); err != nil {
				return err
			}
			region, err := object.AsString(args[0])
			if err != nil {
				return err
			}
			new := c.value.Copy()
			new.Region = region
			return NewConfig(new)
		}), true
	case "copy":
		return object.NewBuiltin("copy", func(ctx context.Context, args ...object.Object) object.Object {
			return NewConfig(c.value.Copy())
		}), true
	case "credentials":
		return object.NewBuiltin("credentials", func(ctx context.Context, args ...object.Object) object.Object {
			creds, err := c.value.Credentials.Retrieve(ctx)
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
	}
	return nil, false
}

func (c *Config) SetAttr(name string, value object.Object) error {
	return fmt.Errorf("attribute error: aws.config object has no attribute %q", name)
}

func (c *Config) Interface() interface{} {
	return c.value
}

func (c *Config) String() string {
	return c.Inspect()
}

func (c *Config) Compare(other object.Object) (int, error) {
	return 0, errz.TypeErrorf("type error: unable to compare aws.config")
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
	return object.EvalErrorf("eval error: unsupported operation for aws.config: %v ", opType)
}

func (c *Config) Cost() int {
	return 0
}

func (c *Config) MarshalJSON() ([]byte, error) {
	return nil, errz.TypeErrorf("type error: unable to marshal aws.config")
}

func (c *Config) Region() string {
	return c.value.Region
}

func NewConfig(cfg aws.Config) *Config {
	return &Config{value: cfg}
}

func NewConfigFromMap(ctx context.Context, m *object.Map) (*Config, error) {
	opts, err := getConfigOptions(m)
	if err != nil {
		return nil, err
	}
	cfg, err := config.LoadDefaultConfig(ctx, opts...)
	if err != nil {
		return nil, err
	}
	return NewConfig(cfg), nil
}

func getConfigOptions(m *object.Map) ([]func(*config.LoadOptions) error, error) {
	// Options:
	// {
	//    "region": "us-east-1",
	//    "credentials": {
	//       "key": "AKID",
	//       "secret": "SECRET",
	//       "session": "SESSION_TOKEN"
	//    },
	//    "profile": "custom_profile",
	//    "credentials_files": ["test/credentials"],
	//    "config_files": ["test/config"],
	// }

	var opts []func(*config.LoadOptions) error

	// Optional region
	region, present, err := mapGetStr(m, "region")
	if err != nil {
		return nil, err
	} else if present {
		opts = append(opts, config.WithRegion(region))
	}

	// Optional credentials
	creds, present, err := mapGetMap(m, "credentials")
	if err != nil {
		return nil, err
	} else if present {
		key, _, err := mapGetStr(creds, "key")
		if err != nil {
			return nil, err
		}
		secret, _, err := mapGetStr(creds, "secret")
		if err != nil {
			return nil, err
		}
		session, _, err := mapGetStr(creds, "session")
		if err != nil {
			return nil, err
		}
		opts = append(opts, config.WithCredentialsProvider(
			credentials.NewStaticCredentialsProvider(key, secret, session)))
	}

	// Optional profile
	profile, present, err := mapGetStr(m, "profile")
	if err != nil {
		return nil, err
	} else if present {
		opts = append(opts, config.WithSharedConfigProfile(profile))
	}

	// Optional shared credentials files
	credentialsFiles, present, err := mapGetStrList(m, "credentials_files")
	if err != nil {
		return nil, err
	} else if present {
		opts = append(opts, config.WithSharedCredentialsFiles(credentialsFiles))
	}

	// Optional shared config files
	configFiles, present, err := mapGetStrList(m, "config_files")
	if err != nil {
		return nil, err
	} else if present {
		opts = append(opts, config.WithSharedConfigFiles(configFiles))
	}
	return opts, nil
}

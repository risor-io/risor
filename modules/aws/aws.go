//go:build aws
// +build aws

package aws

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/risor-io/risor/internal/arg"
	"github.com/risor-io/risor/object"
)

type configuration struct {
	Region      string
	Credentials struct {
		Key     string
		Secret  string
		Session string
	}
	Profile          string
	CredentialsFiles []string
	ConfigFiles      []string
}

func ConfigFunc(ctx context.Context, args ...object.Object) object.Object {
	if err := arg.RequireRange("aws.config", 0, 1, args); err != nil {
		return err
	}
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
	if len(args) == 1 {
		// Configuration options may be passed as a map
		m, err := object.AsMap(args[0])
		if err != nil {
			return err
		}
		// Optional region
		region, present, err := mapGetStr(m, "region")
		if err != nil {
			return err
		} else if present {
			opts = append(opts, config.WithRegion(region))
		}
		// Optional credentials
		creds, present, err := mapGetMap(m, "credentials")
		if err != nil {
			return err
		} else if present {
			key, _, err := mapGetStr(creds, "key")
			if err != nil {
				return err
			}
			secret, _, err := mapGetStr(creds, "secret")
			if err != nil {
				return err
			}
			session, _, err := mapGetStr(creds, "session")
			if err != nil {
				return err
			}
			opts = append(opts, config.WithCredentialsProvider(
				credentials.NewStaticCredentialsProvider(key, secret, session)))
		}
		// Optional profile
		profile, present, err := mapGetStr(m, "profile")
		if err != nil {
			return err
		} else if present {
			opts = append(opts, config.WithSharedConfigProfile(profile))
		}
		// Optional shared credentials files
		credentialsFiles, present, err := mapGetStrList(m, "credentials_files")
		if err != nil {
			return err
		} else if present {
			opts = append(opts, config.WithSharedCredentialsFiles(credentialsFiles))
		}
		// Optional shared config files
		configFiles, present, err := mapGetStrList(m, "config_files")
		if err != nil {
			return err
		} else if present {
			opts = append(opts, config.WithSharedConfigFiles(configFiles))
		}
	}
	cfg, err := config.LoadDefaultConfig(ctx, opts...)
	if err != nil {
		return object.NewError(err)
	}
	return NewConfig(cfg)
}

func Module() *object.Module {
	return object.NewBuiltinsModule("aws", map[string]object.Object{
		"config": object.NewBuiltin("aws.config", ConfigFunc),
	})
}

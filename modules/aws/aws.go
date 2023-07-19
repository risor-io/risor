//go:build aws
// +build aws

package aws

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/risor-io/risor/internal/arg"
	"github.com/risor-io/risor/object"
)

func ConfigFunc(ctx context.Context, args ...object.Object) object.Object {
	if err := arg.RequireRange("aws.config", 0, 1, args); err != nil {
		return err
	}
	var opts []func(*config.LoadOptions) error
	if len(args) == 1 {
		// Configuration options may be passed as a map
		m, err := object.AsMap(args[0])
		if err != nil {
			return err
		}
		cfg, initErr := NewConfigFromMap(ctx, m)
		if initErr != nil {
			return object.NewError(initErr)
		}
		return cfg
	}
	cfg, err := config.LoadDefaultConfig(ctx, opts...)
	if err != nil {
		return object.NewError(err)
	}
	return NewConfig(cfg)
}

func ClientFunc(ctx context.Context, args ...object.Object) object.Object {
	if err := arg.RequireRange("aws.client", 1, 2, args); err != nil {
		return err
	}
	serviceName, err := object.AsString(args[0])
	if err != nil {
		return err
	}
	var cfg *Config
	if len(args) == 2 {
		switch arg := args[1].(type) {
		case *object.Map:
			// Configuration options may be passed as a map
			m, err := object.AsMap(args[1])
			if err != nil {
				return err
			}
			var initErr error
			cfg, initErr = NewConfigFromMap(ctx, m)
			if initErr != nil {
				return object.NewError(initErr)
			}
		case *Config:
			cfg = arg
		default:
			return object.Errorf("aws.client: expected config or map (got %s)", args[1].Type())
		}
	} else {
		awsCfg, err := config.LoadDefaultConfig(ctx)
		if err != nil {
			return object.NewError(err)
		}
		cfg = NewConfig(awsCfg)
	}
	return getClient(serviceName, cfg)
}

func Module() *object.Module {
	return object.NewBuiltinsModule("aws", map[string]object.Object{
		"config": object.NewBuiltin("aws.config", ConfigFunc),
		"client": object.NewBuiltin("aws.client", ClientFunc),
	})
}

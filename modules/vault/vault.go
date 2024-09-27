//go:build vault
// +build vault

package vault

import (
	"context"
	"time"

	"github.com/hashicorp/vault-client-go"
	"github.com/hashicorp/vault-client-go/schema"
	"github.com/risor-io/risor/arg"
	"github.com/risor-io/risor/errz"
	"github.com/risor-io/risor/object"
	"github.com/risor-io/risor/op"
)

const VAULT object.Type = "vault"

type Vault struct {
	client *vault.Client
}

func (v *Vault) Type() object.Type {
	return VAULT
}

func (v *Vault) Inspect() string {
	return "vault.client"
}

func (v *Vault) Interface() interface{} {
	return v.client
}

func (v *Vault) IsTruthy() bool {
	return v.client != nil
}

func (v *Vault) Cost() int {
	return 8
}

func (v *Vault) MarshalJSON() ([]byte, error) {
	return nil, errz.TypeErrorf("type error: unable to marshal %s", VAULT)
}

func (v *Vault) RunOperation(opType op.BinaryOpType, right object.Object) object.Object {
	return object.TypeErrorf("type error: unsupported operation for %s: %v", VAULT, opType)
}

func (v *Vault) Equals(other object.Object) object.Object {
	if other.Type() != VAULT {
		return object.False
	}
	return object.NewBool(v.client == other.(*Vault).client)
}

func (v *Vault) SetAttr(name string, value object.Object) error {
	switch name {
	case "token":
		token, objErr := object.AsString(value)
		if objErr != nil {
			return objErr.Value()
		}
		return v.client.SetToken(token)
	case "namespace":
		namespace, objErr := object.AsString(value)
		if objErr != nil {
			return objErr.Value()
		}
		return v.client.SetNamespace(namespace)
	}
	return object.TypeErrorf("type error: %s object has no attribute %q", VAULT, name)
}

func (v *Vault) GetAttr(name string) (object.Object, bool) {
	switch name {
	case "set_token":
		return object.NewBuiltin("vault.set_token", func(ctx context.Context, args ...object.Object) object.Object {
			if err := arg.Require("vault.set_token", 1, args); err != nil {
				return err
			}
			token, objErr := object.AsString(args[0])
			if objErr != nil {
				return objErr
			}
			if err := v.client.SetToken(token); err != nil {
				return object.NewError(err)
			}
			return object.Nil
		}), true
	case "app_role_login":
		return object.NewBuiltin("vault.app_role_login", func(ctx context.Context, args ...object.Object) object.Object {
			if err := arg.Require("vault.app_role_login", 2, args); err != nil {
				return err
			}
			roleID, objErr := object.AsString(args[0])
			if objErr != nil {
				return objErr
			}
			secretID, objErr := object.AsString(args[1])
			if objErr != nil {
				return objErr
			}
			resp, err := v.client.Auth.AppRoleLogin(ctx, schema.AppRoleLoginRequest{
				RoleId:   roleID,
				SecretId: secretID,
			})
			if err != nil {
				return object.NewError(err)
			}
			if err := v.client.SetToken(resp.Auth.ClientToken); err != nil {
				return object.NewError(err)
			}
			return object.Nil
		}), true
	case "write":
		return object.NewBuiltin("vault.write", func(ctx context.Context, args ...object.Object) object.Object {
			if err := arg.Require("vault.write", 2, args); err != nil {
				return err
			}
			data, objErr := object.AsMap(args[0])
			if objErr != nil {
				return objErr
			}
			path, objErr := object.AsString(args[1])
			if objErr != nil {
				return objErr
			}
			body, err := data.MarshalJSON()
			if err != nil {
				return object.NewError(err)
			}
			resp, err := v.client.WriteFromBytes(ctx, path, body)
			if err != nil {
				return object.NewError(err)
			}
			return object.FromGoType(resp.Data)
		}), true
	case "write_raw":
		return object.NewBuiltin("vault.write", func(ctx context.Context, args ...object.Object) object.Object {
			if err := arg.Require("vault.write", 2, args); err != nil {
				return err
			}
			body, objErr := object.AsBytes(args[0])
			if objErr != nil {
				return objErr
			}
			path, objErr := object.AsString(args[1])
			if objErr != nil {
				return objErr
			}
			resp, err := v.client.WriteFromBytes(ctx, path, body)
			if err != nil {
				return object.NewError(err)
			}
			return object.FromGoType(resp.Data)
		}), true
	case "list":
		return object.NewBuiltin("vault.list", func(ctx context.Context, args ...object.Object) object.Object {
			if err := arg.Require("vault.list", 1, args); err != nil {
				return err
			}
			path, objErr := object.AsString(args[0])
			if objErr != nil {
				return objErr
			}
			resp, err := v.client.List(ctx, path)
			if err != nil {
				return object.NewError(err)
			}
			return object.FromGoType(resp.Data)
		}), true
	case "delete":
		return object.NewBuiltin("vault.delete", func(ctx context.Context, args ...object.Object) object.Object {
			if err := arg.Require("vault.delete", 1, args); err != nil {
				return err
			}
			path, objErr := object.AsString(args[0])
			if objErr != nil {
				return objErr
			}
			resp, err := v.client.Delete(ctx, path)
			if err != nil {
				return object.NewError(err)
			}
			return object.FromGoType(resp.Data)
		}), true
	case "read":
		return object.NewBuiltin("vault.read", func(ctx context.Context, args ...object.Object) object.Object {
			if err := arg.Require("vault.read", 1, args); err != nil {
				return err
			}
			path, objErr := object.AsString(args[0])
			if objErr != nil {
				return objErr
			}
			resp, err := v.client.Read(ctx, path)
			if err != nil {
				return object.NewError(err)
			}
			return object.FromGoType(resp.Data)
		}), true
	case "read_raw":
		return object.NewBuiltin("vault.read_raw", func(ctx context.Context, args ...object.Object) object.Object {
			if err := arg.Require("vault.read_raw", 1, args); err != nil {
				return err
			}
			path, objErr := object.AsString(args[0])
			if objErr != nil {
				return objErr
			}
			resp, err := v.client.ReadRaw(ctx, path)
			if err != nil {
				return object.NewError(err)
			}
			return object.FromGoType(resp)
		}), true
	}
	return nil, false
}

func New(addr string) (*Vault, error) {
	var client *vault.Client

	client, err := vault.New(
		vault.WithAddress(addr),
		vault.WithRequestTimeout(30*time.Second),
	)
	if err != nil {
		return nil, err
	}

	return &Vault{
		client: client,
	}, nil
}

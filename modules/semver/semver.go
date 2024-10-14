//go:build semver
// +build semver

package semver

import (
	"context"
	"strings"

	"github.com/blang/semver/v4"
	"github.com/risor-io/risor/arg"
	"github.com/risor-io/risor/object"
)

func Parse(ctx context.Context, args ...object.Object) object.Object {
	if err := arg.Require("semver.parse", 1, args); err != nil {
		return err
	}

	str, err := object.AsString(args[0])
	if err != nil {
		return err
	}

	v, perr := semver.Make(str)
	if perr != nil {
		return object.NewError(perr)
	}

	pre := []string{}
	for _, p := range v.Pre {
		pre = append(pre, p.String())
	}
	vmap := map[string]object.Object{
		"major": object.NewInt(int64(v.Major)),
		"minor": object.NewInt(int64(v.Minor)),
		"patch": object.NewInt(int64(v.Patch)),
		"pre":   object.NewString(strings.Join(pre, ".")),
		"build": object.NewString(strings.Join(v.Build, ".")),
	}

	return object.NewMap(vmap)
}

func Build(ctx context.Context, args ...object.Object) object.Object {
	if err := arg.Require("semver.build", 1, args); err != nil {
		return err
	}

	str, err := object.AsString(args[0])
	if err != nil {
		return err
	}

	v, perr := semver.Make(str)
	if perr != nil {
		return object.NewError(perr)
	}

	return object.NewString(strings.Join(v.Build, "."))
}

func Pre(ctx context.Context, args ...object.Object) object.Object {
	if err := arg.Require("semver.pre", 1, args); err != nil {
		return err
	}

	str, err := object.AsString(args[0])
	if err != nil {
		return err
	}

	v, perr := semver.Make(str)
	if perr != nil {
		return object.NewError(perr)
	}

	pre := []string{}
	for _, p := range v.Pre {
		pre = append(pre, p.String())
	}

	return object.NewString(strings.Join(pre, "."))
}

func Validate(ctx context.Context, args ...object.Object) object.Object {
	if err := arg.Require("semver.validate", 1, args); err != nil {
		return err
	}

	str, err := object.AsString(args[0])
	if err != nil {
		return err
	}

	v, perr := semver.Make(str)
	if perr != nil {
		return object.NewError(perr)
	}

	validated := v.Validate()
	if validated != nil {
		return object.NewError(validated)
	}

	return nil
}

func Major(ctx context.Context, args ...object.Object) object.Object {
	if err := arg.Require("semver.major", 1, args); err != nil {
		return err
	}

	str, err := object.AsString(args[0])
	if err != nil {
		return err
	}

	v, perr := semver.Make(str)
	if perr != nil {
		return object.NewError(perr)
	}

	return object.NewInt(int64(v.Major))
}

func Minor(ctx context.Context, args ...object.Object) object.Object {
	if err := arg.Require("semver.minor", 1, args); err != nil {
		return err
	}

	str, err := object.AsString(args[0])
	if err != nil {
		return err
	}

	v, perr := semver.Make(str)
	if perr != nil {
		return object.NewError(perr)
	}

	return object.NewInt(int64(v.Minor))
}

func Patch(ctx context.Context, args ...object.Object) object.Object {
	if err := arg.Require("semver.patch", 1, args); err != nil {
		return err
	}

	str, err := object.AsString(args[0])
	if err != nil {
		return err
	}

	v, perr := semver.Make(str)
	if perr != nil {
		return object.NewError(perr)
	}

	return object.NewInt(int64(v.Patch))
}

func Compare(ctx context.Context, args ...object.Object) object.Object {
	if err := arg.Require("semver.compare", 2, args); err != nil {
		return err
	}

	v1s, err := object.AsString(args[0])
	if err != nil {
		return err
	}
	v1, perr := semver.Make(v1s)
	if perr != nil {
		return object.NewError(perr)
	}

	v2s, err := object.AsString(args[1])
	if err != nil {
		return err
	}
	v2, perr := semver.Make(v2s)
	if perr != nil {
		return object.NewError(perr)
	}

	return object.NewInt(int64(v1.Compare(v2)))
}

func Equals(ctx context.Context, args ...object.Object) object.Object {
	if err := arg.Require("semver.equals", 2, args); err != nil {
		return err
	}

	v1s, err := object.AsString(args[0])
	if err != nil {
		return err
	}
	v1, perr := semver.Make(v1s)
	if perr != nil {
		return object.NewError(perr)
	}

	v2s, err := object.AsString(args[1])
	if err != nil {
		return err
	}
	v2, perr := semver.Make(v2s)
	if perr != nil {
		return object.NewError(perr)
	}

	return object.NewBool(v1.Equals(v2))
}

func Module() *object.Module {
	return object.NewBuiltinsModule("semver", map[string]object.Object{
		"build":    object.NewBuiltin("build", Build),
		"compare":  object.NewBuiltin("compare", Compare),
		"equals":   object.NewBuiltin("now", Equals),
		"parse":    object.NewBuiltin("parse", Parse),
		"major":    object.NewBuiltin("major", Major),
		"minor":    object.NewBuiltin("minor", Minor),
		"patch":    object.NewBuiltin("patch", Patch),
		"validate": object.NewBuiltin("validate", Validate),
	})
}

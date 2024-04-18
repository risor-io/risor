package semver

import (
	"context"

	"github.com/risor-io/risor/object"
	"golang.org/x/mod/semver"
)

func IsValid(ctx context.Context, args ...object.Object) object.Object {
	numArgs := len(args)
	if numArgs != 1 {
		return object.NewArgsError("semver.is_valid", 1, numArgs)
	}
	s, err := object.AsString(args[0])
	if err != nil {
		return err
	}
	return object.NewBool(semver.IsValid(s))
}

func Canonical(ctx context.Context, args ...object.Object) object.Object {
	numArgs := len(args)
	if numArgs != 1 {
		return object.NewArgsError("semver.canonical", 1, numArgs)
	}
	s, err := object.AsString(args[0])
	if err != nil {
		return err
	}
	return object.NewString(semver.Canonical(s))
}

func Major(ctx context.Context, args ...object.Object) object.Object {
	numArgs := len(args)
	if numArgs != 1 {
		return object.NewArgsError("semver.major", 1, numArgs)
	}
	s, err := object.AsString(args[0])
	if err != nil {
		return err
	}
	return object.NewString(semver.Major(s))
}

func MajorMinor(ctx context.Context, args ...object.Object) object.Object {
	numArgs := len(args)
	if numArgs != 1 {
		return object.NewArgsError("semver.major_minor", 1, numArgs)
	}
	s, err := object.AsString(args[0])
	if err != nil {
		return err
	}
	return object.NewString(semver.MajorMinor(s))
}

func Prerelease(ctx context.Context, args ...object.Object) object.Object {
	numArgs := len(args)
	if numArgs != 1 {
		return object.NewArgsError("semver.prerelease", 1, numArgs)
	}
	s, err := object.AsString(args[0])
	if err != nil {
		return err
	}
	return object.NewString(semver.Prerelease(s))
}

func Build(ctx context.Context, args ...object.Object) object.Object {
	numArgs := len(args)
	if numArgs != 1 {
		return object.NewArgsError("semver.build", 1, numArgs)
	}
	s, err := object.AsString(args[0])
	if err != nil {
		return err
	}
	return object.NewString(semver.Build(s))
}

func Compare(ctx context.Context, args ...object.Object) object.Object {
	numArgs := len(args)
	if numArgs != 2 {
		return object.NewArgsError("semver.compare", 2, numArgs)
	}
	s1, err := object.AsString(args[0])
	if err != nil {
		return err
	}
	s2, err := object.AsString(args[1])
	if err != nil {
		return err
	}
	return object.NewInt(int64(semver.Compare(s1, s2)))
}

func Module() *object.Module {
	return object.NewBuiltinsModule("semver", map[string]object.Object{
		"build":       object.NewBuiltin("build", Build),
		"canonical":   object.NewBuiltin("canonical", Canonical),
		"compare":     object.NewBuiltin("compare", Compare),
		"is_valid":    object.NewBuiltin("is_valid", IsValid),
		"major_minor": object.NewBuiltin("major_minor", MajorMinor),
		"major":       object.NewBuiltin("major", Major),
		"prerelease":  object.NewBuiltin("prerelease", Prerelease),
	})
}

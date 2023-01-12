package rand

import (
	"context"
	crand "crypto/rand"
	"encoding/binary"
	"fmt"
	"math/rand"

	"github.com/cloudcmds/tamarin/arg"
	"github.com/cloudcmds/tamarin/object"
	"github.com/cloudcmds/tamarin/scope"
)

// Name of this module
const Name = "rand"

func Seed() {
	var b [8]byte
	if _, err := crand.Read(b[:]); err != nil {
		panic("crypto/rand failed: " + err.Error())
	}
	rand.Seed(int64(binary.LittleEndian.Uint64(b[:])))
}

func Float(ctx context.Context, args ...object.Object) object.Object {
	if err := arg.Require("rand.float", 0, args); err != nil {
		return err
	}
	return object.NewFloat(rand.Float64())
}

func Int(ctx context.Context, args ...object.Object) object.Object {
	if err := arg.Require("rand.int", 0, args); err != nil {
		return err
	}
	return object.NewInt(rand.Int63())
}

func IntN(ctx context.Context, args ...object.Object) object.Object {
	if err := arg.Require("rand.intn", 1, args); err != nil {
		return err
	}
	n, err := object.AsInt(args[0])
	if err != nil {
		return err
	}
	return object.NewInt(rand.Int63n(n))
}

func NormFloat(ctx context.Context, args ...object.Object) object.Object {
	if err := arg.Require("rand.norm_float", 0, args); err != nil {
		return err
	}
	return object.NewFloat(rand.NormFloat64())
}

func ExpFloat(ctx context.Context, args ...object.Object) object.Object {
	if err := arg.Require("rand.exp_float", 0, args); err != nil {
		return err
	}
	return object.NewFloat(rand.ExpFloat64())
}

func Shuffle(ctx context.Context, args ...object.Object) object.Object {
	if err := arg.Require("rand.shuffle", 1, args); err != nil {
		return err
	}
	ls, err := object.AsList(args[0])
	if err != nil {
		return err
	}
	items := ls.Value()
	rand.Shuffle(len(items), func(i, j int) {
		items[i], items[j] = items[j], items[i]
	})
	return ls
}

func Module(parentScope *scope.Scope) (*object.Module, error) {
	s := scope.New(scope.Opts{
		Name:   fmt.Sprintf("module:%s", Name),
		Parent: parentScope,
	})

	m := object.NewModule(Name, s)

	if err := s.AddBuiltins([]*object.Builtin{
		object.NewBuiltin("float", Float, m),
		object.NewBuiltin("int", Int, m),
		object.NewBuiltin("intn", IntN, m),
		object.NewBuiltin("norm_float", NormFloat, m),
		object.NewBuiltin("exp_float", ExpFloat, m),
		object.NewBuiltin("shuffle", Shuffle, m),
	}); err != nil {
		return nil, err
	}
	return m, nil
}

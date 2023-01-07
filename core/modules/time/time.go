package time

import (
	"context"
	"fmt"
	"time"

	"github.com/cloudcmds/tamarin/core/arg"
	"github.com/cloudcmds/tamarin/core/object"
	"github.com/cloudcmds/tamarin/core/scope"
)

// Name of this module
const Name = "time"

func Now(ctx context.Context, args ...object.Object) object.Object {
	if err := arg.Require("time.now", 0, args); err != nil {
		return err
	}
	return object.NewTime(time.Now())
}

func Parse(ctx context.Context, args ...object.Object) object.Object {
	if err := arg.Require("time.parse", 2, args); err != nil {
		return err
	}
	layout, err := object.AsString(args[0])
	if err != nil {
		return err
	}
	value, err := object.AsString(args[1])
	if err != nil {
		return err
	}
	t, parseErr := time.Parse(layout, value)
	if parseErr != nil {
		return object.NewErrResult(object.NewError(parseErr))
	}
	return object.NewOkResult(object.NewTime(t))
}

func Sleep(ctx context.Context, args ...object.Object) object.Object {
	if err := arg.Require("time.sleep", 1, args); err != nil {
		return err
	}
	d, err := object.AsFloat(args[0])
	if err != nil {
		return err
	}
	timer := time.NewTimer(time.Duration(d*1000) * time.Millisecond)
	defer timer.Stop()
	select {
	case <-ctx.Done():
	case <-timer.C:
	}
	return object.Nil
}

func Module(parentScope *scope.Scope) (*object.Module, error) {
	s := scope.New(scope.Opts{
		Name:   fmt.Sprintf("module:%s", Name),
		Parent: parentScope,
	})

	m := object.NewModule(Name, s)

	if err := s.AddBuiltins([]*object.Builtin{
		object.NewBuiltin("now", Now, m),
		object.NewBuiltin("parse", Parse, m),
		object.NewBuiltin("sleep", Sleep, m),
	}); err != nil {
		return nil, err
	}

	s.Declare("ANSIC", object.NewString(time.ANSIC), true)
	s.Declare("UnixDate", object.NewString(time.UnixDate), true)
	s.Declare("RubyDate", object.NewString(time.RubyDate), true)
	s.Declare("RFC822", object.NewString(time.RFC822), true)
	s.Declare("RFC822Z", object.NewString(time.RFC822Z), true)
	s.Declare("RFC850", object.NewString(time.RFC850), true)
	s.Declare("RFC1123", object.NewString(time.RFC1123), true)
	s.Declare("RFC1123Z", object.NewString(time.RFC1123Z), true)
	s.Declare("RFC3339", object.NewString(time.RFC3339), true)
	s.Declare("RFC3339Nano", object.NewString(time.RFC3339Nano), true)
	s.Declare("Kitchen", object.NewString(time.Kitchen), true)
	s.Declare("Stamp", object.NewString(time.Stamp), true)
	s.Declare("StampMilli", object.NewString(time.StampMilli), true)
	s.Declare("StampMicro", object.NewString(time.StampMicro), true)
	s.Declare("StampNano", object.NewString(time.StampNano), true)
	return m, nil
}

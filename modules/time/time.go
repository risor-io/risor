package time

import (
	"context"
	"time"

	"github.com/cloudcmds/tamarin/arg"
	"github.com/cloudcmds/tamarin/object"
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

func Module() *object.Module {
	m := object.NewModule(Name)
	m.Register("now", object.NewBuiltin("now", Now, m))
	m.Register("parse", object.NewBuiltin("parse", Parse, m))
	m.Register("sleep", object.NewBuiltin("sleep", Sleep, m))
	m.Register("ANSIC", object.NewString(time.ANSIC))
	m.Register("UnixDate", object.NewString(time.UnixDate))
	m.Register("RubyDate", object.NewString(time.RubyDate))
	m.Register("RFC822", object.NewString(time.RFC822))
	m.Register("RFC822Z", object.NewString(time.RFC822Z))
	m.Register("RFC850", object.NewString(time.RFC850))
	m.Register("RFC1123", object.NewString(time.RFC1123))
	m.Register("RFC1123Z", object.NewString(time.RFC1123Z))
	m.Register("RFC3339", object.NewString(time.RFC3339))
	m.Register("RFC3339Nano", object.NewString(time.RFC3339Nano))
	m.Register("Kitchen", object.NewString(time.Kitchen))
	m.Register("Stamp", object.NewString(time.Stamp))
	m.Register("StampMilli", object.NewString(time.StampMilli))
	m.Register("StampMicro", object.NewString(time.StampMicro))
	m.Register("StampNano", object.NewString(time.StampNano))
	return m
}

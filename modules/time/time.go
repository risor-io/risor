package time

import (
	"context"
	"time"

	"github.com/cloudcmds/tamarin/v2/internal/arg"
	"github.com/cloudcmds/tamarin/v2/object"
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
		return object.NewError(parseErr)
	}
	return object.NewTime(t)
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
	m := object.NewBuiltinsModule(Name, map[string]object.Object{
		"now":         object.NewBuiltin("now", Now),
		"parse":       object.NewBuiltin("parse", Parse),
		"sleep":       object.NewBuiltin("sleep", Sleep),
		"ANSIC":       object.NewString(time.ANSIC),
		"UnixDate":    object.NewString(time.UnixDate),
		"RubyDate":    object.NewString(time.RubyDate),
		"RFC822":      object.NewString(time.RFC822),
		"RFC822Z":     object.NewString(time.RFC822Z),
		"RFC850":      object.NewString(time.RFC850),
		"RFC1123":     object.NewString(time.RFC1123),
		"RFC1123Z":    object.NewString(time.RFC1123Z),
		"RFC3339":     object.NewString(time.RFC3339),
		"RFC3339Nano": object.NewString(time.RFC3339Nano),
		"Kitchen":     object.NewString(time.Kitchen),
		"Stamp":       object.NewString(time.Stamp),
		"StampMilli":  object.NewString(time.StampMilli),
		"StampMicro":  object.NewString(time.StampMicro),
		"StampNano":   object.NewString(time.StampNano),
	})
	return m
}

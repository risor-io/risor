//go:build sched
// +build sched

package sched

import (
	"context"
	"time"

	"codnect.io/chrono"
	"github.com/risor-io/risor/object"
)

// Module returns the sched module.
func Module() *object.Module {
	return object.NewBuiltinsModule(
		"sched", map[string]object.Object{
			"cron":  object.NewBuiltin("sched.cron", Cron),
			"every": object.NewBuiltin("sched.every", Every),
			"once":  object.NewBuiltin("sched.once", Once),
		},
	)
}

// Cron schedules a function to run at a specific time using a cron like expression.
//
// The first argument is the cron expression.
// The second argument is the function to be scheduled.
func Cron(ctx context.Context, args ...object.Object) object.Object {
	if len(args) != 2 {
		return object.Errorf("missing arguments, 2 required")
	}

	cronLine, err := object.AsString(args[0])
	if err != nil {
		return err
	}

	fn, ok := args[1].(*object.Function)
	if !ok {
		return object.Errorf("expected function, got %s", args[1].Type())
	}

	// GetCloneCallFunc returns a function safe to be called from a different goroutine.
	// NOTE: unsure about the above, but seems to be working
	cfunc, ok := object.GetCloneCallFunc(ctx)
	if !ok {
		return object.EvalErrorf("eval error: context did not contain a call function")
	}

	taskScheduler := chrono.NewDefaultTaskScheduler()
	t, nerr := taskScheduler.ScheduleWithCron(func(context.Context) {
		_, _ = cfunc(ctx, fn, nil)
	}, cronLine)
	if nerr != nil {
		return object.NewError(nerr)
	}

	return &task{t: t}
}

// Every schedules a function to run every n seconds.
//
// The first argument is the interval in seconds (float).
// The second argument is the function to be scheduled.
func Every(ctx context.Context, args ...object.Object) object.Object {
	if len(args) != 2 {
		return object.Errorf("missing arguments, 2 required")
	}

	dur, err := object.AsString(args[0])
	if err != nil {
		return err
	}

	fn, ok := args[1].(*object.Function)
	if !ok {
		return object.Errorf("expected function, got %s", args[1].Type())
	}

	// GetCloneCallFunc returns a function safe to be called from a different goroutine.
	// NOTE: unsure about the above, but seems to be working
	cfunc, ok := object.GetCloneCallFunc(ctx)
	if !ok {
		return object.EvalErrorf("eval error: context did not contain a call function")
	}

	duration, nerr := time.ParseDuration(dur)
	if nerr != nil {
		return object.NewError(nerr)
	}

	taskScheduler := chrono.NewDefaultTaskScheduler()
	t, nerr := taskScheduler.ScheduleAtFixedRate(func(context.Context) {
		_, _ = cfunc(ctx, fn, nil)
	}, duration)
	if nerr != nil {
		return object.NewError(err)
	}

	return &task{t: t}
}

// Once schedules a function to run once.
//
// The first argument is the interval in seconds (float).
// The second argument is the function to be scheduled.
func Once(ctx context.Context, args ...object.Object) object.Object {
	if len(args) != 2 {
		return object.Errorf("missing arguments, 2 required")
	}

	dur, err := object.AsString(args[0])
	if err != nil {
		return err
	}

	duration, nerr := time.ParseDuration(dur)
	if nerr != nil {
		return object.NewError(nerr)
	}

	start := time.Now().Add(duration)

	fn, ok := args[1].(*object.Function)
	if !ok {
		return object.Errorf("expected function, got %s", args[1].Type())
	}

	// GetCloneCallFunc returns a function safe to be called from a different goroutine.
	// NOTE: unsure about the above, but seems to be working
	cfunc, ok := object.GetCloneCallFunc(ctx)
	if !ok {
		return object.EvalErrorf("eval error: context did not contain a call function")
	}

	taskScheduler := chrono.NewDefaultTaskScheduler()
	t, nerr := taskScheduler.Schedule(func(ctx context.Context) {
		_, _ = cfunc(ctx, fn, nil)
	}, chrono.WithTime(start))
	if nerr != nil {
		return object.NewError(err)
	}

	return &task{t: t}
}

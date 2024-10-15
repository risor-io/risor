//go:build sched
// +build sched

package sched

import (
	"context"
	"testing"
	"time"

	"github.com/risor-io/risor/object"
	"github.com/stretchr/testify/assert"
)

func TestCron(t *testing.T) {
	var executed int
	callFn := func(ctx context.Context, fn *object.Function, args []object.Object) (object.Object, error) {
		executed++
		return nil, nil
	}
	ctx := object.WithCloneCallFunc(context.Background(), callFn)
	var fn *object.Function

	cronExpr := object.NewString("invalid-cronline")
	task := Cron(ctx, cronExpr, fn)
	assert.True(t, object.IsError(task))

	cronExpr = object.NewString("*/1 * * * * *") // Every second
	task = Cron(ctx, cronExpr, fn)
	assert.False(t, object.IsError(task))
	assert.NotNil(t, task)
	assert.Equal(t, "sched.task", string(task.Type()))
	assert.Equal(t, 0, executed)

	_, ok := task.GetAttr("is_running")
	assert.True(t, ok)

	// Wait a second to allow the cron job to execute
	time.Sleep(1 * time.Second)
	assert.Equal(t, 1, executed)

	_, ok = task.GetAttr("cancel")
	assert.True(t, ok)
}

func TestEvery(t *testing.T) {
	var executed int
	callFn := func(ctx context.Context, fn *object.Function, args []object.Object) (object.Object, error) {
		if executed == 2 {
			return nil, nil
		}
		executed++
		return nil, nil
	}
	ctx := object.WithCloneCallFunc(context.Background(), callFn)
	var fn *object.Function

	// Schedule the function to run every 1 millisecond
	interval := object.NewString("1ms")
	task := Every(ctx, interval, fn)
	assert.False(t, object.IsError(task))
	assert.Equal(t, "sched.task", string(task.Type()))

	// Wait for a few milliseconds to allow the job to execute
	time.Sleep(50 * time.Millisecond)
	assert.Equal(t, 2, executed)

	interval = object.NewString("1foo")
	task = Every(ctx, interval, fn)
	assert.True(t, object.IsError(task))
}

func TestOnce(t *testing.T) {
	var executed int
	callFn := func(ctx context.Context, fn *object.Function, args []object.Object) (object.Object, error) {
		executed++
		return nil, nil
	}
	ctx := object.WithCloneCallFunc(context.Background(), callFn)
	var fn *object.Function

	// Schedule the function to run every 1ms
	interval := object.NewString("1ms")
	task := Once(ctx, interval, fn)
	assert.False(t, object.IsError(task))
	assert.Equal(t, "sched.task", string(task.Type()))

	// Wait some time for the job to complete
	time.Sleep(50 * time.Millisecond)
	assert.Equal(t, 1, executed)
}

func TestEqual(t *testing.T) {
	callFn := func(ctx context.Context, fn *object.Function, args []object.Object) (object.Object, error) {
		return nil, nil
	}
	ctx := object.WithCloneCallFunc(context.Background(), callFn)
	var fn *object.Function

	cronExpr := object.NewString("* * * * * *")
	task := Cron(ctx, cronExpr, fn)
	task2 := Cron(ctx, cronExpr, fn)

	assert.True(t, task.Equals(task).(*object.Bool).Value())
	assert.False(t, task.Equals(task2).(*object.Bool).Value())
}

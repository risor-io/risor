import { Callout } from 'nextra/components';

# Sched

<Callout type="info" emoji="ℹ️">
This module is not included with Risor.
</Callout>

The `sched` module exposes a simple interface to schedule tasks, powered by [chrono](https://github.com/codnect/chrono).

```go
import sched
import time

once := sched.once("1s", func(){
  print("once")
})

cron := sched.cron("*/1 * * * * *", func(){
  print("hola")
})
print(cron.is_running())

every := sched.every("1m", func() {
  print("every 1 minute")
})

for {
  time.sleep(1)
}
```

## Functions

Functions are non-blocking and return a `task` object.

### cron

```go filename="Function signature"
cron(cronline string, fn func)
```

Creates a new cron job.

```go copy filename="Example"
// Run every second
s := sched.cron("*/1 * * * * *", func() {
	print("hello world!")
})
```

### every

```go filename="Function signature"
every(duration string, fn func)
```

Creates a new task that runs regularly, at the specified duration.
The string format is documented in [the standard Go library](https://pkg.go.dev/time#ParseDuration).

```go copy filename="Example"
// Run every minute
s := sched.every("1m"", func() {
	print("hello world!")
})
```

```go copy filename="Example"
// Run every 500 milliseconds
s := sched.every(0.5, func() {
	print("hello world!")
})
```

## Types

### task

A Task object.

#### Attributes

| Name           | Type                           | Description                                  |
| -------------- | ------------------------------ | -------------------------------------------- |
| cancel         | func()                         | Cancels the task                             |
| running        | func()                         | True if the task is running                  |

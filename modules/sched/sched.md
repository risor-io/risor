import { Callout } from 'nextra/components';

# sched

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

The `cronline` string is a space-separated list of 6 fields, representing the time to run the task.

```
SECOND MINUTE HOUR DAY MONTH DAYOFWEEK
   *      *     *   *    *      *
```

Some examples:

- `* * * * * *` every second
- `*/5 * * * * *` every 5 seconds
- `0 * * * * *` every minute
- `0 0 * * * *` every hour
- `20 45 18 5-20/3 * *` every day at 18:45:20, from the 5th to the 20th of the month, every 3 days
- `0 0 0 1 SEP *` every year on September 1st at midnight
- `0 0 0 1 5 SUN` every year on the first Sunday of May at midnight


```go copy filename="Example"
// Run every second
task := sched.cron("*/1 * * * * *", func() {
	print("hello world!")
})
```

Functions run in a separate goroutine, and errors are ignored, so the main program can continue to run. Make sure to handle errors in your function.

### every

```go filename="Function signature"
every(duration string, fn func)
```

Creates a new task that runs regularly, at the specified duration.

The string format is documented in [the standard Go library](https://pkg.go.dev/time#ParseDuration).

```go copy filename="Example"
// Run every minute
task := sched.every("1m"", func() {
	print("hello world!")
})
```

### once

```go filename="Function signature"
once(duration string, fn func)
```

Creates a new task that runs once, after the specified duration.

```go copy filename="Example"
// Run a task in one 1h, once
task := sched.once("1h", func() {
	print("hello world!")
})
```

`cron` can be used to run a task once at a specific time, if more control is needed.

## Types

### task

A Task object returned by `cron`, `every` and `once` functions.

#### Attributes

| Name           | Type                           | Description                                  |
| -------------- | ------------------------------ | -------------------------------------------- |
| cancel         | func()                         | Cancels the task                             |
| is_running     | func()                         | True if the task is running                  |

# carbon

The `carbon` module provides a convenient set of utilities for working with
dates and times. These functions are meant to aid writing more readable and
maintainable date/time code.

The core functionality is provided by
[github.com/golang-module/carbon](https://github.com/golang-module/carbon).

## Module

```go copy filename="Function signature"
carbon(input ...object)
```

The `carbon` module object itself is callable in order to provide a shorthand for
initializing a new carbon object.

Three input variations are supported:

- No input: Returns the current time as a carbon object.
- A string input: Parses the string into a carbon object.
- A time input: Converts the time into a carbon object.

```go copy filename="Example"
>>> carbon()
carbon.carbon(2024-02-03 12:03:55)
>>> carbon("2015-03-01 12:01:00")
carbon.carbon(2015-03-01 12:01:00)
>>> carbon(time.now())
carbon.carbon(2024-02-03 12:03:55)
```

## Functions

### now

```go filename="Function signature"
now(timezone ...string) carbon
```

Returns the current time as a carbon object. A timezone string may optionally
be provided.

```go filename="Example"
>>> carbon.now()
carbon.carbon(2024-02-03 12:03:55)
>>> carbon.now("America/Sao_Paulo").timezone()
"-03"
```

### parse

```go filename="Function signature"
parse(value, timezone ...string) carbon
```

Parses the time string into a carbon object. A timezone may optionally be provided.
If the time string is invalid, an error is raised.

```go filename="Example"
>>> carbon.parse("2021-03-20")
carbon.carbon(2021-03-20 00:00:00)
>>> carbon.parse("2024-02-03 12:03:55")
carbon.carbon(2024-02-03 12:03:55)
>>> carbon.parse("2024-02-03 12:03:55", "America/Sao_Paulo")
carbon.carbon(2024-02-03 12:03:55)
>>> carbon.parse("oops")
// value error: invalid time string (got "oops")
```

### yesterday

```go filename="Function signature"
yesterday(timezone ...string) carbon
```

Returns a time object for yesterday, at the current time of day. A timezone
may optionally be provided.

```go filename="Example"
>>> carbon.yesterday()
carbon.carbon(2024-02-02 12:03:55)
```

### tomorrow

```go filename="Function signature"
tomorrow(timezone ...string) carbon
```

Returns a time object for tomorrow, at the current time of day. A timezone
may optionally be provided.

```go filename="Example"
>>> carbon.tomorrow()
carbon.carbon(2024-02-04 12:03:55)
```

## Types

### carbon

The carbon object represents a point in time and offers a variety of helper
methods for performing date and time operations. Many of the operations reflect
ways humans think about dates and times, such as "add a day" or "get the age".

#### Methods

##### add_day

```go filename="Method signature"
add_day() carbon
```

Returns a new carbon object with the day incremented by one.

```go filename="Example"
>>> c := carbon.now()
>>> c.add_day()
carbon.carbon(2024-02-04 12:03:55)
```

##### add_days

```go filename="Method signature"
add_days(days int) carbon
```

Returns a new carbon object with the day incremented by the specified number of days.

```go filename="Example"
>>> c := carbon.now()
>>> c.add_days(7)
carbon.carbon(2024-02-10 12:03:55)
```

##### sub_day

```go filename="Method signature"
sub_day() carbon
```

Returns a new carbon object with the day decremented by one.

```go filename="Example"
>>> c := carbon.now()
>>> c.sub_day()
carbon.carbon(2024-02-02 12:03:55)
```

##### sub_days

```go filename="Method signature"
sub_days(days int) carbon
```

Returns a new carbon object with the day decremented by the specified number of days.

```go filename="Example"
>>> c := carbon.now()
>>> c.sub_days(7)
carbon.carbon(2024-01-27 12:03:55)
```

##### timezone

```go filename="Method signature"
timezone(tz string) carbon
```

Returns a new carbon object with the timezone set to the specified timezone string.

```go filename="Example"
>>> c := carbon.now()
>>> c.timezone("America/New_York")
carbon.carbon(2024-02-03 07:03:55)
```

##### age

```go filename="Method signature"
age(other carbon) int
```

Returns the age (in years) between the current carbon object and the specified carbon object.

```go filename="Example"
>>> c1 := carbon.parse("2000-01-01")
>>> c2 := carbon.now()
>>> c2.age(c1)
24
```

##### days_in_month

```go filename="Method signature"
days_in_month() int
```

Returns the number of days in the month of the current carbon object.

```go filename="Example"
>>> c := carbon.parse("2024-02-15")
>>> c.days_in_month()
29
```

##### days_in_year

```go filename="Method signature"
days_in_year() int
```

Returns the number of days in the year of the current carbon object.

```go filename="Example"
>>> c := carbon.parse("2024-02-15")
>>> c.days_in_year()
366
```

##### week_of_month

```go filename="Method signature"
week_of_month() int
```

Returns the week of the month for the current carbon object.

```go filename="Example"
>>> c := carbon.parse("2024-02-15")
>>> c.week_of_month()
3
```

##### week_of_year

```go filename="Method signature"
week_of_year() int
```

Returns the week of the year for the current carbon object.

```go filename="Example"
>>> c := carbon.parse("2024-02-15")
>>> c.week_of_year()
7
```

##### timestamp

```go filename="Method signature"
timestamp() int64
```

Returns the Unix timestamp (in seconds) for the current carbon object.

```go filename="Example"
>>> c := carbon.now()
>>> c.timestamp()
1677886635
```

##### string

```go filename="Method signature"
string() string
```

Returns the string representation of the current carbon object in the format
"YYYY-MM-DD HH:MM:SS".

```go filename="Example"
>>> c := carbon.now()
>>> c.string()
"2024-02-03 12:03:55"
```

##### is_valid

```go filename="Method signature"
is_valid() bool
```

Returns true if the current carbon object represents a valid date and time,
false otherwise.

```go filename="Example"
>>> c1 := carbon.parse("2024-02-29")
>>> c1.is_valid()
false
>>> c2 := carbon.parse("2024-02-28")
>>> c2.is_valid()
true
```

##### is_am

```go filename="Method signature"
is_am() bool
```

Returns true if the current carbon object represents a time in the AM
(before noon), false otherwise.

```go filename="Example"
>>> c1 := carbon.parse("2024-02-03 08:00:00")
>>> c1.is_am()
true
>>> c2 := carbon.parse("2024-02-03 15:00:00")
>>> c2.is_am()
false
```

##### is_pm

```go filename="Method signature"
is_pm() bool
```

Returns true if the current carbon object represents a time in the PM (after noon), false otherwise.

```go filename="Example"
>>> c1 := carbon.parse("2024-02-03 08:00:00")
>>> c1.is_pm()
false
>>> c2 := carbon.parse("2024-02-03 15:00:00")
>>> c2.is_pm()
true
```

##### is_leap_year

```go filename="Method signature"
is_leap_year() bool
```

Returns true if the year of the current carbon object is a leap year, false otherwise.

```go filename="Example"
>>> c1 := carbon.parse("2024-02-29")
>>> c1.is_leap_year()
true
>>> c2 := carbon.parse("2023-02-28")
>>> c2.is_leap_year()
false
```

##### is_future

```go filename="Method signature"
is_future() bool
```

Returns true if the current carbon object represents a time in the future
(after the current time), false otherwise.

```go filename="Example"
>>> c1 := carbon.parse("2024-02-03 15:00:00")
>>> c1.is_future()
true
>>> c2 := carbon.now()
>>> c2.is_future()
false
```

##### is_past

```go filename="Method signature"
is_past() bool
```

Returns true if the current carbon object represents a time in the past
(before the current time), false otherwise.

```go filename="Example"
>>> c1 := carbon.parse("2024-02-03 15:00:00")
>>> c1.is_past()
false
>>> c2 := carbon.now().sub_days(1)
>>> c2.is_past()
true
```

##### is_today

```go filename="Method signature"
is_today() bool
```

Returns true if the current carbon object represents a time on the current day,
false otherwise.

```go filename="Example"
>>> c1 := carbon.now()
>>> c1.is_today()
true
>>> c2 := carbon.parse("2024-02-02 15:00:00")
>>> c2.is_today()
false
```

##### is_yesterday

```go filename="Method signature"
is_yesterday() bool
```

Returns true if the current carbon object represents a time on the previous day,
false otherwise.

```go filename="Example"
>>> c1 := carbon.now()
>>> c1.is_yesterday()
false
>>> c2 := carbon.yesterday()
>>> c2.is_yesterday()
true
```

##### day

```go filename="Method signature"
day() int
```

Returns the day of the month for the current carbon object.

```go filename="Example"
>>> c := carbon.parse("2024-02-15 12:00:00")
>>> c.day()
15
```

##### month

```go filename="Method signature"
month() int
```

Returns the month for the current carbon object.

```go filename="Example"
>>> c := carbon.parse("2024-02-15 12:00:00")
>>> c.month()
2
```

##### year

```go filename="Method signature"
year() int
```

Returns the year for the current carbon object.

```go filename="Example"
>>> c := carbon.parse("2024-02-15 12:00:00")
>>> c.year()
2024
```

##### hour

```go filename="Method signature"
hour() int
```

Returns the hour for the current carbon object.

```go filename="Example"
>>> c := carbon.parse("2024-02-15 12:00:00")
>>> c.hour()
12
```

##### minute

```go filename="Method signature"
minute() int
```

Returns the minute for the current carbon object.
```go filename="Example"
>>> c := carbon.parse("2024-02-15 12:30:00")
>>> c.minute()
30
```

##### second

```go filename="Method signature"
second() int
```

Returns the second for the current carbon object.

```go filename="Example"
>>> c := carbon.parse("2024-02-15 12:30:45")
>>> c.second()
45
```

##### diff_for_humans

```go filename="Method signature"
diff_for_humans(other carbon) string
```

Returns a human-readable string representing the difference between the current
carbon object and the specified carbon object.

```go filename="Example"
>>> c1 := carbon.now()
>>> c2 := c1.add_days(7)
>>> c2.diff_for_humans(c1)
"1 week"
```

##### std_time

```go filename="Method signature"
std_time() time.Time
```

Returns a time.Time object representing the current carbon object.

```go filename="Example"
>>> c := carbon.now()
>>> t := c.std_time()
>>> fmt.Println(t)
2024-02-03 12:03:55 +0000 UTC
```

##### to_date_string

```go filename="Method signature"
to_date_string() string
```

Returns a string representation of the current carbon object in the format
"YYYY-MM-DD".

```go filename="Example"
>>> c := carbon.now()
>>> c.to_date_string()
"2024-02-03"
```

##### to_time_string

```go filename="Method signature"
to_time_string() string
```

Returns a string representation of the current carbon object in the format
"HH:MM:SS".

```go filename="Example"
>>> c := carbon.parse("2024-02-03 12:30:45")
>>> c.to_time_string()
"12:30:45"
```

##### to_datetime_string

```go filename="Method signature"
to_datetime_string() string
```

Returns a string representation of the current carbon object in the format
"YYYY-MM-DD HH:MM:SS".

```go filename="Example"
>>> c := carbon.parse("2024-02-03 12:30:45")
>>> c.to_datetime_string()
"2024-02-03 12:30:45"
```

# math

Module `math` provides constants and mathematical functions. It is primarily
a wrapper of the Go [math](https://pkg.go.dev/math) package, however it also
includes a `sum` function for summing a list of numbers.

Risor provides equivalence between float and int types, so many of the
functions in this module accept both float and int as inputs. In the documentation below, "number" is used to refer to either float or int.

## Constants

### PI

```go filename="Constant"
PI float
```

```go copy filename="Example"
>>> math.PI
3.141592653589793
```

### E

```go filename="Constant"
E float
```

```go copy filename="Example"
>>> math.E
2.718281828459045
```

## Functions

### abs

```go filename="Function signature"
abs(x number) number
```

Returns the absolute value of x.

```go copy filename="Example"
>>> math.abs(-2)
2
>>> math.abs(3.3)
3.3
```

### atan2

```go filename="Function signature"
atan2(y, x number) number
```

Returns the arc tangent value of y/x.

```go copy filename="Example"
>>> math.atan2(1, 2)
0.4636476090008061
```

### sqrt

```go filename="Function signature"
sqrt(x number) float
```

Returns the square root of x.

```go copy filename="Example"
>>> math.sqrt(4)
2
>>> math.sqrt(2)
1.4142135623730951
```

### min

```go filename="Function signature"
min(x, y number) float
```

Returns the smaller of x or y.

```go copy filename="Example"
>>> math.min(1, 2)
1.0
>>> math.min(3, -2.5)
-2.5
```

### max

```go filename="Function signature"
max(x, y number) float
```

Returns the larger of x or y.

```go copy filename="Example"
>>> math.max(1, 2)
2.0
>>> math.max(-3, 2.5)
2.5
```

### floor

```go filename="Function signature"
floor(x number) number
```

Returns the largest integer value less than or equal to x.

```go copy filename="Example"
>>> math.floor(2.5)
2
>>> math.floor(-2.5)
-3
>>> math.floor(3)
3
```

### ceil

```go filename="Function signature"
ceil(x number) number
```

Returns the smallest integer value greater than or equal to x.

```go copy filename="Example"
>>> math.ceil(2.5)
3
>>> math.ceil(-2.5)
-2
```

### sin

```go filename="Function signature"
sin(x number) float
```

Returns the sine of x.

```go copy filename="Example"
>>> math.sin(0)
0
>>> math.sin(math.PI / 2)
1
```

### cos

```go filename="Function signature"
cos(x number) float
```

Returns the cosine of x.

```go copy filename="Example"
>>> math.cos(0)
1
>>> math.cos(math.PI / 2)
0
```

### tan

```go filename="Function signature"
tan(x number) float
```

Returns the tangent of x.

```go copy filename="Example"
>>> math.tan(0)
0
>>> math.tan(math.PI / 4)
0.9999999999999998
```

### mod

```go filename="Function signature"
mod(x, y number) float
```

Returns the remainder of x divided by y.

```go copy filename="Example"
>>> math.mod(5, 2)
1
>>> math.mod(5, 3)
2
```

### log

```go filename="Function signature"
log(x number) float
```

Returns the natural logarithm of x.

```go copy filename="Example"
>>> math.log(1)
0
>>> math.log(math.E)
1
```

### log10

```go filename="Function signature"
log10(x number) float
```

Returns the base 10 logarithm of x.

```go copy filename="Example"
>>> math.log10(1)
0
>>> math.log10(10)
1
```

### log2

```go filename="Function signature"
log2(x number) float
```

Returns the base 2 logarithm of x.

```go copy filename="Example"
>>> math.log2(1)
0
>>> math.log2(8)
3
```

### pow

```go filename="Function signature"
pow(x, y number) float
```

Returns x raised to the power of y.

```go copy filename="Example"
>>> math.pow(2, 3)
8
>>> math.pow(2, 0.5)
1.4142135623730951
```

### pow10

```go filename="Function signature"
pow10(x number) float
```

Returns 10 raised to the power of x.

```go copy filename="Example"
>>> math.pow10(0)
1
>>> math.pow10(1)
10
>>> math.pow10(2)
100
```

### inf

```go filename="Function signature"
inf(x number) float 
```

Inf returns positive infinity if sign >= 0, negative infinity if sign < 0. 

```go copy filename="Example"
>>> math.inf()
+Inf
```

### is_inf

```go filename="Function signature"
is_inf(x number) bool
```

Returns true if x is positive or negative infinity.

```go copy filename="Example"
>>> math.is_inf(math.inf)
true
>>> math.is_inf(-math.inf)
true
>>> math.is_inf(0)
false
```

### round

```go filename="Function signature"
round(x number) float
```

Returns x rounded to the nearest integer.

```go copy filename="Example"
>>> math.round(1.4)
1
>>> math.round(1.5)
2
```

### sum

```go filename="Function signature"
sum(list) float
```

Returns the sum of all numbers in a list.

```go copy filename="Example"
>>> math.sum([1, 2, 3])
6
>>> math.sum([])
0
```

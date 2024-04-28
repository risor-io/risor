# bcrypt

The `bcrypt` module provides functions to hash and verify passwords using the
bcrypt algorithm.

The core functionality is provided by
[golang.org/x/crypto/bcrypt](https://pkg.go.dev/golang.org/x/crypto/bcrypt).

## Functions

### hash

```go filename="Function signature"
hash(password string, cost int = bcrypt.default_cost) byte_slice
```

Hash the password using bcrypt with the given cost. The cost is the number of
rounds to use. The cost must be between bcrypt.min_cost and bcrypt.max_cost. If
not provided, the cost defaults to bcrypt.default_cost (10).

```go copy filename="Example"
>>> bcrypt.hash("password")
byte_slice("$2a$10$vGwQDlEmqYug7JP.w/acKOProf3HsIYO3wI9CUxuxOc/RpqwWD0/C")
```

Note that the bcrypt hash is non-deterministic due to the random salt used in
the hashing process.

### compare

```go filename="Function signature"
compare(hash byte_slice, password string) bool
```

Compare the password with the bcrypt hash. Raises an error if the password does
not match the hash.

```go copy filename="Example"
>>> bcrypt.compare("$2a$10$vGwQDlEmqYug7JP.w/acKOProf3HsIYO3wI9CUxuxOc/RpqwWD0/C", "password")
>>> bcrypt.compare("$2a$10$vGwQDlEmqYug7JP.w/acKOProf3HsIYO3wI9CUxuxOc/RpqwWD0/C", "oops")
crypto/bcrypt: hashedPassword is not the hash of the given password
```

## Constants

### min_cost

```go filename="Constant"
min_cost int
```

The minimum cost that can be used for hashing.

```go copy filename="Example"
>>> bcrypt.min_cost
4
```

### max_cost

```go filename="Constant"
max_cost int
```

The maximum cost that can be used for hashing.

```go copy filename="Example"
>>> bcrypt.max_cost
31
```

### default_cost

```go copy filename="Example"
default_cost int
```

The default cost that is used for hashing if not provided.

```go copy filename="Example"
>>> bcrypt.default_cost
10
```

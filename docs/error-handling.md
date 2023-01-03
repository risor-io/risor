# Error Handling

## Operations that may fail

`Result` objects wrap `Ok` and `Err` values for operations that may fail.

```go
obj := json.unmarshal("true")
obj.unwrap() // returns true

failed := json.unmarshal("/not-valid/")
failed.is_err() // returns true
failed.unwrap() // raises error that stops execution
```

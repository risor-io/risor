# Error Traceback Design

## Summary

This document describes the design and implementation of error traceback functionality in Risor, which allows users to inspect the call stack when an error is raised.

## Motivation

Prior to this feature, when an error occurred in a Risor script, users only received the error message without any context about where in the code the error originated or the sequence of function calls that led to it. This made debugging complex scripts difficult, especially when errors occurred deep within nested function calls.

### Goals

1. Capture the call stack at the point where an error is raised
2. Provide access to traceback information through the error object
3. Preserve traceback information when errors are caught and re-raised
4. Include source file information when available
5. Maintain backward compatibility with existing error handling

## Design

### Core Components

#### 1. StackFrame Structure

```go
type StackFrame struct {
    FunctionName string
    FileName     string
    LineNumber   int
}
```

Each frame represents one level in the call stack, containing:
- The name of the function (or `<anonymous>` for unnamed functions, `<module>` for top-level code)
- The source file name (if available)
- The line number (currently always 0, reserved for future enhancement)

#### 2. Error Object Enhancement

The `Error` object was extended with:
- A `traceback []StackFrame` field to store the call stack
- A `WithTraceback()` method to attach a traceback
- A `GetTraceback()` method to retrieve the raw traceback
- A `Traceback()` method that returns a formatted string representation
- An `ErrorfWithTraceback()` constructor for creating errors with tracebacks

#### 3. VM Integration

The VM provides traceback capture through:
- A `captureStackTrace()` method that walks the VM's frame stack
- Registration of this method as a `TraceFunc` in the context
- Preservation of error objects (rather than unwrapping to Go errors) during error propagation

#### 4. Builtin Function Updates

The `error()` builtin was modified to:
- Capture the current stack trace when creating new errors
- Preserve existing tracebacks when re-raising error objects
- Attach the traceback to the created error

The `try()` builtin was updated to:
- Preserve error objects and their tracebacks when catching errors
- Pass the full error object (not just the message) to catch handlers

### Implementation Details

#### Stack Capture Process

1. When `error()` is called, it retrieves the `TraceFunc` from the context
2. The VM's `captureStackTrace()` walks from the current frame pointer (`vm.fp`) down to frame 0
3. For each frame, it extracts:
   - Function name from `frame.fn.Name()` if available
   - Filename from the function's code or frame's code
   - Frame type (function vs module)
4. The resulting stack frames are attached to the error object

#### Filename Propagation

To ensure filenames are available in tracebacks:
- The compiler's `Code.newChild()` method was modified to inherit the parent's filename
- Functions compiled within a file now retain that file's name in their code objects
- The VM checks both function code and frame code for filename information

#### Error Propagation

Key changes to preserve tracebacks during error propagation:
- VM's `callObject()` returns `*object.Error` instead of unwrapping to `error`
- Various VM operations that check `IsRaised()` now return the error object itself
- The `try()` builtin preserves error objects when catching them

### Example Usage

```risor
func inner() {
    error("something went wrong")
}

func outer() {
    inner()
}

try(
    outer,
    func(err) {
        print(err.traceback())
    }
)
```

Output:
```
Traceback (most recent call last):
  File "example.risor", line 0, in inner
  File "example.risor", line 0, in outer
  File "example.risor", line 0, in <module>
Error: something went wrong
```

## Trade-offs and Decisions

### Performance Considerations

- Stack capture adds overhead to error creation, but only when errors are actually raised
- The traceback is captured eagerly (at error creation time) rather than lazily
- No performance impact on the happy path (non-error cases)

### Design Choices

1. **Eager vs Lazy Capture**: We capture the stack immediately when an error is created. This ensures the stack represents the actual error location, not where it's later accessed.

2. **Traceback Format**: We follow Python's traceback format as it's familiar to many developers and clearly shows the call sequence.

3. **Error Object Preservation**: We modified error propagation to preserve `*object.Error` types rather than converting to Go `error` interfaces, ensuring traceback information isn't lost.

4. **Filename Inheritance**: Child code objects (functions) automatically inherit their parent's filename, simplifying the implementation and ensuring consistent file attribution.

## Future Enhancements

1. **Line Number Support**: Currently all line numbers show as 0. Implementing accurate line numbers would require:
   - Tracking source positions during parsing
   - Mapping bytecode instructions to source positions
   - Storing this mapping in compiled code

2. **Column Numbers**: Could add column information for more precise error locations

3. **Stack Frame Variables**: Could capture and display local variables in each frame for better debugging context

4. **Customizable Traceback Format**: Could allow users to customize how tracebacks are displayed

5. **Traceback Filtering**: Could add options to filter or limit traceback depth

## Testing

Comprehensive tests were added to verify:
- Basic traceback capture and display
- Nested function call stacks
- Error re-raising with traceback preservation
- Filename inclusion in tracebacks
- Anonymous function handling
- Module-level error handling

## Conclusion

The error traceback feature significantly improves the debugging experience in Risor by providing clear visibility into the call stack when errors occur. The implementation integrates cleanly with existing error handling mechanisms while maintaining backward compatibility.
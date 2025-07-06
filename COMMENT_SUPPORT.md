# Comment Support in Risor AST

This implementation adds full comment support to the Risor lexer and parser, allowing comments to be captured in the Abstract Syntax Tree (AST). This is essential for building tools like auto-formatters that need to preserve comments in the source code.

## Overview

The implementation captures three types of comments:

1. **Single-line comments using `#`**: `# This is a comment`
2. **Single-line comments using `//`**: `// This is a comment`
3. **Multi-line comments using `/* */`**: `/* This is a multi-line comment */`

Comments are now treated as first-class AST nodes and are preserved during parsing, making them available for tools that need to maintain the original code structure.

## Changes Made

### 1. Token Types (`token/token.go`)

Added two new token types:
- `COMMENT` - for single-line comments (both `#` and `//`)
- `MULTILINE_COMMENT` - for multi-line comments (`/* */`)

### 2. Lexer (`lexer/lexer.go`)

Modified the lexer to capture comments instead of skipping them:
- Replaced `skipComment()` with `readComment()` that returns comment text
- Replaced `skipMultiLineComment()` with `readMultiLineComment()` that returns comment text
- Updated the `Next()` method to return comment tokens instead of skipping them

### 3. AST Node (`ast/statements.go`)

Added a new `Comment` AST node type:
```go
type Comment struct {
    token token.Token
    text  string
}
```

The `Comment` node implements the `Statement` interface and provides:
- `Text()` method to get the comment content
- `String()` method that returns the comment text
- Standard `Token()`, `Literal()`, and other required methods

### 4. Parser (`parser/parser.go`)

Enhanced the parser to handle comment tokens:
- Added `parseComment()` function to create comment AST nodes
- Registered comment token types as prefix parsers
- Updated statement terminators to include comment tokens

## Usage Example

```risor
# This is a single-line comment using hash
// This is a single-line comment using double slash

/* This is a 
   multi-line comment
   that spans several lines */

var x = 5  # End-of-line comment

// Function with comments
func fibonacci(n) {
    # Base case comment
    if n <= 1 {
        return n
    }
    
    /* Recursive case:
       Calculate fibonacci recursively */
    return fibonacci(n-1) + fibonacci(n-2)
}

var result = fibonacci(10)  // Calculate fibonacci of 10
```

## Testing

The implementation includes comprehensive tests:

### Lexer Tests (`lexer/lexer_test.go`)
- Updated existing tests to expect comment tokens instead of skipping them
- Tests for single-line comments, multi-line comments, and shebang handling

### Parser Tests (`parser/parser_test.go`)
- Added `TestComments()` for testing individual comment parsing
- Added `TestCommentsInCode()` for testing comments mixed with code

All tests pass, ensuring backward compatibility while adding the new functionality.

## Benefits for Tool Development

This implementation enables the development of sophisticated Risor tools:

1. **Auto-formatters**: Can preserve comments while reformatting code
2. **Documentation generators**: Can extract comments for documentation
3. **Code analyzers**: Can analyze comment patterns and documentation coverage
4. **IDE support**: Can provide better syntax highlighting and code folding

## Impact on Existing Code

The changes are backward compatible:
- Existing Risor code will continue to work exactly as before
- The parser behavior remains the same for non-comment tokens
- All existing tests pass without modification (except lexer tests that explicitly tested comment skipping)

The only difference is that comments are now preserved in the AST rather than being discarded during lexing.
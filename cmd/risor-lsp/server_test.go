package main

import (
	"context"
	"testing"

	"github.com/jdbaldry/go-language-server-protocol/lsp/protocol"
	"github.com/risor-io/risor/parser"
	"github.com/stretchr/testify/require"
)

// Helper function to set a document in the cache for testing
func setTestDocument(c *cache, uri protocol.DocumentURI, text string) error {
	item := &protocol.TextDocumentItem{
		URI:     uri,
		Text:    text,
		Version: 1,
	}

	doc := &document{
		item:                 *item,
		linesChangedSinceAST: map[int]bool{},
	}

	if text != "" {
		ctx := context.Background()
		doc.ast, doc.err = parser.Parse(ctx, text)
	}

	return c.put(doc)
}

func TestCache_ParseValidRisorCode(t *testing.T) {
	c := newCache()

	// Test valid Risor code
	validCode := `var x = 42
y := "hello"
func add(a, b) {
    return a + b
}`

	uri := protocol.DocumentURI("file:///test.risor")
	err := setTestDocument(c, uri, validCode)
	require.NoError(t, err)

	doc, err := c.get(uri)
	require.NoError(t, err)

	require.NoError(t, doc.err)

	require.NotNil(t, doc.ast)

	// Verify we have statements
	statements := doc.ast.Statements()
	require.NotEmpty(t, statements)
}

func TestCache_ParseInvalidRisorCode(t *testing.T) {
	c := newCache()

	// Test invalid Risor code
	invalidCode := `var x = 
func incomplete(`

	uri := protocol.DocumentURI("file:///test_invalid.risor")
	err := setTestDocument(c, uri, invalidCode)
	require.NoError(t, err)

	doc, err := c.get(uri)
	require.NoError(t, err)

	// Should have a parse error
	require.Error(t, doc.err)
}

func TestCompletionProvider_ExtractVariables(t *testing.T) {
	// Create a test program
	code := `var x = 42
y := "hello"
z = [1, 2, 3]`

	ctx := context.Background()
	prog, err := parser.Parse(ctx, code)
	require.NoError(t, err)

	variables := extractVariables(prog)

	expectedVars := []string{"x", "y", "z"}
	require.Equal(t, len(expectedVars), len(variables))

	// Check that all expected variables are found
	varMap := make(map[string]bool)
	for _, v := range variables {
		varMap[v] = true
	}

	for _, expected := range expectedVars {
		require.True(t, varMap[expected], "Expected variable %s not found in %v", expected, variables)
	}
}

func TestCompletionProvider_ExtractFunctions(t *testing.T) {
	// Create a test program with function assignments
	code := `add := func(a, b) { return a + b }
subtract = func(x, y) { return x - y }`

	ctx := context.Background()
	prog, err := parser.Parse(ctx, code)
	require.NoError(t, err)

	functions := extractFunctions(prog)

	expectedFuncs := []string{"add", "subtract"}
	require.Equal(t, len(expectedFuncs), len(functions))

	// Check that all expected functions are found
	funcMap := make(map[string]bool)
	for _, f := range functions {
		funcMap[f] = true
	}

	for _, expected := range expectedFuncs {
		require.True(t, funcMap[expected], "Expected function %s not found in %v", expected, functions)
	}
}

func TestHoverProvider_FindSymbolAtPosition(t *testing.T) {
	// Create a test program
	code := `var x = 42
y := "hello"`

	ctx := context.Background()
	prog, err := parser.Parse(ctx, code)
	require.NoError(t, err)

	// Test finding symbol at position of variable 'x' (line 1, around column 5)
	symbol := findSymbolAtPosition(prog, 1, 5)
	require.Equal(t, "x", symbol)

	// Test finding symbol at position of variable 'y' (line 2, around column 1)
	symbol = findSymbolAtPosition(prog, 2, 1)
	require.Equal(t, "y", symbol)

	// Test position with no symbol
	symbol = findSymbolAtPosition(prog, 1, 15)
	require.Empty(t, symbol)
}

func TestKeywordsAndBuiltins(t *testing.T) {
	// Test that our keyword list contains expected Risor keywords
	expectedKeywords := []string{"var", "func", "if", "else", "for", "return", "true", "false", "nil"}

	for _, keyword := range expectedKeywords {
		found := false
		for _, k := range risorKeywords {
			if k == keyword {
				found = true
				break
			}
		}
		require.True(t, found, "Expected keyword '%s' not found in risorKeywords", keyword)
	}

	// Test that our builtin list contains expected functions
	expectedBuiltins := []string{"len", "print", "println", "string", "int", "float"}

	for _, builtin := range expectedBuiltins {
		found := false
		for _, b := range risorBuiltins {
			if b == builtin {
				found = true
				break
			}
		}
		require.True(t, found, "Expected builtin '%s' not found in risorBuiltins", builtin)
	}
}

func TestDiagnostics_WithParseError(t *testing.T) {
	// Test code with syntax error
	invalidCode := `var x = 
func incomplete(`

	// Parse the code to get a parse error
	ctx := context.Background()
	_, err := parser.Parse(ctx, invalidCode)
	require.Error(t, err)

	// Verify it's a parse error we can handle
	parseErr, ok := err.(parser.ParserError)
	require.True(t, ok, "Expected parser.ParseError type, got %T", err)

	require.NotEmpty(t, parseErr.Message())

	startPos := parseErr.StartPosition()
	require.Greater(t, startPos.LineNumber(), 0)
}

func TestServer_QueueDiagnostics(t *testing.T) {
	// Create a minimal server for testing
	server := &Server{
		name:    "test-server",
		version: "test",
		cache:   newCache(),
	}

	// This test mainly ensures the method doesn't panic
	// In a full integration test, we'd mock the client and verify the diagnostics
	uri := protocol.DocumentURI("file:///test.risor")

	// Set a document with an error
	err := setTestDocument(server.cache, uri, "var x = \nfunc incomplete(")
	require.NoError(t, err)

	// This should not panic
	server.queueDiagnostics(uri)
}

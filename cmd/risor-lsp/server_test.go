package main

import (
	"context"
	"testing"

	"github.com/jdbaldry/go-language-server-protocol/lsp/protocol"
	"github.com/risor-io/risor/parser"
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
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	
	doc, err := c.get(uri)
	if err != nil {
		t.Fatalf("Expected no error retrieving document, got %v", err)
	}
	
	if doc.err != nil {
		t.Fatalf("Expected no parse error, got %v", doc.err)
	}
	
	if doc.ast == nil {
		t.Fatalf("Expected AST to be parsed, got nil")
	}
	
	// Verify we have statements
	statements := doc.ast.Statements()
	if len(statements) == 0 {
		t.Fatalf("Expected statements in AST, got none")
	}
}

func TestCache_ParseInvalidRisorCode(t *testing.T) {
	c := newCache()
	
	// Test invalid Risor code
	invalidCode := `var x = 
func incomplete(`
	
	uri := protocol.DocumentURI("file:///test_invalid.risor")
	err := setTestDocument(c, uri, invalidCode)
	if err != nil {
		t.Fatalf("Expected no error setting document, got %v", err)
	}
	
	doc, err := c.get(uri)
	if err != nil {
		t.Fatalf("Expected no error retrieving document, got %v", err)
	}
	
	// Should have a parse error
	if doc.err == nil {
		t.Fatalf("Expected parse error for invalid code, got none")
	}
}

func TestCompletionProvider_ExtractVariables(t *testing.T) {
	// Create a test program
	code := `var x = 42
y := "hello"
z = [1, 2, 3]`
	
	ctx := context.Background()
	prog, err := parser.Parse(ctx, code)
	if err != nil {
		t.Fatalf("Failed to parse test code: %v", err)
	}
	
	variables := extractVariables(prog)
	
	expectedVars := []string{"x", "y", "z"}
	if len(variables) != len(expectedVars) {
		t.Fatalf("Expected %d variables, got %d: %v", len(expectedVars), len(variables), variables)
	}
	
	// Check that all expected variables are found
	varMap := make(map[string]bool)
	for _, v := range variables {
		varMap[v] = true
	}
	
	for _, expected := range expectedVars {
		if !varMap[expected] {
			t.Errorf("Expected variable %s not found in %v", expected, variables)
		}
	}
}

func TestCompletionProvider_ExtractFunctions(t *testing.T) {
	// Create a test program with function assignments
	code := `add := func(a, b) { return a + b }
subtract = func(x, y) { return x - y }`
	
	ctx := context.Background()
	prog, err := parser.Parse(ctx, code)
	if err != nil {
		t.Fatalf("Failed to parse test code: %v", err)
	}
	
	functions := extractFunctions(prog)
	
	expectedFuncs := []string{"add", "subtract"}
	if len(functions) != len(expectedFuncs) {
		t.Fatalf("Expected %d functions, got %d: %v", len(expectedFuncs), len(functions), functions)
	}
	
	// Check that all expected functions are found
	funcMap := make(map[string]bool)
	for _, f := range functions {
		funcMap[f] = true
	}
	
	for _, expected := range expectedFuncs {
		if !funcMap[expected] {
			t.Errorf("Expected function %s not found in %v", expected, functions)
		}
	}
}

func TestHoverProvider_FindSymbolAtPosition(t *testing.T) {
	// Create a test program
	code := `var x = 42
y := "hello"`
	
	ctx := context.Background()
	prog, err := parser.Parse(ctx, code)
	if err != nil {
		t.Fatalf("Failed to parse test code: %v", err)
	}
	
	// Test finding symbol at position of variable 'x' (line 1, around column 5)
	symbol := findSymbolAtPosition(prog, 1, 5)
	if symbol != "x" {
		t.Errorf("Expected to find symbol 'x', got '%s'", symbol)
	}
	
	// Test finding symbol at position of variable 'y' (line 2, around column 1)
	symbol = findSymbolAtPosition(prog, 2, 1)
	if symbol != "y" {
		t.Errorf("Expected to find symbol 'y', got '%s'", symbol)
	}
	
	// Test position with no symbol
	symbol = findSymbolAtPosition(prog, 1, 15)
	if symbol != "" {
		t.Errorf("Expected no symbol at position, got '%s'", symbol)
	}
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
		if !found {
			t.Errorf("Expected keyword '%s' not found in risorKeywords", keyword)
		}
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
		if !found {
			t.Errorf("Expected builtin '%s' not found in risorBuiltins", builtin)
		}
	}
}

func TestDiagnostics_WithParseError(t *testing.T) {
	// Test code with syntax error
	invalidCode := `var x = 
func incomplete(`
	
	// Parse the code to get a parse error
	ctx := context.Background()
	_, err := parser.Parse(ctx, invalidCode)
	if err == nil {
		t.Fatalf("Expected parse error for invalid code, got none")
	}
	
	// Verify it's a parse error we can handle
	if parseErr, ok := err.(parser.ParserError); ok {
		if parseErr.Message() == "" {
			t.Errorf("Expected parse error to have a message")
		}
		
		startPos := parseErr.StartPosition()
		if startPos.LineNumber() <= 0 {
			t.Errorf("Expected valid line number in parse error, got %d", startPos.LineNumber())
		}
	} else {
		t.Errorf("Expected parser.ParseError type, got %T", err)
	}
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
	if err != nil {
		t.Fatalf("Failed to set test document: %v", err)
	}
	
	// This should not panic
	server.queueDiagnostics(uri)
}
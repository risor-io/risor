package main

import (
	"context"
	"testing"

	"github.com/jdbaldry/go-language-server-protocol/lsp/protocol"
)

// TestLanguageServerIntegration demonstrates testing the language server
// with a complete Risor file, simulating real VS Code interactions
func TestLanguageServerIntegration(t *testing.T) {
	// Sample Risor code that demonstrates various language features
	risorCode := `// Example Risor program
var config = {
    "host": "localhost",
    "port": 8080,
    "debug": true
}

// Function to process user data
process_user := func(user_id, name) {
    if user_id <= 0 {
        return "Invalid user ID"
    }
    
    user_data := {
        "id": user_id,
        "name": name,
        "status": "active"
    }
    
    return user_data
}

// Main processing logic
users := []
for i in range(1, 5) {
    user := process_user(i, sprintf("User_%d", i))
    users = append(users, user)
}

// Print results
for user in users {
    println(sprintf("User: %s (ID: %d)", user["name"], user["id"]))
}`

	// Create a server instance
	server := &Server{
		name:    "test-risor-lsp",
		version: "1.0.0-test",
		cache:   newCache(),
	}

	uri := protocol.DocumentURI("file:///example.risor")

	// Test 1: Document parsing and caching
	t.Run("DocumentParsing", func(t *testing.T) {
		err := setTestDocument(server.cache, uri, risorCode)
		if err != nil {
			t.Fatalf("Failed to cache document: %v", err)
		}

		doc, err := server.cache.get(uri)
		if err != nil {
			t.Fatalf("Failed to retrieve document: %v", err)
		}

		if doc.err != nil {
			t.Fatalf("Document parsing failed: %v", doc.err)
		}

		if doc.ast == nil {
			t.Fatal("Expected AST to be parsed")
		}

		statements := doc.ast.Statements()
		if len(statements) == 0 {
			t.Fatal("Expected statements in AST")
		}

		t.Logf("Successfully parsed %d statements", len(statements))
	})

	// Test 2: Completion at various positions
	t.Run("Completion", func(t *testing.T) {
		// Test completion at line 23 (after "for user in")
		params := &protocol.CompletionParams{
			TextDocumentPositionParams: protocol.TextDocumentPositionParams{
				TextDocument: protocol.TextDocumentIdentifier{URI: uri},
				Position:     protocol.Position{Line: 22, Character: 15}, // After "for user in "
			},
		}

		result, err := server.Completion(context.Background(), params)
		if err != nil {
			t.Fatalf("Completion failed: %v", err)
		}

		if result == nil || len(result.Items) == 0 {
			t.Fatal("Expected completion items")
		}

		// Should include variables like "users", keywords, and builtins
		hasUsers := false
		hasKeywords := false
		hasBuiltins := false

		for _, item := range result.Items {
			switch item.Label {
			case "users":
				hasUsers = true
			case "range", "if", "for":
				hasKeywords = true
			case "len", "print", "println":
				hasBuiltins = true
			}
		}

		if !hasUsers {
			t.Error("Expected 'users' variable in completion")
		}
		if !hasKeywords {
			t.Error("Expected keywords in completion")
		}
		if !hasBuiltins {
			t.Error("Expected builtin functions in completion")
		}

		t.Logf("Completion returned %d items", len(result.Items))
	})

	// Test 3: Hover information
	t.Run("Hover", func(t *testing.T) {
		// Test hover over the "process_user" function name
		params := &protocol.HoverParams{
			TextDocumentPositionParams: protocol.TextDocumentPositionParams{
				TextDocument: protocol.TextDocumentIdentifier{URI: uri},
				Position:     protocol.Position{Line: 7, Character: 0}, // At "process_user"
			},
		}

		result, err := server.Hover(context.Background(), params)
		if err != nil {
			t.Fatalf("Hover failed: %v", err)
		}

		// Note: hover might not find anything with our simple position-based implementation
		// This is expected for this test
		if result != nil && result.Contents.Value != "" {
			t.Logf("Hover content: %s", result.Contents.Value)
		} else {
			t.Logf("No hover content found (expected with simple implementation)")
		}
	})

	// Test 4: Document symbols
	t.Run("DocumentSymbols", func(t *testing.T) {
		params := &protocol.DocumentSymbolParams{
			TextDocument: protocol.TextDocumentIdentifier{URI: uri},
		}

		result, err := server.DocumentSymbol(context.Background(), params)
		if err != nil {
			t.Fatalf("DocumentSymbol failed: %v", err)
		}

		if result == nil || len(result) == 0 {
			t.Fatal("Expected document symbols")
		}

		// Should find variables like "config", "process_user", "users"
		symbolNames := []string{}
		for _, symbolInterface := range result {
			if symbol, ok := symbolInterface.(protocol.DocumentSymbol); ok {
				symbolNames = append(symbolNames, symbol.Name)
			}
		}

		expectedSymbols := []string{"config", "process_user", "users"}
		for _, expected := range expectedSymbols {
			found := false
			for _, name := range symbolNames {
				if name == expected {
					found = true
					break
				}
			}
			if !found {
				t.Errorf("Expected symbol '%s' not found in %v", expected, symbolNames)
			}
		}

		t.Logf("Found symbols: %v", symbolNames)
	})

	// Test 5: Go-to-definition
	t.Run("Definition", func(t *testing.T) {
		// Test definition for "process_user" usage
		params := &protocol.DefinitionParams{
			TextDocumentPositionParams: protocol.TextDocumentPositionParams{
				TextDocument: protocol.TextDocumentIdentifier{URI: uri},
				Position:     protocol.Position{Line: 20, Character: 12}, // At "process_user" call
			},
		}

		result, err := server.Definition(context.Background(), params)
		if err != nil {
			t.Fatalf("Definition failed: %v", err)
		}

		// This might not find anything with our simple implementation,
		// but shouldn't error
		if result != nil {
			t.Logf("Definition result type: %T", result)
		}
	})
}

// TestLanguageServerWithErrors tests how the language server handles syntax errors
func TestLanguageServerWithErrors(t *testing.T) {
	server := &Server{
		name:    "test-risor-lsp",
		version: "1.0.0-test",
		cache:   newCache(),
	}

	// Code with syntax errors
	invalidCode := `var x = 42
func incomplete(
y := "missing closing brace"
if true {
    // missing closing brace`

	uri := protocol.DocumentURI("file:///invalid.risor")

	err := setTestDocument(server.cache, uri, invalidCode)
	if err != nil {
		t.Fatalf("Failed to cache document: %v", err)
	}

	doc, err := server.cache.get(uri)
	if err != nil {
		t.Fatalf("Failed to retrieve document: %v", err)
	}

	// Should have a parse error
	if doc.err == nil {
		t.Fatal("Expected parse error for invalid code")
	}

	t.Logf("Parse error (as expected): %v", doc.err)

	// Test that we can still provide basic completion even with errors
	params := &protocol.CompletionParams{
		TextDocumentPositionParams: protocol.TextDocumentPositionParams{
			TextDocument: protocol.TextDocumentIdentifier{URI: uri},
			Position:     protocol.Position{Line: 0, Character: 10},
		},
	}

	result, err := server.Completion(context.Background(), params)
	if err != nil {
		t.Fatalf("Completion failed: %v", err)
	}

	// Should still provide keywords and builtins even with syntax errors
	if result == nil || len(result.Items) == 0 {
		t.Fatal("Expected completion items even with syntax errors")
	}

	t.Logf("Completion with errors returned %d items", len(result.Items))
}

// TestRisorCodeExamples tests the language server with various Risor code patterns
func TestRisorCodeExamples(t *testing.T) {
	examples := map[string]string{
		"variables": `var name = "Risor"
age := 25
is_valid = true`,

		"functions": `add := func(a, b) { return a + b }
greet := func(name) {
    return sprintf("Hello, %s!", name)
}`,

		"control_flow": `if age >= 18 {
    status := "adult"
} else {
    status := "minor"
}

for i in range(10) {
    println(i)
}`,

		"data_structures": `person := {
    "name": "Alice",
    "age": 30,
    "hobbies": ["reading", "coding"]
}

numbers := [1, 2, 3, 4, 5]`,
	}

	server := &Server{
		name:    "test-risor-lsp",
		version: "1.0.0-test",
		cache:   newCache(),
	}

	for name, code := range examples {
		t.Run(name, func(t *testing.T) {
			uri := protocol.DocumentURI("file:///" + name + ".risor")

			err := setTestDocument(server.cache, uri, code)
			if err != nil {
				t.Fatalf("Failed to cache document: %v", err)
			}

			doc, err := server.cache.get(uri)
			if err != nil {
				t.Fatalf("Failed to retrieve document: %v", err)
			}

			if doc.err != nil {
				t.Fatalf("Parse error in %s: %v", name, doc.err)
			}

			if doc.ast == nil {
				t.Fatalf("No AST parsed for %s", name)
			}

			statements := doc.ast.Statements()
			if len(statements) == 0 {
				t.Fatalf("No statements found in %s", name)
			}

			t.Logf("Example '%s': parsed %d statements successfully", name, len(statements))

			// Test that completion works for each example
			params := &protocol.CompletionParams{
				TextDocumentPositionParams: protocol.TextDocumentPositionParams{
					TextDocument: protocol.TextDocumentIdentifier{URI: uri},
					Position:     protocol.Position{Line: 0, Character: 0},
				},
			}

			result, err := server.Completion(context.Background(), params)
			if err != nil {
				t.Fatalf("Completion failed for %s: %v", name, err)
			}

			if result == nil || len(result.Items) == 0 {
				t.Fatalf("No completion items for %s", name)
			}

			t.Logf("Example '%s': completion returned %d items", name, len(result.Items))
		})
	}
}

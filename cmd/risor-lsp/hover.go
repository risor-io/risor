package main

import (
	"context"
	"fmt"

	"github.com/jdbaldry/go-language-server-protocol/lsp/protocol"
	"github.com/risor-io/risor/ast"
	"github.com/rs/zerolog/log"
)

func (s *Server) Hover(ctx context.Context, params *protocol.HoverParams) (*protocol.Hover, error) {
	doc, err := s.cache.get(params.TextDocument.URI)
	if err != nil {
		log.Error().Err(err).Str("call", "Hover").Msg("failed to get document")
		return nil, nil
	}

	if doc.ast == nil || doc.err != nil {
		return nil, nil
	}

	// Convert LSP position to 1-based line/column
	line := int(params.Position.Line) + 1
	column := int(params.Position.Character) + 1

	// Find the symbol at the cursor position
	symbol := findSymbolAtPosition(doc.ast, line, column)
	if symbol == "" {
		return nil, nil
	}

	// Generate hover information for the symbol
	info := getSymbolInfo(symbol, doc.ast)
	if info == "" {
		// Check if it's a built-in function
		if contains(risorBuiltins, symbol) {
			info = fmt.Sprintf("**%s()** - Built-in function", symbol)
		} else if contains(risorKeywords, symbol) {
			info = fmt.Sprintf("**%s** - Risor keyword", symbol)
		} else {
			return nil, nil
		}
	}

	return &protocol.Hover{
		Contents: protocol.MarkupContent{
			Kind:  protocol.Markdown,
			Value: info,
		},
	}, nil
}

// findSymbolAtPosition finds the symbol (identifier) at the given line/column position
func findSymbolAtPosition(program *ast.Program, line, column int) string {
	// Simple implementation: look through statements for identifiers at the position
	for _, stmt := range program.Statements() {
		if symbol := findIdentInNode(stmt, line, column); symbol != "" {
			return symbol
		}
	}
	return ""
}

// findIdentInNode recursively searches for identifiers in a node
func findIdentInNode(node ast.Node, line, column int) string {
	switch n := node.(type) {
	case *ast.Var:
		name, _ := n.Value()
		pos := n.Token().StartPosition
		endPos := n.Token().EndPosition
		if pos.LineNumber() == line && pos.ColumnNumber() <= column && column <= endPos.ColumnNumber() {
			return name
		}
	case *ast.Assign:
		name := n.Name()
		pos := n.Token().StartPosition
		endPos := n.Token().EndPosition
		if pos.LineNumber() == line && pos.ColumnNumber() <= column && column <= endPos.ColumnNumber() {
			return name
		}
	case *ast.Func:
		if n.Name() != nil {
			name := n.Name().String()
			pos := n.Token().StartPosition
			endPos := n.Token().EndPosition
			if pos.LineNumber() == line && pos.ColumnNumber() <= column && column <= endPos.ColumnNumber() {
				return name
			}
		}
	}
	return ""
}

// getSymbolInfo returns documentation or information about a symbol
func getSymbolInfo(symbol string, program *ast.Program) string {
	// Check variables and functions in all statements
	for _, stmt := range program.Statements() {
		switch s := stmt.(type) {
		case *ast.Assign:
			if s.Name() == symbol {
				return fmt.Sprintf("**%s** - Variable\n\nAssigned in this file", symbol)
			}
		case *ast.Var:
			name, _ := s.Value()
			if name == symbol {
				return fmt.Sprintf("**%s** - Variable\n\nDeclared in this file", symbol)
			}
		// Check if statement contains expressions with functions (e.g., function assignments)
		default:
			// For now, we don't traverse into expressions within statements
			// This could be enhanced later to find functions assigned to variables
		}
	}

	return ""
}

// contains checks if a slice contains a string
func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}

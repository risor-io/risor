package semver

import (
	"context"
	"testing"

	"github.com/jdbaldry/go-language-server-protocol/lsp/protocol"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestServer_Initialize(t *testing.T) {
	s := &Server{
		name:    "Test Server",
		version: "1.0.0",
		cache:   newCache(),
	}

	ctx := context.Background()
	params := &protocol.ParamInitialize{}

	result, err := s.Initialize(ctx, params)
	require.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, s.name, result.ServerInfo.Name)
	assert.Equal(t, s.version, result.ServerInfo.Version)
	assert.True(t, result.Capabilities.HoverProvider)
	assert.True(t, result.Capabilities.DefinitionProvider)
	assert.True(t, result.Capabilities.DocumentFormattingProvider)
	assert.True(t, result.Capabilities.DocumentSymbolProvider)
}

func TestServer_DidOpen(t *testing.T) {
	s := &Server{
		cache: newCache(),
	}

	ctx := context.Background()
	params := &protocol.DidOpenTextDocumentParams{
		TextDocument: protocol.TextDocumentItem{
			URI:  protocol.DocumentURI("file:///test.risor"),
			Text: "x := 1",
		},
	}

	err := s.DidOpen(ctx, params)
	require.NoError(t, err)

	doc, err := s.cache.get(params.TextDocument.URI)
	require.NoError(t, err)
	assert.NotNil(t, doc)
	assert.NotNil(t, doc.ast)
	assert.NoError(t, doc.err)
}

func TestServer_DidChange(t *testing.T) {
	s := &Server{
		cache: newCache(),
	}

	ctx := context.Background()
	params := &protocol.DidChangeTextDocumentParams{
		TextDocument: protocol.VersionedTextDocumentIdentifier{
			TextDocumentIdentifier: protocol.TextDocumentIdentifier{
				URI: protocol.DocumentURI("file:///test.risor"),
			},
		},
	}

	err := s.DidChange(ctx, params)
	require.NoError(t, err)
	// Add more assertions if DidChange is implemented to do more
}

// Helper function to create a new cache for testing
func newCache() *cache {
	return &cache{
		documents: make(map[protocol.DocumentURI]*document),
	}
}

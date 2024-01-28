package main

import (
	"context"

	"github.com/jdbaldry/go-language-server-protocol/lsp/protocol"
)

func (s *Server) Hover(ctx context.Context, params *protocol.HoverParams) (*protocol.Hover, error) {
	var name string
	var doc *document
	for k, v := range s.cache.docs {
		doc = v
		name = k.SpanURI().Filename()
		break
	}

	if doc == nil {
		return nil, nil
	}

	return &protocol.Hover{
		Range: protocol.Range{
			Start: protocol.Position{Line: 1, Character: 1},
			End:   protocol.Position{Line: 2, Character: 1},
		},
		Contents: protocol.MarkupContent{
			Kind:  protocol.Markdown,
			Value: name,
		},
	}, nil
}

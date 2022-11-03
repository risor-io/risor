package main

import (
	"context"

	"github.com/jdbaldry/go-language-server-protocol/lsp/protocol"
	"github.com/myzie/tamarin/ast"
	"github.com/rs/zerolog/log"
)

func (s *Server) DocumentSymbol(ctx context.Context, params *protocol.DocumentSymbolParams) ([]interface{}, error) {
	doc, err := s.cache.get(params.TextDocument.URI)
	if err != nil {
		log.Error().Err(err).Str("call", "DocumentSymbol").Msg("failed to get document")
		return nil, nil
	}
	if doc.err != nil {
		log.Error().Err(err).Str("call", "DocumentSymbol").Msg("document has error")
		return nil, nil
	}

	var symbols []protocol.DocumentSymbol
	for i, stmt := range doc.ast.Statements {
		switch stmt := stmt.(type) {
		case *ast.LetStatement:
			symbols = append(symbols, protocol.DocumentSymbol{
				Name: stmt.Name.Value,
				Kind: protocol.Variable,
				Range: protocol.Range{
					Start: protocol.Position{
						Line:      uint32(i),
						Character: uint32(1),
					},
					End: protocol.Position{
						Line:      uint32(i),
						Character: uint32(5),
					},
				},
				Detail: "Let Statement",
			})
			// log.Info().
			// 	Str("call", "DocumentSymbol").
			// 	Str("name", stmt.Name.Value).
			// 	Str("stmt", stmt.String()).
			// 	Msg("let statement found")
		}
	}
	log.Info().
		Str("call", "DocumentSymbol").
		Int("count", len(symbols)).
		Msg("document statement found")

	result := make([]interface{}, len(symbols))
	for i, symbol := range symbols {
		result[i] = symbol
	}
	return result, nil
}

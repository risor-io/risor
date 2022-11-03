package main

import (
	"context"

	"github.com/jdbaldry/go-language-server-protocol/lsp/protocol"
)

func (s *Server) Formatting(ctx context.Context, params *protocol.DocumentFormattingParams) ([]protocol.TextEdit, error) {
	return nil, nil
}

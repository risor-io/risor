package main

import (
	"context"

	"github.com/jdbaldry/go-language-server-protocol/lsp/protocol"
)

func (s *Server) Completion(ctx context.Context, params *protocol.CompletionParams) (*protocol.CompletionList, error) {
	return &protocol.CompletionList{IsIncomplete: false, Items: nil}, nil
}

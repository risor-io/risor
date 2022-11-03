package main

import (
	"context"

	"github.com/jdbaldry/go-language-server-protocol/lsp/protocol"
)

type Configuration struct {
	EnableEvalDiagnostics bool
	EnableLintDiagnostics bool
}

func (s *Server) DidChangeConfiguration(ctx context.Context, params *protocol.DidChangeConfigurationParams) error {
	return nil
}

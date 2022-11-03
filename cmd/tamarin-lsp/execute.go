package main

import (
	"context"
	"fmt"

	"github.com/jdbaldry/go-language-server-protocol/lsp/protocol"
)

func (s *Server) ExecuteCommand(ctx context.Context, params *protocol.ExecuteCommandParams) (interface{}, error) {
	return nil, fmt.Errorf("unknown command: %s", params.Command)
}

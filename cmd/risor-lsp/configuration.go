package main

import (
	"context"

	"github.com/jdbaldry/go-language-server-protocol/lsp/protocol"
	"github.com/rs/zerolog/log"
)

type Configuration struct {
	EnableEvalDiagnostics bool `json:"enableEvalDiagnostics"`
	EnableLintDiagnostics bool `json:"enableLintDiagnostics"`
	MaxNumberOfProblems   int  `json:"maxNumberOfProblems"`
}

func (s *Server) DidChangeConfiguration(ctx context.Context, params *protocol.DidChangeConfigurationParams) error {
	log.Info().Msg("Configuration changed")
	
	// For now, we'll use default settings
	// In a full implementation, you'd parse params.Settings to extract configuration
	config := Configuration{
		EnableEvalDiagnostics: false,
		EnableLintDiagnostics: true,
		MaxNumberOfProblems:   100,
	}
	
	// Store configuration in server (you'd add a config field to Server struct)
	log.Info().
		Bool("evalDiagnostics", config.EnableEvalDiagnostics).
		Bool("lintDiagnostics", config.EnableLintDiagnostics).
		Int("maxProblems", config.MaxNumberOfProblems).
		Msg("Updated configuration")
	
	return nil
}

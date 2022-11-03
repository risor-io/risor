package main

import (
	"context"

	"github.com/jdbaldry/go-language-server-protocol/lsp/protocol"
	"github.com/myzie/tamarin/internal/parser"
	"github.com/rs/zerolog/log"
)

type Server struct {
	name    string
	version string
	client  protocol.ClientCloser
	cache   *cache
}

func (s *Server) queueDiagnostics(uri protocol.DocumentURI) {}

func (s *Server) DidChange(ctx context.Context, params *protocol.DidChangeTextDocumentParams) error {
	defer s.queueDiagnostics(params.TextDocument.URI)
	return nil
}

func (s *Server) DidOpen(ctx context.Context, params *protocol.DidOpenTextDocumentParams) (err error) {
	defer s.queueDiagnostics(params.TextDocument.URI)
	log.Info().Str("filename", params.TextDocument.URI.SpanURI().Filename()).Msg("DidOpen")
	doc := &document{
		item:                 params.TextDocument,
		linesChangedSinceAST: map[int]bool{},
	}
	if params.TextDocument.Text != "" {
		doc.ast, doc.err = parser.ParseProgram(params.TextDocument.Text)
		if doc.err != nil {
			log.Error().Err(doc.err).Msg("parse program failed")
		} else {
			log.Info().Msg("parse program ok")
		}
	}
	return s.cache.put(doc)
}

func (s *Server) Initialize(ctx context.Context, params *protocol.ParamInitialize) (*protocol.InitializeResult, error) {
	log.Info().Msg("Initialize")
	return &protocol.InitializeResult{
		Capabilities: protocol.ServerCapabilities{
			CompletionProvider: protocol.CompletionOptions{
				TriggerCharacters: []string{"."},
			},
			HoverProvider:              true,
			DefinitionProvider:         true,
			DocumentFormattingProvider: true,
			DocumentSymbolProvider:     true,
			ExecuteCommandProvider: protocol.ExecuteCommandOptions{
				Commands: []string{},
			},
			TextDocumentSync: &protocol.TextDocumentSyncOptions{
				Change:    protocol.Full,
				OpenClose: true,
				Save: protocol.SaveOptions{
					IncludeText: false,
				},
			},
		},
		ServerInfo: struct {
			Name    string `json:"name"`
			Version string `json:"version,omitempty"`
		}{
			Name:    s.name,
			Version: s.version,
		},
	}, nil
}

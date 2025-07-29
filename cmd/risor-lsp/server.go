package main

import (
	"context"

	"github.com/jdbaldry/go-language-server-protocol/lsp/protocol"
	"github.com/risor-io/risor/parser"
	"github.com/rs/zerolog/log"
)

type Server struct {
	name    string
	version string
	client  protocol.ClientCloser
	cache   *cache
}

func (s *Server) queueDiagnostics(uri protocol.DocumentURI) {
	go s.publishDiagnostics(uri)
}

func (s *Server) publishDiagnostics(uri protocol.DocumentURI) {
	doc, err := s.cache.get(uri)
	if err != nil {
		log.Error().Err(err).Str("uri", string(uri)).Msg("failed to get document for diagnostics")
		return
	}

	var diagnostics []protocol.Diagnostic

	// Check for parse errors
	if doc.err != nil {
		if parseErr, ok := doc.err.(parser.ParserError); ok {
			startPos := parseErr.StartPosition()
			endPos := parseErr.EndPosition()

			diagnostic := protocol.Diagnostic{
				Range: protocol.Range{
					Start: protocol.Position{
						Line:      uint32(startPos.LineNumber() - 1),   // LSP uses 0-based line numbers
						Character: uint32(startPos.ColumnNumber() - 1), // LSP uses 0-based column numbers
					},
					End: protocol.Position{
						Line:      uint32(endPos.LineNumber() - 1),
						Character: uint32(endPos.ColumnNumber() - 1),
					},
				},
				Severity: 1, // Error
				Source:   "risor-lsp",
				Message:  parseErr.Message(),
			}
			diagnostics = append(diagnostics, diagnostic)
		} else {
			// Generic error handling for non-parser errors
			diagnostic := protocol.Diagnostic{
				Range: protocol.Range{
					Start: protocol.Position{Line: 0, Character: 0},
					End:   protocol.Position{Line: 0, Character: 0},
				},
				Severity: 1, // Error
				Source:   "risor-lsp",
				Message:  doc.err.Error(),
			}
			diagnostics = append(diagnostics, diagnostic)
		}
	}

	// Publish diagnostics to the client
	params := &protocol.PublishDiagnosticsParams{
		URI:         uri,
		Diagnostics: diagnostics,
	}

	err = s.client.PublishDiagnostics(context.Background(), params)
	if err != nil {
		log.Error().Err(err).Str("uri", string(uri)).Msg("failed to publish diagnostics")
	}
}

func (s *Server) DidChange(ctx context.Context, params *protocol.DidChangeTextDocumentParams) error {
	doc, err := s.cache.get(params.TextDocument.URI)
	if err != nil {
		log.Error().Err(err).Str("uri", string(params.TextDocument.URI)).Msg("failed to get document for change")
		return err
	}

	// Apply changes to the document
	for _, change := range params.ContentChanges {
		if change.Range == nil {
			// Full document update
			doc.item.Text = change.Text
		} else {
			// Incremental update - for now, we'll treat it as full update
			// TODO: Implement proper incremental updates
			doc.item.Text = change.Text
		}
	}

	// Update version
	doc.item.Version = params.TextDocument.Version

	// Reparse the document
	doc.ast, doc.err = parser.Parse(ctx, doc.item.Text)
	if doc.err != nil {
		log.Error().Err(doc.err).Msg("parse program failed after change")
	} else {
		log.Info().Msg("parse program ok after change")
	}

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
		doc.ast, doc.err = parser.Parse(ctx, params.TextDocument.Text)
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

func (s *Server) DidSave(ctx context.Context, params *protocol.DidSaveTextDocumentParams) error {
	log.Info().Str("filename", params.TextDocument.URI.SpanURI().Filename()).Msg("DidSave")
	defer s.queueDiagnostics(params.TextDocument.URI)
	return nil
}

func (s *Server) DidClose(ctx context.Context, params *protocol.DidCloseTextDocumentParams) error {
	log.Info().Str("filename", params.TextDocument.URI.SpanURI().Filename()).Msg("DidClose")
	// Clear diagnostics for closed document
	err := s.client.PublishDiagnostics(context.Background(), &protocol.PublishDiagnosticsParams{
		URI:         params.TextDocument.URI,
		Diagnostics: []protocol.Diagnostic{},
	})
	if err != nil {
		log.Error().Err(err).Msg("failed to clear diagnostics")
	}
	return nil
}

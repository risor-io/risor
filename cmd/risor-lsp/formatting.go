package main

import (
	"context"
	"strings"

	"github.com/jdbaldry/go-language-server-protocol/lsp/protocol"
	"github.com/rs/zerolog/log"
)

func (s *Server) Formatting(ctx context.Context, params *protocol.DocumentFormattingParams) ([]protocol.TextEdit, error) {
	doc, err := s.cache.get(params.TextDocument.URI)
	if err != nil {
		log.Error().Err(err).Str("call", "Formatting").Msg("failed to get document")
		return nil, nil
	}

	// For now, implement basic formatting: ensure proper indentation and spacing
	text := doc.item.Text
	lines := strings.Split(text, "\n")

	var formattedLines []string
	indentLevel := 0

	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if trimmed == "" {
			formattedLines = append(formattedLines, "")
			continue
		}

		// Decrease indent for closing braces/brackets
		if strings.HasPrefix(trimmed, "}") || strings.HasPrefix(trimmed, "]") {
			if indentLevel > 0 {
				indentLevel--
			}
		}

		// Apply indentation
		indent := strings.Repeat("  ", indentLevel) // 2 spaces per level
		formattedLine := indent + trimmed
		formattedLines = append(formattedLines, formattedLine)

		// Increase indent for opening braces/brackets
		if strings.HasSuffix(trimmed, "{") || strings.HasSuffix(trimmed, "[") {
			indentLevel++
		}
	}

	formattedText := strings.Join(formattedLines, "\n")

	// If no changes needed, return nil
	if formattedText == text {
		return nil, nil
	}

	// Calculate the range for the entire document
	lines = strings.Split(text, "\n")
	lastLine := len(lines) - 1
	lastChar := len(lines[lastLine])

	return []protocol.TextEdit{
		{
			Range: protocol.Range{
				Start: protocol.Position{Line: 0, Character: 0},
				End:   protocol.Position{Line: uint32(lastLine), Character: uint32(lastChar)},
			},
			NewText: formattedText,
		},
	}, nil
}

package main

import (
	"context"
	"fmt"

	"github.com/jdbaldry/go-language-server-protocol/jsonrpc2"
	"github.com/jdbaldry/go-language-server-protocol/lsp/protocol"
	"github.com/rs/zerolog/log"
)

func (s *Server) Initialized(context.Context, *protocol.InitializedParams) error {
	return nil
}

func (s *Server) CodeAction(context.Context, *protocol.CodeActionParams) ([]protocol.CodeAction, error) {
	return nil, notImplemented("CodeAction")
}

func (s *Server) CodeLens(ctx context.Context, params *protocol.CodeLensParams) ([]protocol.CodeLens, error) {
	return []protocol.CodeLens{}, nil
}

func (s *Server) CodeLensRefresh(context.Context) error {
	return notImplemented("CodeLensRefresh")
}

func (s *Server) ColorPresentation(context.Context, *protocol.ColorPresentationParams) ([]protocol.ColorPresentation, error) {
	return nil, notImplemented("ColorPresentation")
}

func (s *Server) Declaration(context.Context, *protocol.DeclarationParams) (protocol.Declaration, error) {
	return nil, notImplemented("Declaration")
}

func (s *Server) DidRenameFiles(context.Context, *protocol.RenameFilesParams) error {
	return notImplemented("DidRenameFiles")
}

func (s *Server) DocumentColor(context.Context, *protocol.DocumentColorParams) ([]protocol.ColorInformation, error) {
	return nil, notImplemented("DocumentColor")
}

func (s *Server) DocumentHighlight(context.Context, *protocol.DocumentHighlightParams) ([]protocol.DocumentHighlight, error) {
	return nil, notImplemented("DocumentHighlight")
}

func (s *Server) Exit(context.Context) error {
	return notImplemented("Exit")
}

func (s *Server) FoldingRange(context.Context, *protocol.FoldingRangeParams) ([]protocol.FoldingRange, error) {
	return nil, notImplemented("FoldingRange")
}

func (s *Server) Implementation(context.Context, *protocol.ImplementationParams) (protocol.Definition, error) {
	return nil, notImplemented("Implementation")
}

func (s *Server) IncomingCalls(context.Context, *protocol.CallHierarchyIncomingCallsParams) ([]protocol.CallHierarchyIncomingCall, error) {
	return nil, notImplemented("IncomingCalls")
}

func (s *Server) LinkedEditingRange(context.Context, *protocol.LinkedEditingRangeParams) (*protocol.LinkedEditingRanges, error) {
	return nil, notImplemented("LinkedEditingRange")
}

func (s *Server) LogTrace(context.Context, *protocol.LogTraceParams) error {
	return notImplemented("LogTrace")
}

func (s *Server) Moniker(context.Context, *protocol.MonikerParams) ([]protocol.Moniker, error) {
	return nil, notImplemented("Moniker")
}

func (s *Server) NonstandardRequest(context.Context, string, interface{}) (interface{}, error) {
	return nil, notImplemented("NonstandardRequest")
}

func (s *Server) OnTypeFormatting(context.Context, *protocol.DocumentOnTypeFormattingParams) ([]protocol.TextEdit, error) {
	return nil, notImplemented("OnTypeFormatting")
}

func (s *Server) OutgoingCalls(context.Context, *protocol.CallHierarchyOutgoingCallsParams) ([]protocol.CallHierarchyOutgoingCall, error) {
	return nil, notImplemented("OutgoingCalls")
}

func (s *Server) PrepareCallHierarchy(context.Context, *protocol.CallHierarchyPrepareParams) ([]protocol.CallHierarchyItem, error) {
	return nil, notImplemented("PrepareCallHierarchy")
}

func (s *Server) PrepareRename(context.Context, *protocol.PrepareRenameParams) (*protocol.Range, error) {
	return nil, notImplemented("PrepareRange")
}

func (s *Server) PrepareTypeHierarchy(context.Context, *protocol.TypeHierarchyPrepareParams) ([]protocol.TypeHierarchyItem, error) {
	return nil, notImplemented("PrepareTypeHierarchy")
}

func (s *Server) RangeFormatting(context.Context, *protocol.DocumentRangeFormattingParams) ([]protocol.TextEdit, error) {
	return nil, notImplemented("RangeFormatting")
}

func (s *Server) References(context.Context, *protocol.ReferenceParams) ([]protocol.Location, error) {
	return nil, notImplemented("References")
}

func (s *Server) Rename(context.Context, *protocol.RenameParams) (*protocol.WorkspaceEdit, error) {
	return nil, notImplemented("Rename")
}

func (s *Server) Resolve(context.Context, *protocol.CompletionItem) (*protocol.CompletionItem, error) {
	return nil, notImplemented("Resolve")
}

func (s *Server) ResolveCodeAction(context.Context, *protocol.CodeAction) (*protocol.CodeAction, error) {
	return nil, notImplemented("ResolveCodeAction")
}

func (s *Server) ResolveCodeLens(context.Context, *protocol.CodeLens) (*protocol.CodeLens, error) {
	return nil, notImplemented("ResolveCodeLens")
}

func (s *Server) ResolveDocumentLink(context.Context, *protocol.DocumentLink) (*protocol.DocumentLink, error) {
	return nil, notImplemented("ResolveDocumentLink")
}

func (s *Server) SelectionRange(context.Context, *protocol.SelectionRangeParams) ([]protocol.SelectionRange, error) {
	return nil, notImplemented("SelectionRange")
}

func (s *Server) SemanticTokensFull(ctx context.Context, p *protocol.SemanticTokensParams) (*protocol.SemanticTokens, error) {
	log.Info().
		Str("call", "SemanticTokensFull").
		Interface("params", p).
		Msg("semantic tokens")
	return nil, notImplemented("SemanticTokensFull")
}

func (s *Server) SemanticTokensFullDelta(context.Context, *protocol.SemanticTokensDeltaParams) (interface{}, error) {
	return nil, notImplemented("SemanticTokensFullDelta")
}

func (s *Server) SemanticTokensRange(context.Context, *protocol.SemanticTokensRangeParams) (*protocol.SemanticTokens, error) {
	return nil, notImplemented("SemanticTokensRange")
}

func (s *Server) SemanticTokensRefresh(context.Context) error {
	return notImplemented("SemanticTokensRefresh")
}

func (s *Server) SetTrace(context.Context, *protocol.SetTraceParams) error {
	return notImplemented("SetTrace")
}

func (s *Server) Shutdown(context.Context) error {
	return nil
}

func (s *Server) SignatureHelp(context.Context, *protocol.SignatureHelpParams) (*protocol.SignatureHelp, error) {
	return nil, notImplemented("SignatureHelp")
}

func (s *Server) Subtypes(context.Context, *protocol.TypeHierarchySubtypesParams) ([]protocol.TypeHierarchyItem, error) {
	return nil, notImplemented("Subtypes")
}

func (s *Server) Supertypes(context.Context, *protocol.TypeHierarchySupertypesParams) ([]protocol.TypeHierarchyItem, error) {
	return nil, notImplemented("Supertypes")
}

func (s *Server) Symbol(context.Context, *protocol.WorkspaceSymbolParams) ([]protocol.SymbolInformation, error) {
	return nil, notImplemented("Symbol")
}

func (s *Server) TypeDefinition(context.Context, *protocol.TypeDefinitionParams) (protocol.Definition, error) {
	return nil, notImplemented("TypeDefinition")
}

func (s *Server) WillCreateFiles(context.Context, *protocol.CreateFilesParams) (*protocol.WorkspaceEdit, error) {
	return nil, notImplemented("WillCreateFiles")
}

func (s *Server) WillDeleteFiles(context.Context, *protocol.DeleteFilesParams) (*protocol.WorkspaceEdit, error) {
	return nil, notImplemented("WillDeleteFiles")
}

func (s *Server) WillRenameFiles(context.Context, *protocol.RenameFilesParams) (*protocol.WorkspaceEdit, error) {
	return nil, notImplemented("WillRenameFiles")
}

func (s *Server) WillSave(context.Context, *protocol.WillSaveTextDocumentParams) error {
	return notImplemented("WillSave")
}

func (s *Server) WillSaveWaitUntil(context.Context, *protocol.WillSaveTextDocumentParams) ([]protocol.TextEdit, error) {
	return nil, notImplemented("WillSaveWaitUntil")
}

func (s *Server) WorkDoneProgressCancel(context.Context, *protocol.WorkDoneProgressCancelParams) error {
	return notImplemented("WorkDoneProgressCancel")
}

func (s *Server) Diagnostic(context.Context, *string) (*string, error) {
	return nil, notImplemented("Diagnostic")
}

func (s *Server) DiagnosticRefresh(context.Context) error {
	return notImplemented("DiagnosticRefresh")
}

func (s *Server) DiagnosticWorkspace(context.Context, *protocol.WorkspaceDiagnosticParams) (*protocol.WorkspaceDiagnosticReport, error) {
	return nil, notImplemented("DiagnosticWorkspace")
}

func (s *Server) DidChangeWatchedFiles(context.Context, *protocol.DidChangeWatchedFilesParams) error {
	return notImplemented("DidChangeWatchedFiles")
}

func (s *Server) DidChangeWorkspaceFolders(context.Context, *protocol.DidChangeWorkspaceFoldersParams) error {
	return notImplemented("DidChangeWorkspaceFolders")
}

func (s *Server) DidCreateFiles(context.Context, *protocol.CreateFilesParams) error {
	return notImplemented("DidCreateFiles")
}

func (s *Server) DidDeleteFiles(context.Context, *protocol.DeleteFilesParams) error {
	return notImplemented("DidDeleteFiles")
}

func (s *Server) DocumentLink(context.Context, *protocol.DocumentLinkParams) ([]protocol.DocumentLink, error) {
	return nil, nil
}

func notImplemented(method string) error {
	return fmt.Errorf("%w: %q not yet implemented", jsonrpc2.ErrMethodNotFound, method)
}

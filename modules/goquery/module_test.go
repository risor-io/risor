package goquery

import (
	"bytes"
	"context"
	"io"
	"strings"
	"testing"

	"github.com/risor-io/risor/object"
	"github.com/risor-io/risor/op"
	"github.com/stretchr/testify/require"
)

func TestModuleInitialization(t *testing.T) {
	module := Module()
	require.NotNil(t, module)
	require.Equal(t, object.MODULE, module.Type())

	// Check module has parse function
	parseObj, ok := module.GetAttr("parse")
	require.True(t, ok)
	require.NotNil(t, parseObj)

	// Verify parse is a builtin function
	parseBuiltin, ok := parseObj.(*object.Builtin)
	require.True(t, ok)
	require.Equal(t, "parse", parseBuiltin.Name())
}

func TestParseWithString(t *testing.T) {
	ctx := context.Background()
	html := `<html><body><div id="test">Hello World</div></body></html>`
	result := Parse(ctx, object.NewString(html))

	// Check result is a Document
	doc, ok := result.(*Document)
	require.True(t, ok)
	require.NotNil(t, doc)
	require.NotNil(t, doc.Value())

	// Check document content
	require.Contains(t, doc.String(), "Hello World")
}

func TestParseWithByteSlice(t *testing.T) {
	ctx := context.Background()
	html := []byte(`<html><body><div id="test">Hello World</div></body></html>`)
	result := Parse(ctx, object.NewByteSlice(html))

	// Check result is a Document
	doc, ok := result.(*Document)
	require.True(t, ok)
	require.NotNil(t, doc)
	require.NotNil(t, doc.Value())

	// Check document content
	require.Contains(t, doc.String(), "Hello World")
}

func TestParseWithFile(t *testing.T) {
	ctx := context.Background()
	html := `<html><body><div id="test">Hello World</div></body></html>`
	mockFile := &mockFile{
		reader: strings.NewReader(html),
	}
	result := Parse(ctx, mockFile)

	// Check result is a Document
	doc, ok := result.(*Document)
	require.True(t, ok)
	require.NotNil(t, doc)
	require.NotNil(t, doc.Value())

	// Check document content
	require.Contains(t, doc.String(), "Hello World")
}

func TestParseWithReader(t *testing.T) {
	ctx := context.Background()
	html := `<html><body><div id="test">Hello World</div></body></html>`
	reader := bytes.NewReader([]byte(html))
	result := Parse(ctx, &mockReader{reader: reader})

	// Check result is a Document
	doc, ok := result.(*Document)
	require.True(t, ok)
	require.NotNil(t, doc)
	require.NotNil(t, doc.Value())

	// Check document content
	require.Contains(t, doc.String(), "Hello World")
}

func TestParseWithInvalidType(t *testing.T) {
	ctx := context.Background()
	result := Parse(ctx, object.NewInt(123))

	// Check result is an error
	errObj, ok := result.(*object.Error)
	require.True(t, ok)
	require.Contains(t, errObj.Error(), "type error: expected reader")
}

func TestParseWithInvalidHTML(t *testing.T) {
	ctx := context.Background()
	// This is not truly invalid HTML, but let's try to force a parse error
	html := strings.Repeat("a", 1000000) // Very large input
	result := Parse(ctx, object.NewString(html))

	// Still should parse without error, just verify we get a document back
	doc, ok := result.(*Document)
	require.True(t, ok)
	require.NotNil(t, doc)
	require.NotNil(t, doc.Value())
}

// Test invalid number of arguments
func TestParseWithoutArguments(t *testing.T) {
	ctx := context.Background()
	result := Parse(ctx)

	// Check result is an error
	errObj, ok := result.(*object.Error)
	require.True(t, ok)
	require.Contains(t, errObj.Error(), "args error: goquery.parse() takes exactly 1 argument")
}

// Helper mock file for testing
type mockFile struct {
	reader io.Reader
	object.Object
}

func (f *mockFile) Read(p []byte) (n int, err error) {
	return f.reader.Read(p)
}

func (f *mockFile) Type() object.Type {
	return object.FILE
}

func (f *mockFile) Interface() interface{} {
	return f
}

func (f *mockFile) Inspect() string {
	return "mock_file"
}

func (f *mockFile) String() string {
	return f.Inspect()
}

func (f *mockFile) IsTruthy() bool {
	return true
}

func (f *mockFile) Equals(other object.Object) object.Object {
	if f == other.(object.Object) {
		return object.True
	}
	return object.False
}

func (f *mockFile) GetAttr(name string) (object.Object, bool) {
	return nil, false
}

func (f *mockFile) SetAttr(name string, value object.Object) error {
	return object.TypeErrorf("type error: cannot set attribute on mock_file")
}

func (f *mockFile) RunOperation(opType op.BinaryOpType, right object.Object) object.Object {
	return object.TypeErrorf("type error: unsupported operation for mock_file")
}

func (f *mockFile) Cost() int {
	return 0
}

// Helper mock reader for testing
type mockReader struct {
	reader *bytes.Reader
	object.Object
}

func (r *mockReader) Read(p []byte) (n int, err error) {
	return r.reader.Read(p)
}

func (r *mockReader) Type() object.Type {
	return "mock_reader"
}

func (r *mockReader) Interface() interface{} {
	return r
}

func (r *mockReader) Inspect() string {
	return "mock_reader"
}

func (r *mockReader) String() string {
	return r.Inspect()
}

func (r *mockReader) IsTruthy() bool {
	return true
}

func (r *mockReader) Equals(other object.Object) object.Object {
	if r == other.(object.Object) {
		return object.True
	}
	return object.False
}

func (r *mockReader) GetAttr(name string) (object.Object, bool) {
	return nil, false
}

func (r *mockReader) SetAttr(name string, value object.Object) error {
	return object.TypeErrorf("type error: cannot set attribute on mock_reader")
}

func (r *mockReader) RunOperation(opType op.BinaryOpType, right object.Object) object.Object {
	return object.TypeErrorf("type error: unsupported operation for mock_reader")
}

func (r *mockReader) Cost() int {
	return 0
}

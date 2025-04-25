package goquery

import (
	"context"
	"strings"
	"testing"

	"github.com/risor-io/risor/object"
	"github.com/risor-io/risor/op"
	"github.com/stretchr/testify/require"
)

func TestNewDocumentFromReader(t *testing.T) {
	html := `<html><body><div id="test">Hello World</div></body></html>`
	reader := strings.NewReader(html)

	doc, err := NewDocumentFromReader(reader)
	require.NoError(t, err)
	require.NotNil(t, doc)
	require.Equal(t, DOCUMENT, doc.Type())
	require.True(t, doc.IsTruthy())
}

func TestDocumentType(t *testing.T) {
	html := `<html><body><div id="test">Hello World</div></body></html>`
	reader := strings.NewReader(html)
	doc, err := NewDocumentFromReader(reader)
	require.NoError(t, err)

	require.Equal(t, DOCUMENT, doc.Type())
	require.Equal(t, "goquery.document()", doc.Inspect())
}

func TestDocumentString(t *testing.T) {
	html := `<html><body><div id="test">Hello World</div></body></html>`
	reader := strings.NewReader(html)
	doc, err := NewDocumentFromReader(reader)
	require.NoError(t, err)

	require.Contains(t, doc.String(), "Hello World")
}

func TestDocumentEquals(t *testing.T) {
	html := `<html><body><div id="test">Hello World</div></body></html>`
	reader := strings.NewReader(html)
	doc1, err := NewDocumentFromReader(reader)
	require.NoError(t, err)

	reader = strings.NewReader(html)
	doc2, err := NewDocumentFromReader(reader)
	require.NoError(t, err)

	// Different document instances should not be equal
	require.Equal(t, object.False, doc1.Equals(doc2))
	// Same document instance should be equal to itself
	require.Equal(t, object.True, doc1.Equals(doc1))
	// Different types should not be equal
	require.Equal(t, object.False, doc1.Equals(object.NewString("test")))
}

func TestDocumentSetAttr(t *testing.T) {
	html := `<html><body><div id="test">Hello World</div></body></html>`
	reader := strings.NewReader(html)
	doc, err := NewDocumentFromReader(reader)
	require.NoError(t, err)

	err = doc.SetAttr("test", object.NewString("value"))
	require.Error(t, err)
	require.Contains(t, err.Error(), "cannot set")
}

func TestDocumentRunOperation(t *testing.T) {
	html := `<html><body><div id="test">Hello World</div></body></html>`
	reader := strings.NewReader(html)
	doc, err := NewDocumentFromReader(reader)
	require.NoError(t, err)

	result := doc.RunOperation(op.Add, object.NewString("test"))
	_, ok := result.(*object.Error)
	require.True(t, ok)
}

func TestDocumentGetAttr(t *testing.T) {
	html := `<html><body><div id="test">Hello World</div></body></html>`
	reader := strings.NewReader(html)
	doc, err := NewDocumentFromReader(reader)
	require.NoError(t, err)

	// Test valid attributes
	find, ok := doc.GetAttr("find")
	require.True(t, ok)
	require.NotNil(t, find)

	html_, ok := doc.GetAttr("html")
	require.True(t, ok)
	require.NotNil(t, html_)

	text, ok := doc.GetAttr("text")
	require.True(t, ok)
	require.NotNil(t, text)

	// Test invalid attribute
	invalid, ok := doc.GetAttr("invalid")
	require.False(t, ok)
	require.Nil(t, invalid)
}

func TestDocumentFindMethod(t *testing.T) {
	html := `<html><body><div id="test">Hello World</div><div class="other">Other</div></body></html>`
	reader := strings.NewReader(html)
	doc, err := NewDocumentFromReader(reader)
	require.NoError(t, err)

	// Get the find method
	find, ok := doc.GetAttr("find")
	require.True(t, ok)
	builtin, ok := find.(*object.Builtin)
	require.True(t, ok)

	// Call find method with valid selector
	ctx := context.Background()
	result := builtin.Call(ctx, object.NewString("#test"))
	selection, ok := result.(*Selection)
	require.True(t, ok)
	require.Equal(t, 1, selection.Value().Length())

	// Call find method with invalid argument type
	result = builtin.Call(ctx, object.NewInt(123))
	_, ok = result.(*object.Error)
	require.True(t, ok)

	// Call find method with invalid number of arguments
	result = builtin.Call(ctx)
	_, ok = result.(*object.Error)
	require.True(t, ok)
}

func TestDocumentHTMLMethod(t *testing.T) {
	html := `<html><body><div id="test">Hello World</div></body></html>`
	reader := strings.NewReader(html)
	doc, err := NewDocumentFromReader(reader)
	require.NoError(t, err)

	// Get the html method
	htmlMethod, ok := doc.GetAttr("html")
	require.True(t, ok)
	builtin, ok := htmlMethod.(*object.Builtin)
	require.True(t, ok)

	// Call html method
	ctx := context.Background()
	result := builtin.Call(ctx)
	htmlStr, ok := result.(*object.String)
	require.True(t, ok)
	require.Contains(t, htmlStr.Value(), "Hello World")

	// Call html method with invalid number of arguments
	result = builtin.Call(ctx, object.NewString("extra"))
	_, ok = result.(*object.Error)
	require.True(t, ok)
}

func TestDocumentTextMethod(t *testing.T) {
	html := `<html><body><div id="test">Hello World</div></body></html>`
	reader := strings.NewReader(html)
	doc, err := NewDocumentFromReader(reader)
	require.NoError(t, err)

	// Get the text method
	textMethod, ok := doc.GetAttr("text")
	require.True(t, ok)
	builtin, ok := textMethod.(*object.Builtin)
	require.True(t, ok)

	// Call text method
	ctx := context.Background()
	result := builtin.Call(ctx)
	textStr, ok := result.(*object.String)
	require.True(t, ok)
	require.Contains(t, textStr.Value(), "Hello World")

	// Call text method with invalid number of arguments
	result = builtin.Call(ctx, object.NewString("extra"))
	_, ok = result.(*object.Error)
	require.True(t, ok)
}

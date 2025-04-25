package goquery

import (
	"context"
	"fmt"
	"strings"
	"testing"

	"github.com/risor-io/risor/object"
	"github.com/risor-io/risor/op"
	"github.com/stretchr/testify/require"
)

func createTestSelection(t *testing.T) *Selection {
	html := `<html><body>
		<div id="test">Hello World</div>
		<div class="items">
			<p class="item">Item 1</p>
			<p class="item">Item 2</p>
			<p class="item">Item 3</p>
		</div>
		<div id="nested">
			<span class="special">Special</span>
		</div>
		<a href="https://example.com">Link</a>
	</body></html>`
	reader := strings.NewReader(html)
	doc, err := NewDocumentFromReader(reader)
	require.NoError(t, err)

	selection := NewSelection(doc.Value().Find("body"))
	return selection
}

func TestSelectionType(t *testing.T) {
	sel := createTestSelection(t)
	require.Equal(t, SELECTION, sel.Type())
	require.Equal(t, "goquery.selection()", sel.Inspect())
}

func TestSelectionString(t *testing.T) {
	sel := createTestSelection(t)
	require.Contains(t, sel.String(), "Hello World")
}

func TestSelectionIsTruthy(t *testing.T) {
	sel := createTestSelection(t)
	require.True(t, sel.IsTruthy())

	// Test non-truthy selection (empty selection)
	emptySelection := NewSelection(sel.Value().Find("#nonexistent"))
	require.False(t, emptySelection.IsTruthy())
}

func TestSelectionEquals(t *testing.T) {
	sel1 := createTestSelection(t)
	// Get a copy of the same selection
	sel2 := NewSelection(sel1.Value())
	// Different selection objects but wrapping the same value
	sel3 := NewSelection(sel1.Value().Find("#test"))

	require.Equal(t, object.True, sel1.Equals(sel1))
	// sel1 and sel2 are different objects but point to the same selection
	require.Equal(t, object.True, sel1.Equals(sel2))
	// sel3 is a different selection
	require.Equal(t, object.False, sel1.Equals(sel3))
	// Different type
	require.Equal(t, object.False, sel1.Equals(object.NewString("test")))
}

func TestSelectionSetAttr(t *testing.T) {
	sel := createTestSelection(t)
	err := sel.SetAttr("test", object.NewString("value"))
	require.Error(t, err)
	require.Contains(t, err.Error(), "cannot set")
}

func TestSelectionRunOperation(t *testing.T) {
	sel := createTestSelection(t)
	result := sel.RunOperation(op.Add, object.NewString("test"))
	_, ok := result.(*object.Error)
	require.True(t, ok)
}

func TestSelectionGetAttr(t *testing.T) {
	sel := createTestSelection(t)

	// Test valid attributes
	methods := []string{
		"find", "attr", "html", "text", "each", "eq",
		"length", "first", "last", "parent", "children",
		"filter", "not", "has_class",
	}

	for _, methodName := range methods {
		method, ok := sel.GetAttr(methodName)
		require.True(t, ok, "Method %s should exist", methodName)
		require.NotNil(t, method)
		_, ok = method.(*object.Builtin)
		require.True(t, ok, "Method %s should be a builtin", methodName)
	}

	// Test invalid attribute
	invalid, ok := sel.GetAttr("invalid")
	require.False(t, ok)
	require.Nil(t, invalid)
}

func TestSelectionFindMethod(t *testing.T) {
	sel := createTestSelection(t)

	// Get the find method
	find, ok := sel.GetAttr("find")
	require.True(t, ok)
	builtin, ok := find.(*object.Builtin)
	require.True(t, ok)

	ctx := context.Background()

	// Find by ID
	result := builtin.Call(ctx, object.NewString("#test"))
	testSel, ok := result.(*Selection)
	require.True(t, ok)
	require.Equal(t, 1, testSel.Value().Length())

	// Find by class
	result = builtin.Call(ctx, object.NewString(".item"))
	itemsSel, ok := result.(*Selection)
	require.True(t, ok)
	require.Equal(t, 3, itemsSel.Value().Length())

	// Invalid argument type
	result = builtin.Call(ctx, object.NewInt(123))
	_, ok = result.(*object.Error)
	require.True(t, ok)

	// Missing argument
	result = builtin.Call(ctx)
	_, ok = result.(*object.Error)
	require.True(t, ok)
}

func TestSelectionAttrMethod(t *testing.T) {
	sel := createTestSelection(t)
	ctx := context.Background()

	// First find an element with an attribute
	findMethod, ok := sel.GetAttr("find")
	require.True(t, ok)
	findBuiltin, ok := findMethod.(*object.Builtin)
	require.True(t, ok)

	linkSel := findBuiltin.Call(ctx, object.NewString("a"))
	require.IsType(t, &Selection{}, linkSel)

	// Get the attr method from the link selection
	attrMethod, ok := linkSel.(*Selection).GetAttr("attr")
	require.True(t, ok)
	attrBuiltin, ok := attrMethod.(*object.Builtin)
	require.True(t, ok)

	// Get href attribute from the link
	result := attrBuiltin.Call(ctx, object.NewString("href"))
	href, ok := result.(*object.String)
	require.True(t, ok)
	require.Equal(t, "https://example.com", href.Value())

	// Get an attribute that doesn't exist
	result = attrBuiltin.Call(ctx, object.NewString("nonexistent"))
	require.Equal(t, object.Nil, result)

	// Invalid argument type
	result = attrBuiltin.Call(ctx, object.NewInt(123))
	_, ok = result.(*object.Error)
	require.True(t, ok)

	// Missing argument
	result = attrBuiltin.Call(ctx)
	_, ok = result.(*object.Error)
	require.True(t, ok)
}

func TestSelectionHTMLMethod(t *testing.T) {
	sel := createTestSelection(t)

	// Get the html method
	htmlMethod, ok := sel.GetAttr("html")
	require.True(t, ok)
	builtin, ok := htmlMethod.(*object.Builtin)
	require.True(t, ok)

	ctx := context.Background()

	// Get HTML of the selection
	result := builtin.Call(ctx)
	html, ok := result.(*object.String)
	require.True(t, ok)
	require.Contains(t, html.Value(), "Hello World")

	// With extra argument
	result = builtin.Call(ctx, object.NewString("extra"))
	_, ok = result.(*object.Error)
	require.True(t, ok)
}

func TestSelectionTextMethod(t *testing.T) {
	sel := createTestSelection(t)

	// Get the text method
	textMethod, ok := sel.GetAttr("text")
	require.True(t, ok)
	builtin, ok := textMethod.(*object.Builtin)
	require.True(t, ok)

	ctx := context.Background()

	// Get text of the selection
	result := builtin.Call(ctx)
	text, ok := result.(*object.String)
	require.True(t, ok)
	require.Contains(t, text.Value(), "Hello World")

	// With extra argument
	result = builtin.Call(ctx, object.NewString("extra"))
	_, ok = result.(*object.Error)
	require.True(t, ok)
}

func TestSelectionEachMethod(t *testing.T) {
	sel := createTestSelection(t)
	ctx := context.Background()

	// Find elements to iterate over
	findMethod, ok := sel.GetAttr("find")
	require.True(t, ok)
	findBuiltin, ok := findMethod.(*object.Builtin)
	require.True(t, ok)

	itemsSel := findBuiltin.Call(ctx, object.NewString(".item"))
	require.IsType(t, &Selection{}, itemsSel)

	// Get the each method from the items selection
	eachMethod, ok := itemsSel.(*Selection).GetAttr("each")
	require.True(t, ok)
	eachBuiltin, ok := eachMethod.(*object.Builtin)
	require.True(t, ok)

	// Create a test function using object.NewBuiltin
	count := 0
	mockFn := object.NewBuiltin("mockFn", func(ctx context.Context, args ...object.Object) object.Object {
		count++
		index, ok := args[0].(*object.Int)
		require.True(t, ok)
		require.Equal(t, int64(count-1), index.Value())

		// Second argument should be a selection
		itemSel, ok := args[1].(*Selection)
		require.True(t, ok)
		require.Equal(t, 1, itemSel.Value().Length())

		return object.Nil
	})

	// Call each
	result := eachBuiltin.Call(ctx, mockFn)
	fmt.Println(result)
	require.Equal(t, object.Nil, result)
	require.Equal(t, 3, count)

	// With invalid argument (not a function)
	result = eachBuiltin.Call(ctx, object.NewString("not a function"))
	_, ok = result.(*object.Error)
	require.True(t, ok)

	// Missing argument
	result = eachBuiltin.Call(ctx)
	_, ok = result.(*object.Error)
	require.True(t, ok)
}

func TestSelectionEqMethod(t *testing.T) {
	sel := createTestSelection(t)
	ctx := context.Background()

	// Find elements to work with
	findMethod, ok := sel.GetAttr("find")
	require.True(t, ok)
	findBuiltin, ok := findMethod.(*object.Builtin)
	require.True(t, ok)

	itemsSel := findBuiltin.Call(ctx, object.NewString(".item"))
	require.IsType(t, &Selection{}, itemsSel)

	// Get the eq method from the items selection
	eqMethod, ok := itemsSel.(*Selection).GetAttr("eq")
	require.True(t, ok)
	eqBuiltin, ok := eqMethod.(*object.Builtin)
	require.True(t, ok)

	// Get first item by index
	result := eqBuiltin.Call(ctx, object.NewInt(0))
	item0, ok := result.(*Selection)
	require.True(t, ok)
	require.Equal(t, 1, item0.Value().Length())

	// Get text of that item
	textMethod, ok := item0.GetAttr("text")
	require.True(t, ok)
	textBuiltin, ok := textMethod.(*object.Builtin)
	require.True(t, ok)

	textResult := textBuiltin.Call(ctx)
	text, ok := textResult.(*object.String)
	require.True(t, ok)
	require.Equal(t, "Item 1", text.Value())

	// With invalid argument type
	result = eqBuiltin.Call(ctx, object.NewString("not an int"))
	_, ok = result.(*object.Error)
	require.True(t, ok)

	// Missing argument
	result = eqBuiltin.Call(ctx)
	_, ok = result.(*object.Error)
	require.True(t, ok)
}

func TestSelectionLengthMethod(t *testing.T) {
	sel := createTestSelection(t)
	ctx := context.Background()

	// Find elements to work with
	findMethod, ok := sel.GetAttr("find")
	require.True(t, ok)
	findBuiltin, ok := findMethod.(*object.Builtin)
	require.True(t, ok)

	itemsSel := findBuiltin.Call(ctx, object.NewString(".item"))
	require.IsType(t, &Selection{}, itemsSel)

	// Get the length method from the items selection
	lengthMethod, ok := itemsSel.(*Selection).GetAttr("length")
	require.True(t, ok)
	lengthBuiltin, ok := lengthMethod.(*object.Builtin)
	require.True(t, ok)

	// Get length
	result := lengthBuiltin.Call(ctx)
	length, ok := result.(*object.Int)
	require.True(t, ok)
	require.Equal(t, int64(3), length.Value())

	// With extra argument
	result = lengthBuiltin.Call(ctx, object.NewString("extra"))
	_, ok = result.(*object.Error)
	require.True(t, ok)
}

func TestSelectionFirstLastMethod(t *testing.T) {
	sel := createTestSelection(t)
	ctx := context.Background()

	// Find elements to work with
	findMethod, ok := sel.GetAttr("find")
	require.True(t, ok)
	findBuiltin, ok := findMethod.(*object.Builtin)
	require.True(t, ok)

	itemsSel := findBuiltin.Call(ctx, object.NewString(".item"))
	require.IsType(t, &Selection{}, itemsSel)

	// Get first method from items selection
	firstMethod, ok := itemsSel.(*Selection).GetAttr("first")
	require.True(t, ok)
	firstBuiltin, ok := firstMethod.(*object.Builtin)
	require.True(t, ok)

	// Get last method from items selection
	lastMethod, ok := itemsSel.(*Selection).GetAttr("last")
	require.True(t, ok)
	lastBuiltin, ok := lastMethod.(*object.Builtin)
	require.True(t, ok)

	// Get text method
	textMethod, ok := sel.GetAttr("text")
	require.True(t, ok)
	textBuiltin, ok := textMethod.(*object.Builtin)
	require.True(t, ok)

	// Get first item
	firstResult := firstBuiltin.Call(ctx)
	firstItem, ok := firstResult.(*Selection)
	require.True(t, ok)

	firstTextResult := textBuiltin.Call(ctx, firstItem)
	firstText, ok := firstTextResult.(*object.String)
	require.True(t, ok)
	require.Equal(t, "Item 1", firstText.Value())

	// Get last item
	lastResult := lastBuiltin.Call(ctx)
	lastItem, ok := lastResult.(*Selection)
	require.True(t, ok)

	lastTextResult := textBuiltin.Call(ctx, lastItem)
	lastText, ok := lastTextResult.(*object.String)
	require.True(t, ok)
	require.Equal(t, "Item 3", lastText.Value())

	// With extra argument
	result := firstBuiltin.Call(ctx, object.NewString("extra"))
	_, ok = result.(*object.Error)
	require.True(t, ok)
}

func TestSelectionParentMethod(t *testing.T) {
	sel := createTestSelection(t)
	ctx := context.Background()

	// Find nested element
	findMethod, ok := sel.GetAttr("find")
	require.True(t, ok)
	findBuiltin, ok := findMethod.(*object.Builtin)
	require.True(t, ok)

	specialSel := findBuiltin.Call(ctx, object.NewString(".special"))
	require.IsType(t, &Selection{}, specialSel)

	// Get parent method from the special selection
	parentMethod, ok := specialSel.(*Selection).GetAttr("parent")
	require.True(t, ok)
	parentBuiltin, ok := parentMethod.(*object.Builtin)
	require.True(t, ok)

	// Get parent
	parentResult := parentBuiltin.Call(ctx)
	parent, ok := parentResult.(*Selection)
	require.True(t, ok)

	// Check parent has id 'nested'
	attrMethod, ok := parent.GetAttr("attr")
	require.True(t, ok)
	attrBuiltin, ok := attrMethod.(*object.Builtin)
	require.True(t, ok)

	idResult := attrBuiltin.Call(ctx, object.NewString("id"))
	id, ok := idResult.(*object.String)
	require.True(t, ok)
	require.Equal(t, "nested", id.Value())

	// With extra argument
	result := parentBuiltin.Call(ctx, object.NewString("extra"))
	_, ok = result.(*object.Error)
	require.True(t, ok)
}

func TestSelectionChildrenMethod(t *testing.T) {
	sel := createTestSelection(t)
	ctx := context.Background()

	// Find element with children
	findMethod, ok := sel.GetAttr("find")
	require.True(t, ok)
	findBuiltin, ok := findMethod.(*object.Builtin)
	require.True(t, ok)

	itemsContainer := findBuiltin.Call(ctx, object.NewString(".items"))
	require.IsType(t, &Selection{}, itemsContainer)

	// Get children method from the container
	childrenMethod, ok := itemsContainer.(*Selection).GetAttr("children")
	require.True(t, ok)
	childrenBuiltin, ok := childrenMethod.(*object.Builtin)
	require.True(t, ok)

	// Get all children
	childrenResult := childrenBuiltin.Call(ctx)
	children, ok := childrenResult.(*Selection)
	require.True(t, ok)

	// Should have 3 children
	require.Equal(t, 3, children.Value().Length())

	// Get children with filter
	childrenWithFilterResult := childrenBuiltin.Call(ctx, object.NewString(".item"))
	childrenWithFilter, ok := childrenWithFilterResult.(*Selection)
	require.True(t, ok)

	// Should still have 3 children
	require.Equal(t, 3, childrenWithFilter.Value().Length())

	// With invalid argument type
	result := childrenBuiltin.Call(ctx, object.NewInt(123))
	_, ok = result.(*object.Error)
	require.True(t, ok)

	// With too many arguments
	result = childrenBuiltin.Call(ctx, object.NewString(".item"), object.NewString("extra"))
	_, ok = result.(*object.Error)
	require.True(t, ok)
}

func TestSelectionFilterMethod(t *testing.T) {
	sel := createTestSelection(t)
	ctx := context.Background()

	// Find all divs
	findMethod, ok := sel.GetAttr("find")
	require.True(t, ok)
	findBuiltin, ok := findMethod.(*object.Builtin)
	require.True(t, ok)

	divs := findBuiltin.Call(ctx, object.NewString("div"))
	require.IsType(t, &Selection{}, divs)

	// Get filter method from divs
	filterMethod, ok := divs.(*Selection).GetAttr("filter")
	require.True(t, ok)
	filterBuiltin, ok := filterMethod.(*object.Builtin)
	require.True(t, ok)

	// Filter divs with id
	filteredResult := filterBuiltin.Call(ctx, object.NewString("[id]"))
	filtered, ok := filteredResult.(*Selection)
	require.True(t, ok)
	require.Equal(t, 2, filtered.Value().Length())

	// With invalid argument type
	result := filterBuiltin.Call(ctx, object.NewInt(123))
	_, ok = result.(*object.Error)
	require.True(t, ok)

	// Missing argument
	result = filterBuiltin.Call(ctx)
	_, ok = result.(*object.Error)
	require.True(t, ok)
}

func TestSelectionNotMethod(t *testing.T) {
	sel := createTestSelection(t)
	ctx := context.Background()

	// Find all divs
	findMethod, ok := sel.GetAttr("find")
	require.True(t, ok)
	findBuiltin, ok := findMethod.(*object.Builtin)
	require.True(t, ok)

	divs := findBuiltin.Call(ctx, object.NewString("div"))
	require.IsType(t, &Selection{}, divs)

	// Get not method from divs
	notMethod, ok := divs.(*Selection).GetAttr("not")
	require.True(t, ok)
	notBuiltin, ok := notMethod.(*object.Builtin)
	require.True(t, ok)

	// Get divs without id=test
	notResult := notBuiltin.Call(ctx, object.NewString("#test"))
	notDivs, ok := notResult.(*Selection)
	require.True(t, ok)

	// Should have 2 divs without id=test
	require.Equal(t, 2, notDivs.Value().Length())

	// With invalid argument type
	result := notBuiltin.Call(ctx, object.NewInt(123))
	_, ok = result.(*object.Error)
	require.True(t, ok)

	// Missing argument
	result = notBuiltin.Call(ctx)
	_, ok = result.(*object.Error)
	require.True(t, ok)
}

func TestSelectionHasClassMethod(t *testing.T) {
	sel := createTestSelection(t)
	ctx := context.Background()

	// Find special span
	findMethod, ok := sel.GetAttr("find")
	require.True(t, ok)
	findBuiltin, ok := findMethod.(*object.Builtin)
	require.True(t, ok)

	specialSpan := findBuiltin.Call(ctx, object.NewString(".special"))
	require.IsType(t, &Selection{}, specialSpan)

	// Get has_class method from special span
	hasClassMethod, ok := specialSpan.(*Selection).GetAttr("has_class")
	require.True(t, ok)
	hasClassBuiltin, ok := hasClassMethod.(*object.Builtin)
	require.True(t, ok)

	// Check it has special class
	hasClassResult := hasClassBuiltin.Call(ctx, object.NewString("special"))
	hasClass, ok := hasClassResult.(*object.Bool)
	require.True(t, ok)
	require.Equal(t, true, hasClass.Value())

	// Check it doesn't have another class
	hasClassResult = hasClassBuiltin.Call(ctx, object.NewString("other"))
	hasClass, ok = hasClassResult.(*object.Bool)
	require.True(t, ok)
	require.Equal(t, false, hasClass.Value())

	// With invalid argument type
	result := hasClassBuiltin.Call(ctx, object.NewInt(123))
	_, ok = result.(*object.Error)
	require.True(t, ok)

	// Missing argument
	result = hasClassBuiltin.Call(ctx)
	_, ok = result.(*object.Error)
	require.True(t, ok)
}

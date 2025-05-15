package htmltomarkdown

import (
	"context"
	"testing"

	"github.com/risor-io/risor/object"
	"github.com/stretchr/testify/assert"
)

func TestConvert(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		input    string
		expected string
	}{
		{
			input:    "<strong>Bold Text</strong>",
			expected: "**Bold Text**",
		},
		{
			input:    "<em>Italic Text</em>",
			expected: "*Italic Text*",
		},
		{
			input:    "<h1>Heading 1</h1>",
			expected: "# Heading 1",
		},
		{
			input:    "<a href=\"https://example.com\">Link</a>",
			expected: "[Link](https://example.com)",
		},
	}

	for _, tt := range tests {
		result := Convert(ctx, object.NewString(tt.input))
		str, ok := result.(*object.String)
		assert.True(t, ok)
		assert.Equal(t, tt.expected, str.Value())
	}
}

func TestConvertInvalidArgs(t *testing.T) {
	ctx := context.Background()

	// Test with wrong number of arguments
	result := Convert(ctx)
	_, ok := result.(*object.Error)
	assert.True(t, ok)

	// Test with non-string argument
	result = Convert(ctx, object.NewInt(123))
	_, ok = result.(*object.Error)
	assert.True(t, ok)
}

func TestCreateModule(t *testing.T) {
	module := Module()
	assert.NotNil(t, module)

	convert, ok := module.GetAttr("convert")
	assert.True(t, ok)
	assert.NotNil(t, convert)
	assert.IsType(t, &object.Builtin{}, convert)
}

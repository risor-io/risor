package parser_test

import (
	"testing"

	"github.com/cloudcmds/tamarin/parser"
	"github.com/stretchr/testify/require"
)

func TestParseString(t *testing.T) {
	tests := []struct {
		input string
		want  []*parser.StringTemplateFragment
	}{
		{
			"Hello ${name}!",
			[]*parser.StringTemplateFragment{
				{Value: "Hello ", IsVariable: false},
				{Value: "name", IsVariable: true},
				{Value: "!", IsVariable: false},
			},
		},
		{
			"abc Def {foo} $bar $} baz",
			[]*parser.StringTemplateFragment{
				{Value: "abc Def {foo} $bar $} baz", IsVariable: false},
			},
		},
		{
			"${ hi + 3 }${h[0]+foo.bar()}X{}${}",
			[]*parser.StringTemplateFragment{
				{Value: " hi + 3 ", IsVariable: true},
				{Value: "h[0]+foo.bar()", IsVariable: true},
				{Value: "X{}", IsVariable: false},
				{Value: "", IsVariable: true},
			},
		},
		{
			`\${1}`,
			[]*parser.StringTemplateFragment{
				{Value: "${1}", IsVariable: false},
			},
		},
	}
	for _, tc := range tests {
		res, err := parser.ParseStringTemplate(tc.input)
		require.Nil(t, err)
		require.Equal(t, tc.input, res.Value)
		require.Equal(t, tc.want, res.Fragments)
	}
}

func TestParseStringErrors(t *testing.T) {
	tests := []struct {
		input   string
		wantErr string
	}{
		{"${ ", `unterminated variable in template: "${ "`},
		{"a${0} ${cd", `unterminated variable in template: "a${0} ${cd"`},
		{"${bar + ${0}}", `invalid nesting in template: "${bar + ${0}}"`},
	}
	for _, tc := range tests {
		_, err := parser.ParseStringTemplate(tc.input)
		require.NotNil(t, err)
		require.Equal(t, tc.wantErr, err.Error())
	}
}

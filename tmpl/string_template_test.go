package tmpl_test

import (
	"testing"

	"github.com/cloudcmds/tamarin/tmpl"
	"github.com/stretchr/testify/require"
)

func TestParseString(t *testing.T) {
	tests := []struct {
		input string
		want  []*tmpl.Fragment
	}{
		{
			"Hello ${name}!",
			[]*tmpl.Fragment{
				{Value: "Hello ", IsVariable: false},
				{Value: "name", IsVariable: true},
				{Value: "!", IsVariable: false},
			},
		},
		{
			"abc Def {foo} $bar $} baz",
			[]*tmpl.Fragment{
				{Value: "abc Def {foo} $bar $} baz", IsVariable: false},
			},
		},
		{
			"${ hi + 3 }${h[0]+foo.bar()}X{}${}",
			[]*tmpl.Fragment{
				{Value: " hi + 3 ", IsVariable: true},
				{Value: "h[0]+foo.bar()", IsVariable: true},
				{Value: "X{}", IsVariable: false},
				{Value: "", IsVariable: true},
			},
		},
		{
			`\${1}`,
			[]*tmpl.Fragment{
				{Value: "${1}", IsVariable: false},
			},
		},
	}
	for _, tc := range tests {
		res, err := tmpl.Parse(tc.input)
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
		_, err := tmpl.Parse(tc.input)
		require.NotNil(t, err)
		require.Equal(t, tc.wantErr, err.Error())
	}
}

package tmpl_test

import (
	"testing"

	"github.com/cloudcmds/tamarin/core/tmpl"
	"github.com/stretchr/testify/require"
)

func TestParseString(t *testing.T) {
	tests := []struct {
		input string
		want  []*tmpl.Fragment
	}{
		{
			"Hello {name}!",
			[]*tmpl.Fragment{
				{Value: "Hello ", IsVariable: false},
				{Value: "name", IsVariable: true},
				{Value: "!", IsVariable: false},
			},
		},
		{
			"ab{{c}} {foo} $bar baz\t",
			[]*tmpl.Fragment{
				{Value: "ab{c} ", IsVariable: false},
				{Value: "foo", IsVariable: true},
				{Value: " $bar baz\t", IsVariable: false},
			},
		},
		{
			"{ hi + 3 }{h[0]+foo.bar()}X{}${}",
			[]*tmpl.Fragment{
				{Value: " hi + 3 ", IsVariable: true},
				{Value: "h[0]+foo.bar()", IsVariable: true},
				{Value: "X", IsVariable: false},
				{Value: "", IsVariable: true},
				{Value: "$", IsVariable: false},
				{Value: "", IsVariable: true},
			},
		},
		{
			`{{1}}`,
			[]*tmpl.Fragment{
				{Value: "{1}", IsVariable: false},
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
		{"{", `missing '}' in template: {`},
		{"a{0} {cd", `missing '}' in template: a{0} {cd`},
		{`{ x.update({"foo": 1}) }`, `invalid '{' in template: { x.update({"foo": 1}) }`},
		{"{a}}", `invalid '}' in template: {a}}`},
		{"}a", `invalid '}' in template: }a`},
	}
	for _, tc := range tests {
		_, err := tmpl.Parse(tc.input)
		require.NotNil(t, err)
		require.Equal(t, tc.wantErr, err.Error())
	}
}

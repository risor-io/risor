package rdoc

import (
	"fmt"
	"os"
	"testing"

	"github.com/russross/blackfriday/v2"
	"github.com/stretchr/testify/require"
)

func TestRegexpDocs(t *testing.T) {
	md, err := os.ReadFile("./fixtures/regexp.md")
	if err != nil {
		t.Fatal(err)
	}
	doc := Parse(string(md))

	require.Len(t, doc.Functions, 2)

	compileFn := doc.Functions[0]

	require.Equal(t, "regexp", compileFn.Module)

	require.Equal(t, "Module regexp provides regular expression matching.\n\nMore info here.", doc.Description)

	require.Equal(t, "compile", compileFn.Name)
	require.Equal(t, "compile(expr string) regexp", compileFn.Signature)
	require.Equal(t, "Compiles a regular expression string into a regexp object.", compileFn.Description)
	require.Len(t, compileFn.Examples, 1)

	example := compileFn.Examples[0]
	require.Equal(t, "Example", example.Name)
	require.Len(t, example.Statements, 2)

	s1 := example.Statements[0]
	require.Equal(t, `r := regexp.compile("a+"); r.match("a")`, s1.Code)
	require.Equal(t, `true`, s1.Result)

	s2 := example.Statements[1]
	require.Equal(t, `r := regexp.compile("[0-9]+"); r.match("nope")`, s2.Code)
	require.Equal(t, `false`, s2.Result)
}

func TestParseParagraph(t *testing.T) {

	// paragraph := "# regexp"
	// paragraph += "\n\n"
	// paragraph += "Module `regexp` provides regular expression matching."
	// paragraph += "\n\n"
	// paragraph += "More info here."

	md, err := os.ReadFile("./fixtures/regexp.md")
	if err != nil {
		t.Fatal(err)
	}

	// extensions := blackfriday.CommonExtensions
	// blackfriday.WithExtensions(extensions)

	node := blackfriday.New().Parse([]byte(md))

	current := node.FirstChild
	for {
		if current == nil {
			break
		}
		// fmt.Println(current)
		switch current.Type {
		case blackfriday.Heading:
			// fmt.Println("heading", getNodeText(current), current.Level)
		case blackfriday.Paragraph:
			if current.FirstChild != nil && current.FirstChild.Type == blackfriday.Code {
				fmt.Println("code paragraph", current.FirstChild.Type, getNodeText(current))
			} else {
				fmt.Println("norm paragraph", current.FirstChild.Type, getNodeText(current))
			}
		case blackfriday.Code:
			fmt.Println("code", getNodeText(current))
		case blackfriday.Text:
			fmt.Println("text", string(current.Literal), "child:", current.FirstChild.Literal)
		}
		current = current.Next
	}

	require.False(t, true)
}

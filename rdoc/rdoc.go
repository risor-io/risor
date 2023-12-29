package rdoc

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/russross/blackfriday/v2"
)

var filenameRegex = regexp.MustCompile(`filename="(.*)"`)

// Statement represents a single statement evaluation from an example.
type Statement struct {
	Code   string
	Result string
}

// Example of code execution parsed from markdown.
type Example struct {
	Name       string
	Statements []Statement
}

// Function description parsed from markdown.
type Function struct {
	Module      string
	Name        string
	Signature   string
	Description string
	Examples    []Example
}

type Module struct {
	Name        string
	Description string
	Functions   []*Function
}

func (d *Module) GetFunction(name string) (*Function, bool) {
	for _, f := range d.Functions {
		if f.Name == name {
			return f, true
		}
	}
	return nil, false
}

func ParseSignature(signature string) string {
	lines := strings.Split(signature, "\n")
	if len(lines) < 2 {
		return ""
	}
	return strings.TrimSpace(lines[1])
}

func ParseFilename(s string) string {
	// go copy filename="Example" -> Example
	matches := filenameRegex.FindStringSubmatch(s)
	if len(matches) < 2 {
		return ""
	}
	return matches[1]
}

func ParseExample(code string) Example {
	lines := strings.Split(code, "\n")
	if len(lines) < 2 {
		return Example{}
	}
	name := ParseFilename(lines[0])
	var statements []Statement
	for i := 2; i < len(lines)-1; i++ {
		line := lines[i]
		if strings.HasPrefix(line, ">>> ") {
			statements = append(statements, Statement{
				Code:   strings.TrimSpace(line[len(">>> "):]),
				Result: strings.TrimSpace(lines[i+1]),
			})
			i++
		}
	}
	return Example{Name: name, Statements: statements}
}

func getNodeText(node *blackfriday.Node) string {
	var text string
	node.Walk(func(node *blackfriday.Node, entering bool) blackfriday.WalkStatus {
		if entering {
			switch node.Type {
			case blackfriday.Text:
				text += string(node.Literal)
			case blackfriday.Code:
				text += fmt.Sprintf("`%s`", string(node.Literal))
			}
		}
		return blackfriday.GoToNext
	})
	return text
}

func Parse(text string) *Module {

	var functions []*Function
	var currentFunc *Function
	var moduleName string
	var moduleDesc []string
	inHeader := true

	node := blackfriday.New().Parse([]byte(text))

	node.Walk(func(node *blackfriday.Node, entering bool) blackfriday.WalkStatus {
		if node.Type == blackfriday.Heading {
			level := node.HeadingData.Level
			if entering {
				if level == 1 {
					moduleName = string(node.FirstChild.Literal)
				} else {
					inHeader = false
				}
				if level < 3 {
					currentFunc = nil
				} else if level == 3 { // Function name
					text := string(node.FirstChild.Literal)
					currentFunc = &Function{
						Module: moduleName,
						Name:   text,
					}
					functions = append(functions, currentFunc)
				}
			}
		}
		if node.Type == blackfriday.Code && entering {
			if currentFunc != nil {
				code := string(node.Literal)
				if strings.Contains(code, "Function signature") {
					currentFunc.Signature = ParseSignature(code)
				} else if strings.Contains(code, "Example") {
					currentFunc.Examples = append(currentFunc.Examples, ParseExample(code))
				}
			} else {
				fmt.Println("CODE:", string(node.Literal))
			}
		}
		if node.Type == blackfriday.Text && entering {
			parent := node.Parent
			if inHeader {
				if parent.Type == blackfriday.Paragraph {
					text := strings.TrimSpace(string(node.Literal))
					if text != "" {
						moduleDesc = append(moduleDesc, text)
					}
				}
			} else if currentFunc != nil {
				text := strings.TrimSpace(string(node.Literal))
				if parent.Type == blackfriday.Paragraph && text != "" {
					if currentFunc.Description != "" {
						currentFunc.Description += "\n\n"
					}
					currentFunc.Description += text
				}
			}
		}
		return blackfriday.GoToNext
	})
	return &Module{
		Name:        moduleName,
		Description: strings.Join(moduleDesc, "\n\n"),
		Functions:   functions,
	}
}

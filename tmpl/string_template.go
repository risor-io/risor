package tmpl

import "fmt"

type Fragment struct {
	Value      string
	IsVariable bool
}

type Template struct {
	Value     string
	Fragments []*Fragment
}

func Parse(s string) (*Template, error) {

	runes := []rune(s)

	getChar := func(index int) rune {
		if index < 0 || index >= len(runes) {
			return 0
		}
		return runes[index]
	}

	template := &Template{Value: s}

	// Iterate through all runes in the string to find ${variable}s. We build up
	// a list of string "fragments", which are either raw text or variables.
	var curFragment *Fragment
	for i := 0; i < len(runes); i++ {
		char := getChar(i)
		peekChar := getChar(i + 1)
		if char == '{' && peekChar == '{' {
			// Escaped { literal
			char = '{'
			i++
		} else if char == '}' {
			if curFragment != nil && curFragment.IsVariable {
				// Closed expression
				curFragment = nil
				continue
			}
			if peekChar == '}' {
				// Escaped } literal
				char = '}'
				i++
			} else {
				// Unescaped } literal is illegal
				return nil, fmt.Errorf("invalid '}' in template: %v", s)
			}
		} else if char == '{' {
			// Start of an expression. Error if we're already in an expression.
			if curFragment != nil && curFragment.IsVariable {
				return nil, fmt.Errorf("invalid '{' in template: %v", s)
			}
			curFragment = &Fragment{
				IsVariable: true,
				Value:      "",
			}
			template.Fragments = append(template.Fragments, curFragment)
			continue
		} else if curFragment != nil && curFragment.IsVariable && char == '}' {
			// End of an expression
			curFragment = nil
			continue
		}
		// Append current character to the current fragment
		if curFragment == nil {
			curFragment = &Fragment{
				IsVariable: false,
				Value:      "",
			}
			template.Fragments = append(template.Fragments, curFragment)
		}
		curFragment.Value += string(char)
	}

	if curFragment != nil && curFragment.IsVariable {
		return nil, fmt.Errorf("missing '}' in template: %v", s)
	}
	return template, nil
}

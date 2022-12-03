package parser

import "fmt"

type StringTemplateFragment struct {
	Value      string
	IsVariable bool
}

type StringTemplate struct {
	Value     string
	Fragments []*StringTemplateFragment
}

func ParseStringTemplate(s string) (*StringTemplate, error) {

	runes := []rune(s)

	getChar := func(index int) rune {
		if index < 0 || index >= len(runes) {
			return 0
		}
		return runes[index]
	}

	template := &StringTemplate{Value: s}

	// Iterate through all runes in the string to find ${variable}s. We build up
	// a list of string "fragments", which are either raw text or variables.
	var curFragment *StringTemplateFragment
	for i := 0; i < len(runes); i++ {
		char := getChar(i)
		prevChar := getChar(i - 1)
		nextChar := getChar(i + 1)
		if char == '\\' && nextChar == '$' {
			// Escaped $, so skip forward one and treat the $ as a literal
			char = '$'
			i++
		} else if char == '$' {
			if prevChar != '\\' && nextChar == '{' {
				if curFragment != nil && curFragment.IsVariable {
					return nil, fmt.Errorf("invalid nesting in template: \"%s\"", s)
				}
				curFragment = &StringTemplateFragment{
					IsVariable: true,
					Value:      "",
				}
				template.Fragments = append(template.Fragments, curFragment)
				i += 1 // Skip the following { character
				continue
			}
		} else if curFragment != nil && curFragment.IsVariable && char == '}' {
			curFragment = nil
			continue
		}
		if curFragment == nil {
			curFragment = &StringTemplateFragment{
				IsVariable: false,
				Value:      "",
			}
			template.Fragments = append(template.Fragments, curFragment)
		}
		curFragment.Value += string(char)
	}

	if curFragment != nil && curFragment.IsVariable {
		return nil, fmt.Errorf("unterminated variable in template: \"%s\"", s)
	}
	return template, nil
}

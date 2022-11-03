package evaluator

import (
	"fmt"
	"regexp"

	"github.com/myzie/tamarin/internal/object"
	"github.com/myzie/tamarin/internal/scope"
)

func matches(left, right object.Object, s *scope.Scope) object.Object {

	str := left.Inspect()

	if right.Type() != object.REGEXP_OBJ {
		return newError("regexp required for regexp-match, given %s", right.Type())
	}

	val := right.(*object.Regexp).Value
	if right.(*object.Regexp).Flags != "" {
		val = "(?" + right.(*object.Regexp).Flags + ")" + val
	}

	// Compile the regular expression.
	r, err := regexp.Compile(val)

	// Ensure it compiled
	if err != nil {
		return newError("error compiling regexp '%s': %s", right.Inspect(), err)
	}

	res := r.FindStringSubmatch(str)

	// Do we have any captures?
	if len(res) > 1 {
		for i := 1; i < len(res); i++ {
			s.Update(fmt.Sprintf("$%d", i), &object.String{Value: res[i]})
		}
	}

	// Test if it matched
	if len(res) > 0 {
		return object.TRUE
	}

	return object.FALSE
}

func notMatches(left, right object.Object) object.Object {
	str := left.Inspect()

	if right.Type() != object.REGEXP_OBJ {
		return newError("regexp required for regexp-match, given %s", right.Type())
	}

	val := right.(*object.Regexp).Value
	if right.(*object.Regexp).Flags != "" {
		val = "(?" + right.(*object.Regexp).Flags + ")" + val
	}

	// Compile the regular expression.
	r, err := regexp.Compile(val)

	// Ensure it compiled
	if err != nil {
		return newError("error compiling regexp '%s': %s", right.Inspect(), err)
	}

	// Test if it matched
	if r.MatchString(str) {
		return object.FALSE
	}

	return object.TRUE
}

package template

import (
	"fmt"
	"regexp"
	"strings"

	"k8s.io/client-go/util/jsonpath"
)

var jsonRegexp = regexp.MustCompile(`^\{\.?([^{}]+)\}$|^\.?([^{}]+)$`)

// taken from https://github.com/kubernetes/kubectl/blob/master/pkg/cmd/get/customcolumn.go
// relaxedJSONPathExpression attempts to be flexible with JSONPath expressions, it accepts:
//   - metadata.name (no leading '.' or curly braces '{...}'
//   - {metadata.name} (no leading '.')
//   - .metadata.name (no curly braces '{...}')
//   - {.metadata.name} (complete expression)
//
// And transforms them all into a valid jsonpath expression:
//
//	{.metadata.name}
func relaxedJSONPathExpression(pathExpression string) (string, error) {
	if len(pathExpression) == 0 {
		return pathExpression, nil
	}
	submatches := jsonRegexp.FindStringSubmatch(pathExpression)
	if submatches == nil {
		return "", fmt.Errorf("unexpected path string, expected a 'name1.name2' or '.name1.name2' or '{name1.name2}' or '{.name1.name2}'")
	}
	if len(submatches) != 3 {
		return "", fmt.Errorf("unexpected submatch list: %v", submatches)
	}
	var fieldSpec string
	if len(submatches[1]) != 0 {
		fieldSpec = submatches[1]
	} else {
		fieldSpec = submatches[2]
	}
	return fmt.Sprintf("{.%s}", fieldSpec), nil
}

func jsonPath(path string, obj any) (any, error) {
	key, err := relaxedJSONPathExpression(path)
	if err != nil {
		return nil, err
	}

	j := jsonpath.New("").AllowMissingKeys(false)

	if err := j.Parse(key); err != nil {
		return nil, fmt.Errorf("error parsing jsonpath: %w", err)
	}

	buf := new(strings.Builder)

	if err := j.Execute(buf, obj); err != nil {
		return nil, fmt.Errorf("error executing jsonpath: %w", err)
	}

	return strings.Trim(buf.String(), "\n"), nil
}

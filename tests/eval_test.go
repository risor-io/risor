package tests

import (
	"context"
	"errors"
	"os"
	"strings"
	"testing"

	"github.com/cloudcmds/tamarin/evaluator"
	"github.com/cloudcmds/tamarin/exec"
	"github.com/cloudcmds/tamarin/object"
	"github.com/stretchr/testify/require"
)

type TestCase struct {
	Name          string
	Text          string
	ExpectedValue string
	ExpectedType  string
}

func readFile(name string) string {
	data, err := os.ReadFile(name)
	if err != nil {
		panic(err)
	}
	return string(data)
}

func parseExpectedValue(filename, text string) (value string, typ string, err error) {
	lines := strings.Split(text, "\n")
	for _, line := range lines {
		if !strings.HasPrefix(line, "// ") {
			continue
		}
		if strings.HasPrefix(line, "// expected value:") {
			value = strings.SplitN(line, ":", 2)[1]
		} else if strings.HasPrefix(line, "// expected type:") {
			typ = strings.SplitN(line, ":", 2)[1]
		}
	}
	if value != "" {
		return strings.TrimSpace(value), strings.TrimSpace(typ), nil
	}
	return "", "", errors.New("test file did not define an expected result")
}

func getTestCase(name string) (TestCase, error) {
	input := readFile(name)
	expectedValue, expectedType, err := parseExpectedValue(name, input)
	if err != nil {
		return TestCase{}, err
	}
	return TestCase{
		Name:          name,
		Text:          input,
		ExpectedValue: expectedValue,
		ExpectedType:  expectedType,
	}, nil
}

func execute(ctx context.Context, input string) (object.Object, error) {
	result, err := exec.Execute(ctx, exec.Opts{
		Input:    string(input),
		Importer: &evaluator.SimpleImporter{},
	})
	if err != nil {
		return nil, err
	}
	return result, nil
}

func listTestFiles() []string {
	files, err := os.ReadDir(".")
	if err != nil {
		panic(err)
	}
	var testFiles []string
	for _, f := range files {
		if strings.HasSuffix(f.Name(), ".tm") {
			testFiles = append(testFiles, f.Name())
		}
	}
	return testFiles
}

func TestFiles(t *testing.T) {
	for _, name := range listTestFiles() {
		if !strings.HasSuffix(name, ".tm") {
			continue
		}
		t.Run(name, func(t *testing.T) {
			tc, err := getTestCase(name)
			require.Nil(t, err)

			ctx := context.Background()
			result, err := execute(ctx, tc.Text)
			require.Nil(t, err)

			expectedType := object.Type(tc.ExpectedType)

			require.Equal(t, tc.ExpectedValue, result.Inspect())
			require.Equal(t, expectedType, result.Type())
		})
	}
}

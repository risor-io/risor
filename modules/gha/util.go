package gha

import (
	"fmt"
	"io"
	"math/rand"
	"strings"

	"github.com/risor-io/risor/object"
	"github.com/risor-io/risor/os"
)

func printableValue(obj object.Object) any {
	if iface := obj.Interface(); iface != nil {
		return iface
	}
	switch obj := obj.(type) {
	case fmt.Stringer:
		return obj.String()
	default:
		return obj.Inspect()
	}
}

var workflowCommandReplacer = strings.NewReplacer(
	// Trick to get newlines included
	// https://github.com/actions/toolkit/issues/193#issuecomment-605394935
	"\r\n", "%0A",
	"\r", "%0A",
	"\n", "%0A",
	"::", "%3A%3A",
)

func sanitizeCommandValue(value string) string {
	return workflowCommandReplacer.Replace(strings.TrimSuffix(value, "\n"))
}

func asCommandProps(obj *object.Map) map[string]any {
	props := make(map[string]any)
	for _, key := range obj.StringKeys() {
		cmdKey, ok := commandPropKey(key)
		if !ok {
			continue
		}
		props[cmdKey] = printableValue(obj.Get(key))
	}
	return props
}

func commandPropKey(key string) (string, bool) {
	switch key {
	case "title", "file", "line":
		return key, true
	case "column":
		return "col", true
	case "end_line":
		return "endLine", true
	case "end_column":
		return "endColumn", true
	default:
		return "", false
	}
}

func stringifyCommandProps(props map[string]any) string {
	var sb strings.Builder
	for key, value := range props {
		if sb.Len() == 0 {
			sb.WriteByte(' ')
		} else {
			sb.WriteByte(',')
		}
		sb.WriteString(key)
		sb.WriteByte('=')
		sb.WriteString(sanitizeCommandValue(fmt.Sprint(value)))
	}
	return sb.String()
}

func runWorkflowCommand(w io.Writer, cmd string, value any, props map[string]any) object.Object {
	if _, err := fmt.Fprintf(w, "::%s%s::%s\n",
		cmd,
		stringifyCommandProps(props),
		sanitizeCommandValue(fmt.Sprint(value))); err != nil {
		return object.Errorf("io error: %v", err)
	}
	return object.Nil
}

func appendWorkflowFile(ros os.OS, path string, message string) object.Object {
	file, err := ros.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		return object.Errorf("io error: %v", err)
	}
	defer file.Close()

	if _, err := fmt.Fprintln(file, strings.TrimSuffix(message, "\n")); err != nil {
		return object.Errorf("io error: %v", err)
	}
	return object.Nil
}

func workflowFileKeyValue(key, value any) string {
	keyStr := fmt.Sprint(key)
	valueStr := fmt.Sprint(value)
	eof := workflowFileKeyValueEOF(keyStr, valueStr)
	return fmt.Sprintf("%s<<%s\n%s\n%s", keyStr, eof, valueStr, eof)
}

func workflowFileKeyValueEOF(key, value string) string {
	const charset = "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	var eof = "EOF"
	for strings.Contains(key, eof) || strings.Contains(value, eof) {
		eof += string(charset[rand.Intn(len(charset))])
	}
	return eof
}

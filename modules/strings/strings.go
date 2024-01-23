package strings

import (
	"strings"
)

//risor:generate

//risor:export
func contains(s, substr string) bool {
	return strings.Contains(s, substr)
}

//risor:export has_prefix
func hasPrefix(s, prefix string) bool {
	return strings.HasPrefix(s, prefix)
}

//risor:export has_prefix
func hasSuffix(s, suffix string) bool {
	return strings.HasSuffix(s, suffix)
}

//risor:export
func count(s, substr string) int {
	return strings.Count(s, substr)
}

//risor:export
func compare(a, b string) int {
	return strings.Compare(a, b)
}

//risor:export
func repeat(s string, count int) string {
	return strings.Repeat(s, count)
}

//risor:export
func join(list []string, sep string) string {
	return strings.Join(list, sep)
}

//risor:export
func split(s, sep string) []string {
	return strings.Split(s, sep)
}

//risor:export
func fields(s string) []string {
	return strings.Fields(s)
}

//risor:export
func index(s, substr string) int {
	return strings.Index(s, substr)
}

//risor:export last_index
func lastIndex(s, substr string) int {
	return strings.LastIndex(s, substr)
}

//risor:export replace_all
func replaceAll(s, old, new string) string {
	return strings.ReplaceAll(s, old, new)
}

//risor:export to_lower
func toLower(s string) string {
	return strings.ToLower(s)
}

//risor:export to_upper
func toUpper(s string) string {
	return strings.ToUpper(s)
}

//risor:export
func trim(s, cutset string) string {
	return strings.Trim(s, cutset)
}

//risor:export trim_prefix
func trimPrefix(s, prefix string) string {
	return strings.TrimPrefix(s, prefix)
}

//risor:export trim_suffix
func trimSuffix(s, prefix string) string {
	return strings.TrimSuffix(s, prefix)
}

//risor:export trim_space
func trimSpace(s string) string {
	return strings.TrimSpace(s)
}

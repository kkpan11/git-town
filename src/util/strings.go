package util

import "strings"

// Indent outputs the given string with the given level of indentation
// on each line. Each level of indentation is two spaces.
func Indent(message string) string {
	return "  " + strings.Replace(message, "\n", "\n  ", -1)
}

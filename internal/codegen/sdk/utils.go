package sdk

import (
	"strings"
	"unicode"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

func LowerTitle(s string) string {
	if s == "" {
		return s
	}

	a := []rune(s)
	a[0] = unicode.ToLower(a[0])
	return string(a)
}

func Title(s string) string {

	if s == "" {
		return s
	}

	// If the first character is a digit return s
	//
	// When a string starts with a digit cases.Title skips all the digits and title case
	// the first character it finds.
	if unicode.IsDigit(rune(s[0])) {
		return s
	}
	return cases.Title(language.English, cases.NoLower).String(s)
}

// Go string literals cannot contain backtick. If a string contains
// a backtick, replace it the following way:
//
// input:
// 	SELECT `group` FROM foo
//
// output:
// 	SELECT ` + "`" + `group` + "`" + ` FROM foo
//
// The escaped string must be rendered inside an existing string literal
//
// A string cannot be escaped twice
func EscapeBacktick(s string) string {
	return strings.Replace(s, "`", "`+\"`\"+`", -1)
}

func DoubleSlashComment(s string) string {
	return "// " + strings.ReplaceAll(s, "\n", "\n// ")
}

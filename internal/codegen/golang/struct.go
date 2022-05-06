package golang

import (
	"strings"
	"unicode"
	"unicode/utf8"

	"github.com/kyleconroy/sqlc/internal/plugin"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

type Struct struct {
	Table   plugin.Identifier
	Name    string
	Fields  []Field
	Comment string
}

// StructName constructs a valid camel case value from a snake case
func StructName(name, rename string) string {

	if rename != "" {
		return rename
	}

	out := ""

	for _, p := range strings.Split(name, "_") {
		if p == "id" {
			out += "ID"
		} else {
			out += cases.Title(language.English).String(p)
		}
	}

	// If a name has a digit as its first char, prepand an underscore to make it a valid Go name.
	r, _ := utf8.DecodeRuneInString(out)
	if unicode.IsDigit(r) {
		return "_" + out
	} else {
		return out
	}
}

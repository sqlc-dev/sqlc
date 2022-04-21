package golang

import (
	"strings"

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

	return out
}

package golang

import (
	"strings"
	"unicode"
	"unicode/utf8"

	"github.com/sqlc-dev/sqlc/internal/codegen/golang/options"
	"github.com/sqlc-dev/sqlc/internal/plugin"
)

type Struct struct {
	Table   *plugin.Identifier
	Name    string
	Fields  []Field
	Comment string
}

func StructName(name string, opts *options.Options) string {
	if rename := opts.Rename[name]; rename != "" {
		return rename
	}
	out := ""
	name = strings.Map(func(r rune) rune {
		if unicode.IsLetter(r) {
			return r
		}
		if unicode.IsDigit(r) {
			return r
		}
		return rune('_')
	}, name)

	for _, p := range strings.Split(name, "_") {
		if p == "id" {
			out += "ID"
		} else {
			out += strings.Title(p)
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

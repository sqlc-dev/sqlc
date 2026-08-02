package golang

import (
	"strings"
	"unicode"
	"unicode/utf8"

	"github.com/sqlc-dev/sqlc/internal/codegen/golang/opts"
	"github.com/sqlc-dev/sqlc/internal/plugin"
)

type Struct struct {
	Table   *plugin.Identifier
	Name    string
	Fields  []Field
	Comment string
	// IsModel is true for table structs that live in the models file. When
	// the models file is generated into a different Go package, references
	// to these types from query files must be qualified.
	IsModel bool
}

func StructName(name string, options *opts.Options) string {
	if rename := options.Rename[name]; rename != "" {
		return rename
	}
	var out strings.Builder
	name = strings.Map(func(r rune) rune {
		if unicode.IsLetter(r) {
			return r
		}
		if unicode.IsDigit(r) {
			return r
		}
		return rune('_')
	}, name)

	for p := range strings.SplitSeq(name, "_") {
		if _, found := options.InitialismsMap[p]; found {
			out.WriteString(strings.ToUpper(p))
		} else {
			out.WriteString(strings.Title(p))
		}
	}

	// If a name has a digit as its first char, prepand an underscore to make it a valid Go name.
	r, _ := utf8.DecodeRuneInString(out.String())
	if unicode.IsDigit(r) {
		return "_" + out.String()
	} else {
		return out.String()
	}
}

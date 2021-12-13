package golang

import (
	"strings"

	"github.com/kyleconroy/sqlc/internal/config"
	"github.com/kyleconroy/sqlc/internal/core"
)

type Struct struct {
	Table   core.FQN
	Name    string
	Fields  []Field
	Comment string
}

func StructName(name string, settings config.CombinedSettings) string {
	if rename := settings.Rename[name]; rename != "" {
		return rename
	}
	out := ""
	for _, p := range strings.Split(name, "_") {
		if p == "id" {
			out += "ID"
		} else {
			out += strings.Title(p)
		}
	}
	return out
}

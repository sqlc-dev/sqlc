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

func doID(name string, settings config.CombinedSettings) (out string) {
	for _, p := range strings.Split(name, "_") {
		if p == "id" {
			out += "ID"
		} else {
			out += strings.Title(p)
		}
	}
	return
}

func doRename(name string, settings config.CombinedSettings) string {
	if rename := settings.Rename[name]; rename != "" {
		return rename
	}
	return name
}

func JSONTagName(name string, settings config.CombinedSettings) string {
	if settings.Go.RenameJSONTags {
		name = doRename(name, settings)
	}
	return name
}

func StructName(name string, settings config.CombinedSettings) string {
	if !settings.Go.SkipStructRename {
		name = doRename(name, settings)
	}
	return doID(name, settings)
}

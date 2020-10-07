package golang

import (
	"log"
	"strings"

	"github.com/kyleconroy/sqlc/internal/compiler"
	"github.com/kyleconroy/sqlc/internal/config"
)

func sqliteType(r *compiler.Result, col *compiler.Column, settings config.CombinedSettings) string {
	dt := col.DataType
	notNull := col.NotNull || col.IsArray
	pointer := ""
	if !notNull && settings.Go.UsePointers {
		pointer = "*"
	}

	switch dt {

	case "integer":
		if notNull || pointer != "" {
			return pointer + "int32"
		}
		return "sql.NullInt32"

	case "any":
		return "interface{}"

	}

	switch {

	case strings.HasPrefix(dt, "varchar"):
		if notNull || pointer != "" {
			return pointer + "string"
		}
		return "sql.NullString"

	default:
		log.Printf("unknown SQLite type: %s\n", dt)
		return "interface{}"

	}
}

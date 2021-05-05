package golang

import (
	"log"
	"strings"

	"github.com/kyleconroy/sqlc/internal/compiler"
	"github.com/kyleconroy/sqlc/internal/config"
)

func sqliteType(r *compiler.Result, col *compiler.Column, settings config.CombinedSettings) string {
	dt := strings.ToLower(col.DataType)
	notNull := col.NotNull || col.IsArray

	switch dt {

	case "int", "integer", "tinyint", "smallint", "mediumint", "bigint", "unsignedbigint", "int2", "int8", "numeric", "decimal":
		if notNull {
			return "int64"
		}
		return "sql.NullInt64"

	case "blob":
		if notNull {
			return "[]uint8"
		}
		return "*[]uint8"

	case "real", "double", "doubleprecision", "float":
		if notNull {
			return "float64"
		}
		return "sql.NullFloat64"

	case "boolean":
		if notNull {
			return "bool"
		}
		return "sql.NullBool"

	case "date", "datetime", "timestamp":
		if notNull {
			return "time.Time"
		}
		return "sql.NullTime"

	case "any":
		return "interface{}"

	}

	switch {

	case strings.HasPrefix(dt, "character"),
		strings.HasPrefix(dt, "varchar"),
		strings.HasPrefix(dt, "varyingcharacter"),
		strings.HasPrefix(dt, "nchar"),
		strings.HasPrefix(dt, "nativecharacter"),
		strings.HasPrefix(dt, "nvarchar"),
		dt == "text",
		dt == "clob":
		if notNull {
			return "string"
		}
		return "sql.NullString"

	default:
		log.Printf("unknown SQLite type: %s\n", dt)
		return "interface{}"

	}
}

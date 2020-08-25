package golang

import (
	"log"

	"github.com/kyleconroy/sqlc/internal/compiler"
	"github.com/kyleconroy/sqlc/internal/config"
	"github.com/kyleconroy/sqlc/internal/debug"
	"github.com/kyleconroy/sqlc/internal/sql/catalog"
)

func mysqlType(r *compiler.Result, col *compiler.Column, settings config.CombinedSettings) string {
	columnType := col.DataType
	notNull := col.NotNull || col.IsArray

	switch columnType {

	case "varchar", "text", "char", "tinytext", "mediumtext", "longtext":
		if notNull {
			return "string"
		}
		return "sql.NullString"

	case "int", "integer", "smallint", "mediumint", "year":
		if notNull {
			return "int32"
		}
		return "sql.NullInt32"

	case "bigint":
		if notNull {
			return "int64"
		}
		return "sql.NullInt64"

	case "blob", "binary", "varbinary", "tinyblob", "mediumblob", "longblob":
		return "[]byte"

	case "double", "double precision", "real":
		if notNull {
			return "float64"
		}
		return "sql.NullFloat64"

	case "decimal", "dec", "fixed":
		if notNull {
			return "string"
		}
		return "sql.NullString"

	case "enum":
		// TODO: Proper Enum support
		return "string"

	case "date", "timestamp", "datetime", "time":
		if notNull {
			return "time.Time"
		}
		return "sql.NullTime"

	case "boolean", "bool", "tinyint":
		if notNull {
			return "bool"
		}
		return "sql.NullBool"

	case "json":
		return "json.RawMessage"

	case "any":
		return "interface{}"

	default:
		for _, schema := range r.Catalog.Schemas {
			for _, typ := range schema.Types {
				switch t := typ.(type) {
				case *catalog.Enum:
					if t.Name == columnType {
						if schema.Name == r.Catalog.DefaultSchema {
							return StructName(t.Name, settings)
						}
						return StructName(schema.Name+"_"+t.Name, settings)
					}
				}
			}
		}
		if debug.Active {
			log.Printf("Unknown MySQL type: %s\n", columnType)
		}
		return "interface{}"

	}
}

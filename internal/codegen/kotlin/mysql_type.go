package kotlin

import (
	"log"

	"github.com/kyleconroy/sqlc/internal/compiler"
	"github.com/kyleconroy/sqlc/internal/config"
	"github.com/kyleconroy/sqlc/internal/debug"
	"github.com/kyleconroy/sqlc/internal/sql/catalog"
)

func mysqlType(r *compiler.Result, col *compiler.Column, settings config.CombinedSettings) (string, bool) {
	columnType := col.DataType

	switch columnType {

	case "varchar", "text", "char", "tinytext", "mediumtext", "longtext":
		return "String", false

	case "int", "integer", "smallint", "mediumint", "year":
		return "Int", false

	case "bigint":
		return "Long", false

	case "blob", "binary", "varbinary", "tinyblob", "mediumblob", "longblob":
		return "String", false

	case "double", "double precision", "real":
		return "Double", false

	case "decimal", "dec", "fixed":
		return "String", false

	case "enum":
		// TODO: Proper Enum support
		return "String", false

	case "date", "datetime", "time":
		return "LocalDateTime", false

	case "timestamp":
		return "Instant", false

	case "boolean", "bool", "tinyint":
		return "Boolean", false

	case "json":
		return "String", false

	case "any":
		return "Any", false

	default:
		for _, schema := range r.Catalog.Schemas {
			for _, typ := range schema.Types {
				switch t := typ.(type) {
				case *catalog.Enum:
					if t.Name == columnType {
						if schema.Name == r.Catalog.DefaultSchema {
							return DataClassName(t.Name, settings), true
						}
						return DataClassName(schema.Name+"_"+t.Name, settings), true
					}
				}
			}
		}
		if debug.Active {
			log.Printf("Unknown MySQL type: %s\n", columnType)
		}
		return "Any", false

	}
}

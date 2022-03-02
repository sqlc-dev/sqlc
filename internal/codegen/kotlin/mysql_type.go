package kotlin

import (
	"log"

	"github.com/kyleconroy/sqlc/internal/codegen/sdk"
	"github.com/kyleconroy/sqlc/internal/debug"
	"github.com/kyleconroy/sqlc/internal/plugin"
)

func mysqlType(req *plugin.CodeGenRequest, col *plugin.Column) (string, bool) {
	columnType := sdk.DataType(col.Type)

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
		for _, schema := range req.Catalog.Schemas {
			for _, enum := range schema.Enums {
				if columnType == enum.Name {
					if schema.Name == req.Catalog.DefaultSchema {
						return dataClassName(enum.Name, req.Settings), true
					}
					return dataClassName(schema.Name+"_"+enum.Name, req.Settings), true
				}
			}
		}
		if debug.Active {
			log.Printf("Unknown MySQL type: %s\n", columnType)
		}
		return "Any", false

	}
}

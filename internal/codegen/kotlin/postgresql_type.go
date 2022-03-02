package kotlin

import (
	"log"

	"github.com/kyleconroy/sqlc/internal/codegen/sdk"
	"github.com/kyleconroy/sqlc/internal/plugin"
)

func postgresType(req *plugin.CodeGenRequest, col *plugin.Column) (string, bool) {
	columnType := sdk.DataType(col.Type)

	switch columnType {
	case "serial", "pg_catalog.serial4":
		return "Int", false

	case "bigserial", "pg_catalog.serial8":
		return "Long", false

	case "smallserial", "pg_catalog.serial2":
		return "Short", false

	case "integer", "int", "int4", "pg_catalog.int4":
		return "Int", false

	case "bigint", "pg_catalog.int8":
		return "Long", false

	case "smallint", "pg_catalog.int2":
		return "Short", false

	case "float", "double precision", "pg_catalog.float8":
		return "Double", false

	case "real", "pg_catalog.float4":
		return "Float", false

	case "pg_catalog.numeric":
		return "java.math.BigDecimal", false

	case "bool", "pg_catalog.bool":
		return "Boolean", false

	case "jsonb":
		// TODO: support json and byte types
		return "String", false

	case "bytea", "blob", "pg_catalog.bytea":
		return "String", false

	case "date":
		// Date and time mappings from https://jdbc.postgresql.org/documentation/head/java8-date-time.html
		return "LocalDate", false

	case "pg_catalog.time", "pg_catalog.timetz":
		return "LocalTime", false

	case "pg_catalog.timestamp":
		return "LocalDateTime", false

	case "pg_catalog.timestamptz", "timestamptz":
		// TODO
		return "OffsetDateTime", false

	case "text", "pg_catalog.varchar", "pg_catalog.bpchar", "string":
		return "String", false

	case "uuid":
		return "UUID", false

	case "inet":
		// TODO
		return "net.IP", false

	case "void":
		// TODO
		// A void value always returns NULL. Since there is no built-in NULL
		// value into the SQL package, we'll use sql.NullBool
		return "sql.NullBool", false

	case "any":
		// TODO
		return "Any", false

	default:
		for _, schema := range req.Catalog.Schemas {
			if schema.Name == "pg_catalog" {
				continue
			}
			for _, enum := range schema.Enums {
				if columnType == enum.Name {
					if schema.Name == req.Catalog.DefaultSchema {
						return dataClassName(enum.Name, req.Settings), true
					}
					return dataClassName(schema.Name+"_"+enum.Name, req.Settings), true
				}
			}
		}
		log.Printf("unknown PostgreSQL type: %s\n", columnType)
		return "Any", false
	}
}

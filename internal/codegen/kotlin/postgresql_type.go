package kotlin

import (
	"log"

	"github.com/kyleconroy/sqlc/internal/compiler"
	"github.com/kyleconroy/sqlc/internal/config"
	"github.com/kyleconroy/sqlc/internal/sql/catalog"
)

func postgresType(r *compiler.Result, col *compiler.Column, settings config.CombinedSettings) (string, bool) {
	columnType := col.DataType

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
		// TODO
		return "uuid.UUID", false

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
		for _, schema := range r.Catalog.Schemas {
			if schema.Name == "pg_catalog" {
				continue
			}
			for _, typ := range schema.Types {
				enum, ok := typ.(*catalog.Enum)
				if !ok {
					continue
				}
				if columnType == enum.Name {
					if schema.Name == r.Catalog.DefaultSchema {
						return DataClassName(enum.Name, settings), true
					}
					return DataClassName(schema.Name+"_"+enum.Name, settings), true
				}
			}
		}
		log.Printf("unknown PostgreSQL type: %s\n", columnType)
		return "Any", false
	}
}

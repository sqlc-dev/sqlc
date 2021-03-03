package python

import (
	"github.com/kyleconroy/sqlc/internal/compiler"
	"github.com/kyleconroy/sqlc/internal/config"
	"github.com/kyleconroy/sqlc/internal/sql/catalog"
	"log"
)

func postgresType(r *compiler.Result, col *compiler.Column, settings config.CombinedSettings) string {
	columnType := col.DataType

	switch columnType {
	case "serial", "serial4", "pg_catalog.serial4", "bigserial", "serial8", "pg_catalog.serial8", "smallserial", "serial2", "pg_catalog.serial2", "integer", "int", "int4", "pg_catalog.int4", "bigint", "int8", "pg_catalog.int8", "smallint", "int2", "pg_catalog.int2":
		return "int"
	case "float", "double precision", "float8", "pg_catalog.float8", "real", "float4", "pg_catalog.float4":
		return "float"
	case "numeric", "pg_catalog.numeric", "money":
		return "decimal.Decimal"
	case "boolean", "bool", "pg_catalog.bool":
		return "bool"
	case "json", "jsonb":
		return "Any"
	case "bytea", "blob", "pg_catalog.bytea":
		return "memoryview"
	case "date":
		return "datetime.date"
	case "pg_catalog.time", "pg_catalog.timetz":
		return "datetime.time"
	case "pg_catalog.timestamp", "pg_catalog.timestamptz", "timestamptz":
		return "datetime.datetime"
	case "interval", "pg_catalog.interval":
		return "datetime.timedelta"
	case "text", "pg_catalog.varchar", "pg_catalog.bpchar", "string", "citext":
		return "str"
	case "uuid":
		return "uuid.UUID"
	case "inet", "cidr", "macaddr", "macaddr8":
		// psycopg2 does have support for ipaddress objects, but it is not enabled by default
		//
		// https://www.psycopg.org/docs/extras.html#adapt-network
		return "str"
	case "ltree", "lquery", "ltxtquery":
		return "str"
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
						return "models." + ModelName(enum.Name, settings)
					}
					return "models." + ModelName(schema.Name+"_"+enum.Name, settings)
				}
			}
		}
		log.Printf("unknown PostgreSQL type: %s\n", columnType)
		return "Any"
	}
}

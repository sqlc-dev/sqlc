package golang

import (
	"log"

	"github.com/kyleconroy/sqlc/internal/compiler"
	"github.com/kyleconroy/sqlc/internal/config"
	"github.com/kyleconroy/sqlc/internal/debug"
	"github.com/kyleconroy/sqlc/internal/sql/catalog"
)

func postgresType(r *compiler.Result, col *compiler.Column, settings config.CombinedSettings) string {
	columnType := col.DataType
	notNull := col.NotNull || col.IsArray

	switch columnType {
	case "serial", "serial4", "pg_catalog.serial4":
		if notNull {
			return "int32"
		}
		return "sql.NullInt32"

	case "bigserial", "serial8", "pg_catalog.serial8":
		if notNull {
			return "int64"
		}
		return "sql.NullInt64"

	case "smallserial", "serial2", "pg_catalog.serial2":
		return "int16"

	case "integer", "int", "int4", "pg_catalog.int4":
		if notNull {
			return "int32"
		}
		return "sql.NullInt32"

	case "bigint", "int8", "pg_catalog.int8":
		if notNull {
			return "int64"
		}
		return "sql.NullInt64"

	case "smallint", "int2", "pg_catalog.int2":
		return "int16"

	case "float", "double precision", "float8", "pg_catalog.float8":
		if notNull {
			return "float64"
		}
		return "sql.NullFloat64"

	case "real", "float4", "pg_catalog.float4":
		if notNull {
			return "float32"
		}
		return "sql.NullFloat64" // TODO: Change to sql.NullFloat32 after updating the go.mod file

	case "numeric", "pg_catalog.numeric", "money":
		// Since the Go standard library does not have a decimal type, lib/pq
		// returns numerics as strings.
		//
		// https://github.com/lib/pq/issues/648
		if notNull {
			return "string"
		}
		return "sql.NullString"

	case "boolean", "bool", "pg_catalog.bool":
		if notNull {
			return "bool"
		}
		return "sql.NullBool"

	case "json", "jsonb":
		return "json.RawMessage"

	case "bytea", "blob", "pg_catalog.bytea":
		return "[]byte"

	case "date":
		if notNull {
			return "time.Time"
		}
		return "sql.NullTime"

	case "pg_catalog.time", "pg_catalog.timetz":
		if notNull {
			return "time.Time"
		}
		return "sql.NullTime"

	case "pg_catalog.timestamp", "pg_catalog.timestamptz", "timestamptz":
		if notNull {
			return "time.Time"
		}
		return "sql.NullTime"

	case "text", "pg_catalog.varchar", "pg_catalog.bpchar", "string":
		if notNull {
			return "string"
		}
		return "sql.NullString"

	case "uuid":
		return "uuid.UUID"

	case "inet", "cidr":
		return "net.IP"

	case "macaddr", "macaddr8":
		return "net.HardwareAddr"

	case "ltree", "lquery", "ltxtquery":
		// This module implements a data type ltree for representing labels
		// of data stored in a hierarchical tree-like structure. Extensive
		// facilities for searching through label trees are provided.
		//
		// https://www.postgresql.org/docs/current/ltree.html
		if notNull {
			return "string"
		}
		return "sql.NullString"

	case "interval", "pg_catalog.interval":
		if notNull {
			return "int64"
		}
		return "sql.NullInt64"

	case "void":
		// A void value always returns NULL. Since there is no built-in NULL
		// value into the SQL package, we'll use sql.NullBool
		return "sql.NullBool"

	case "any":
		return "interface{}"

	default:
		rel, err := compiler.ParseRelationString(columnType)
		if err != nil {
			// TODO: Should this actually return an error here?
			return "interface{}"
		}
		if rel.Schema == "" {
			rel.Schema = r.Catalog.DefaultSchema
		}

		for _, schema := range r.Catalog.Schemas {
			if schema.Name == "pg_catalog" {
				continue
			}
			for _, typ := range schema.Types {
				switch t := typ.(type) {
				case *catalog.Enum:
					if rel.Name == t.Name && rel.Schema == schema.Name {
						if schema.Name == r.Catalog.DefaultSchema {
							return StructName(t.Name, settings)
						}
						return StructName(schema.Name+"_"+t.Name, settings)
					}
				case *catalog.CompositeType:
					if notNull {
						return "string"
					}
					return "sql.NullString"
				}
			}
		}
		if debug.Active {
			log.Printf("unknown PostgreSQL type: %s\n", columnType)
		}
		return "interface{}"
	}
}

package golang

import (
	"fmt"
	"log"
	"strings"

	"github.com/kyleconroy/sqlc/internal/codegen/sdk"
	"github.com/kyleconroy/sqlc/internal/debug"
	"github.com/kyleconroy/sqlc/internal/plugin"
)

func parseIdentifierString(name string) (*plugin.Identifier, error) {
	parts := strings.Split(name, ".")
	switch len(parts) {
	case 1:
		return &plugin.Identifier{
			Name: parts[0],
		}, nil
	case 2:
		return &plugin.Identifier{
			Schema: parts[0],
			Name:   parts[1],
		}, nil
	case 3:
		return &plugin.Identifier{
			Catalog: parts[0],
			Schema:  parts[1],
			Name:    parts[2],
		}, nil
	default:
		return nil, fmt.Errorf("invalid name: %s", name)
	}
}

func postgresType(req *plugin.CodeGenRequest, col *plugin.Column) string {
	columnType := sdk.DataType(col.Type)
	notNull := col.NotNull || col.IsArray
	driver := parseDriver(req.Settings)

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
		if notNull {
			return "int16"
		}
		return "sql.NullInt16"

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
		if notNull {
			return "int16"
		}
		return "sql.NullInt16"

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
		if driver == SQLDriverPGXV4 {
			return "pgtype.Numeric"
		}
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

	case "json":
		switch driver {
		case SQLDriverPGXV4:
			return "pgtype.JSON"
		case SQLDriverLibPQ:
			if notNull {
				return "json.RawMessage"
			} else {
				return "pqtype.NullRawMessage"
			}
		default:
			return "interface{}"
		}

	case "jsonb":
		switch driver {
		case SQLDriverPGXV4:
			return "pgtype.JSONB"
		case SQLDriverLibPQ:
			if notNull {
				return "json.RawMessage"
			} else {
				return "pqtype.NullRawMessage"
			}
		default:
			return "interface{}"
		}

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
		if notNull {
			return "uuid.UUID"
		}
		return "uuid.NullUUID"

	case "inet":
		switch driver {
		case SQLDriverPGXV4:
			return "pgtype.Inet"
		case SQLDriverLibPQ:
			return "pqtype.Inet"
		default:
			return "interface{}"
		}

	case "cidr":
		switch driver {
		case SQLDriverPGXV4:
			return "pgtype.CIDR"
		case SQLDriverLibPQ:
			return "pqtype.CIDR"
		default:
			return "interface{}"
		}

	case "macaddr", "macaddr8":
		switch driver {
		case SQLDriverPGXV4:
			return "pgtype.Macaddr"
		case SQLDriverLibPQ:
			return "pqtype.Macaddr"
		default:
			return "interface{}"
		}

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

	case "daterange":
		if driver == SQLDriverPGXV4 {
			return "pgtype.Daterange"
		}
		return "interface{}"

	case "tsrange":
		if driver == SQLDriverPGXV4 {
			return "pgtype.Tsrange"
		}
		return "interface{}"

	case "tstzrange":
		if driver == SQLDriverPGXV4 {
			return "pgtype.Tstzrange"
		}
		return "interface{}"

	case "numrange":
		if driver == SQLDriverPGXV4 {
			return "pgtype.Numrange"
		}
		return "interface{}"

	case "int4range":
		if driver == SQLDriverPGXV4 {
			return "pgtype.Int4range"
		}
		return "interface{}"

	case "int8range":
		if driver == SQLDriverPGXV4 {
			return "pgtype.Int8range"
		}
		return "interface{}"

	case "hstore":
		if driver == SQLDriverPGXV4 {
			return "pgtype.Hstore"
		}
		return "interface{}"

	case "void":
		// A void value can only be scanned into an empty interface.
		return "interface{}"

	case "any":
		return "interface{}"

	default:
		rel, err := parseIdentifierString(columnType)
		if err != nil {
			// TODO: Should this actually return an error here?
			return "interface{}"
		}
		if rel.Schema == "" {
			rel.Schema = req.Catalog.DefaultSchema
		}

		for _, schema := range req.Catalog.Schemas {
			if schema.Name == "pg_catalog" {
				continue
			}

			for _, enum := range schema.Enums {
				if rel.Name == enum.Name && rel.Schema == schema.Name {
					if schema.Name == req.Catalog.DefaultSchema {
						return StructName(enum.Name, req.Settings)
					}
					return StructName(schema.Name+"_"+enum.Name, req.Settings)
				}
			}

			for _, ct := range schema.CompositeTypes {
				if rel.Name == ct.Name && rel.Schema == schema.Name {
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

package golang

import (
	"fmt"
	"log"
	"strings"

	"github.com/sqlc-dev/sqlc/internal/codegen/golang/opts"
	"github.com/sqlc-dev/sqlc/internal/codegen/sdk"
	"github.com/sqlc-dev/sqlc/internal/debug"
	"github.com/sqlc-dev/sqlc/internal/plugin"
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

func postgresType(req *plugin.GenerateRequest, options *opts.Options, col *plugin.Column) string {
	columnType := sdk.DataType(col.Type)
	notNull := col.NotNull || col.IsArray
	driver := parseDriver(options.SqlPackage)
	emitPointersForNull := driver.IsPGX() && options.EmitPointersForNullTypes

	switch columnType {
	case "serial", "serial4", "pg_catalog.serial4":
		if notNull {
			return "int32"
		}
		if emitPointersForNull {
			return "*int32"
		}
		if driver == opts.SQLDriverPGXV5 {
			return "pgtype.Int4"
		}
		return "sql.NullInt32"

	case "bigserial", "serial8", "pg_catalog.serial8":
		if notNull {
			return "int64"
		}
		if emitPointersForNull {
			return "*int64"
		}
		if driver == opts.SQLDriverPGXV5 {
			return "pgtype.Int8"
		}
		return "sql.NullInt64"

	case "smallserial", "serial2", "pg_catalog.serial2":
		if notNull {
			return "int16"
		}
		if emitPointersForNull {
			return "*int16"
		}
		if driver == opts.SQLDriverPGXV5 {
			return "pgtype.Int2"
		}
		return "sql.NullInt16"

	case "integer", "int", "int4", "pg_catalog.int4":
		if notNull {
			return "int32"
		}
		if emitPointersForNull {
			return "*int32"
		}
		if driver == opts.SQLDriverPGXV5 {
			return "pgtype.Int4"
		}
		return "sql.NullInt32"

	case "bigint", "int8", "pg_catalog.int8":
		if notNull {
			return "int64"
		}
		if emitPointersForNull {
			return "*int64"
		}
		if driver == opts.SQLDriverPGXV5 {
			return "pgtype.Int8"
		}
		return "sql.NullInt64"

	case "smallint", "int2", "pg_catalog.int2":
		if notNull {
			return "int16"
		}
		if emitPointersForNull {
			return "*int16"
		}
		if driver == opts.SQLDriverPGXV5 {
			return "pgtype.Int2"
		}
		return "sql.NullInt16"

	case "float", "double precision", "float8", "pg_catalog.float8":
		if notNull {
			return "float64"
		}
		if emitPointersForNull {
			return "*float64"
		}
		if driver == opts.SQLDriverPGXV5 {
			return "pgtype.Float8"
		}
		return "sql.NullFloat64"

	case "real", "float4", "pg_catalog.float4":
		if notNull {
			return "float32"
		}
		if emitPointersForNull {
			return "*float32"
		}
		if driver == opts.SQLDriverPGXV5 {
			return "pgtype.Float4"
		}
		return "sql.NullFloat64" // TODO: Change to sql.NullFloat32 after updating the go.mod file

	case "numeric", "pg_catalog.numeric", "money":
		if driver.IsPGX() {
			return "pgtype.Numeric"
		}
		// Since the Go standard library does not have a decimal type, lib/pq
		// returns numerics as strings.
		//
		// https://github.com/lib/pq/issues/648
		if notNull {
			return "string"
		}
		if emitPointersForNull {
			return "*string"
		}
		return "sql.NullString"

	case "boolean", "bool", "pg_catalog.bool":
		if notNull {
			return "bool"
		}
		if emitPointersForNull {
			return "*bool"
		}
		if driver == opts.SQLDriverPGXV5 {
			return "pgtype.Bool"
		}
		return "sql.NullBool"

	case "json", "pg_catalog.json":
		switch driver {
		case opts.SQLDriverPGXV5:
			return "[]byte"
		case opts.SQLDriverPGXV4:
			return "pgtype.JSON"
		case opts.SQLDriverLibPQ:
			if notNull {
				return "json.RawMessage"
			} else {
				return "pqtype.NullRawMessage"
			}
		default:
			return "interface{}"
		}

	case "jsonb", "pg_catalog.jsonb":
		switch driver {
		case opts.SQLDriverPGXV5:
			return "[]byte"
		case opts.SQLDriverPGXV4:
			return "pgtype.JSONB"
		case opts.SQLDriverLibPQ:
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
		if driver == opts.SQLDriverPGXV5 {
			return "pgtype.Date"
		}
		if notNull {
			return "time.Time"
		}
		if emitPointersForNull {
			return "*time.Time"
		}
		return "sql.NullTime"

	case "pg_catalog.time":
		if driver == opts.SQLDriverPGXV5 {
			return "pgtype.Time"
		}
		if notNull {
			return "time.Time"
		}
		if emitPointersForNull {
			return "*time.Time"
		}
		return "sql.NullTime"

	case "pg_catalog.timetz":
		if notNull {
			return "time.Time"
		}
		if emitPointersForNull {
			return "*time.Time"
		}
		return "sql.NullTime"

	case "pg_catalog.timestamp", "timestamp":
		if driver == opts.SQLDriverPGXV5 {
			return "pgtype.Timestamp"
		}
		if notNull {
			return "time.Time"
		}
		if emitPointersForNull {
			return "*time.Time"
		}
		return "sql.NullTime"

	case "pg_catalog.timestamptz", "timestamptz":
		if driver == opts.SQLDriverPGXV5 {
			return "pgtype.Timestamptz"
		}
		if notNull {
			return "time.Time"
		}
		if emitPointersForNull {
			return "*time.Time"
		}
		return "sql.NullTime"

	case "text", "pg_catalog.varchar", "pg_catalog.bpchar", "string", "citext", "name":
		if notNull {
			return "string"
		}
		if emitPointersForNull {
			return "*string"
		}
		if driver == opts.SQLDriverPGXV5 {
			return "pgtype.Text"
		}
		return "sql.NullString"

	case "uuid":
		if driver == opts.SQLDriverPGXV5 {
			return "pgtype.UUID"
		}
		if notNull {
			return "uuid.UUID"
		}
		if emitPointersForNull {
			return "*uuid.UUID"
		}
		return "uuid.NullUUID"

	case "inet":
		switch driver {
		case opts.SQLDriverPGXV5:
			if notNull {
				return "netip.Addr"
			}
			return "*netip.Addr"
		case opts.SQLDriverPGXV4:
			return "pgtype.Inet"
		case opts.SQLDriverLibPQ:
			return "pqtype.Inet"
		default:
			return "interface{}"
		}

	case "cidr":
		switch driver {
		case opts.SQLDriverPGXV5:
			if notNull {
				return "netip.Prefix"
			}
			return "*netip.Prefix"
		case opts.SQLDriverPGXV4:
			return "pgtype.CIDR"
		case opts.SQLDriverLibPQ:
			return "pqtype.CIDR"
		default:
			return "interface{}"
		}

	case "macaddr", "macaddr8":
		switch driver {
		case opts.SQLDriverPGXV5:
			return "net.HardwareAddr"
		case opts.SQLDriverPGXV4:
			return "pgtype.Macaddr"
		case opts.SQLDriverLibPQ:
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
		if emitPointersForNull {
			return "*string"
		}
		if driver == opts.SQLDriverPGXV5 {
			return "pgtype.Text"
		}
		return "sql.NullString"

	case "interval", "pg_catalog.interval":
		if driver == opts.SQLDriverPGXV5 {
			return "pgtype.Interval"
		}
		if notNull {
			return "int64"
		}
		if emitPointersForNull {
			return "*int64"
		}
		return "sql.NullInt64"

	case "daterange":
		switch driver {
		case opts.SQLDriverPGXV4:
			return "pgtype.Daterange"
		case opts.SQLDriverPGXV5:
			return "pgtype.Range[pgtype.Date]"
		default:
			return "interface{}"
		}

	case "datemultirange":
		switch driver {
		case opts.SQLDriverPGXV5:
			return "pgtype.Multirange[pgtype.Range[pgtype.Date]]"
		default:
			return "interface{}"
		}

	case "tsrange":
		switch driver {
		case opts.SQLDriverPGXV4:
			return "pgtype.Tsrange"
		case opts.SQLDriverPGXV5:
			return "pgtype.Range[pgtype.Timestamp]"
		default:
			return "interface{}"
		}

	case "tsmultirange":
		switch driver {
		case opts.SQLDriverPGXV5:
			return "pgtype.Multirange[pgtype.Range[pgtype.Timestamp]]"
		default:
			return "interface{}"
		}

	case "tstzrange":
		switch driver {
		case opts.SQLDriverPGXV4:
			return "pgtype.Tstzrange"
		case opts.SQLDriverPGXV5:
			return "pgtype.Range[pgtype.Timestamptz]"
		default:
			return "interface{}"
		}

	case "tstzmultirange":
		switch driver {
		case opts.SQLDriverPGXV5:
			return "pgtype.Multirange[pgtype.Range[pgtype.Timestamptz]]"
		default:
			return "interface{}"
		}

	case "numrange":
		switch driver {
		case opts.SQLDriverPGXV4:
			return "pgtype.Numrange"
		case opts.SQLDriverPGXV5:
			return "pgtype.Range[pgtype.Numeric]"
		default:
			return "interface{}"
		}

	case "nummultirange":
		switch driver {
		case opts.SQLDriverPGXV5:
			return "pgtype.Multirange[pgtype.Range[pgtype.Numeric]]"
		default:
			return "interface{}"
		}

	case "int4range":
		switch driver {
		case opts.SQLDriverPGXV4:
			return "pgtype.Int4range"
		case opts.SQLDriverPGXV5:
			return "pgtype.Range[pgtype.Int4]"
		default:
			return "interface{}"
		}

	case "int4multirange":
		switch driver {
		case opts.SQLDriverPGXV5:
			return "pgtype.Multirange[pgtype.Range[pgtype.Int4]]"
		default:
			return "interface{}"
		}

	case "int8range":
		switch driver {
		case opts.SQLDriverPGXV4:
			return "pgtype.Int8range"
		case opts.SQLDriverPGXV5:
			return "pgtype.Range[pgtype.Int8]"
		default:
			return "interface{}"
		}

	case "int8multirange":
		switch driver {
		case opts.SQLDriverPGXV5:
			return "pgtype.Multirange[pgtype.Range[pgtype.Int8]]"
		default:
			return "interface{}"
		}

	case "hstore":
		if driver.IsPGX() {
			return "pgtype.Hstore"
		}
		return "interface{}"

	case "bit", "varbit", "pg_catalog.bit", "pg_catalog.varbit":
		if driver == opts.SQLDriverPGXV5 {
			return "pgtype.Bits"
		}
		if driver == opts.SQLDriverPGXV4 {
			return "pgtype.Varbit"
		}

	case "cid":
		if driver == opts.SQLDriverPGXV5 {
			return "pgtype.Uint32"
		}
		if driver == opts.SQLDriverPGXV4 {
			return "pgtype.CID"
		}

	case "oid":
		if driver == opts.SQLDriverPGXV5 {
			return "pgtype.Uint32"
		}
		if driver == opts.SQLDriverPGXV4 {
			return "pgtype.OID"
		}

	case "tid":
		if driver.IsPGX() {
			return "pgtype.TID"
		}

	case "xid":
		if driver == opts.SQLDriverPGXV5 {
			return "pgtype.Uint32"
		}
		if driver == opts.SQLDriverPGXV4 {
			return "pgtype.XID"
		}

	case "box":
		if driver.IsPGX() {
			return "pgtype.Box"
		}

	case "circle":
		if driver.IsPGX() {
			return "pgtype.Circle"
		}

	case "line":
		if driver.IsPGX() {
			return "pgtype.Line"
		}

	case "lseg":
		if driver.IsPGX() {
			return "pgtype.Lseg"
		}

	case "path":
		if driver.IsPGX() {
			return "pgtype.Path"
		}

	case "point":
		if driver.IsPGX() {
			return "pgtype.Point"
		}

	case "polygon":
		if driver.IsPGX() {
			return "pgtype.Polygon"
		}

	case "vector":
		if driver == opts.SQLDriverPGXV5 {
			if emitPointersForNull {
				return "*pgvector.Vector"
			} else {
				return "pgvector.Vector"
			}
		}

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
			if schema.Name == "pg_catalog" || schema.Name == "information_schema" {
				continue
			}

			for _, enum := range schema.Enums {
				if rel.Name == enum.Name && rel.Schema == schema.Name {
					if notNull {
						if schema.Name == req.Catalog.DefaultSchema {
							return StructName(enum.Name, options)
						}
						return StructName(schema.Name+"_"+enum.Name, options)
					} else {
						if schema.Name == req.Catalog.DefaultSchema {
							return "Null" + StructName(enum.Name, options)
						}
						return "Null" + StructName(schema.Name+"_"+enum.Name, options)
					}
				}
			}

			for _, ct := range schema.CompositeTypes {
				if rel.Name == ct.Name && rel.Schema == schema.Name {
					if notNull {
						return "string"
					}
					if emitPointersForNull {
						return "*string"
					}
					return "sql.NullString"
				}
			}
		}
	}

	if debug.Active {
		log.Printf("unknown PostgreSQL type: %s\n", columnType)
	}
	return "interface{}"
}

package golang

import (
	"github.com/sqlc-dev/sqlc/internal/codegen/golang/opts"
	"github.com/sqlc-dev/sqlc/internal/plugin"
)

// duckdbType maps a DuckDB type to the Go type github.com/duckdb/duckdb-go
// hands database/sql for it. The driver returns a LIST or ARRAY as a slice
// of any, a STRUCT as a map by field name and a MAP as its own ordered
// map, and offers duckdb.Composite to decode a list into a typed slice, so
// a list column is a Composite of the slice its elements make. A Composite
// is not accepted as a query argument, while the driver binds a slice as a
// list and its own types as themselves, so a list parameter is the slice
// itself. DECIMAL and INTERVAL come
// back as the driver's structs and map to those; a UUID comes back as its
// bytes, which uuid.UUID scans. An ENUM is reported by the name it was
// created with, which the catalog does not yet carry, so it maps to any
// until it does.
func duckdbType(req *plugin.GenerateRequest, options *opts.Options, col *plugin.Column, param bool) string {
	t, nullable := columnType(col)
	if param {
		return duckdbGoType(options, t, nullable, duckdbParam)
	}
	return duckdbGoType(options, t, nullable, duckdbScan)
}

// duckdbPlace is where a value of the type goes: scanned from a row,
// passed as a query argument, or decoded inside a Composite, where the
// driver's own types are not decoded and a nullable value is a pointer.
type duckdbPlace int

const (
	duckdbScan duckdbPlace = iota
	duckdbParam
	duckdbNested
)

func duckdbGoType(options *opts.Options, t *plugin.TypeExpr, nullable bool, place duckdbPlace) string {
	null := func(typ, nullTyp string) string {
		if place == duckdbNested {
			nullTyp = ""
		}
		return nullableGoType(options, nullable, typ, nullTyp)
	}
	driver := func(typ, nullTyp string) string {
		if place == duckdbNested {
			return "any"
		}
		return nullableGoType(options, nullable, typ, nullTyp)
	}

	switch typeExprName(t) {
	case "array", "list":
		elem := typeExprArg(t, 0)
		var typ string
		if elem == nil {
			typ = "[]any"
		} else {
			typ = "[]" + duckdbGoType(options, elem, elem.Nullable, duckdbNested)
		}
		if place == duckdbScan {
			return "duckdb.Composite[" + typ + "]"
		}
		return typ

	case "struct":
		return null("map[string]any", "")

	case "map":
		return driver("duckdb.Map", "")

	case "boolean", "bool", "logical":
		return null("bool", "sql.NullBool")

	case "tinyint", "int1":
		return null("int8", "sql.NullInt16")

	case "smallint", "int2", "short":
		return null("int16", "sql.NullInt16")

	case "integer", "int", "int4", "signed":
		return null("int32", "sql.NullInt32")

	case "bigint", "int8", "long":
		return null("int64", "sql.NullInt64")

	case "utinyint":
		return null("uint8", "sql.NullByte")

	case "usmallint":
		return null("uint16", "")

	case "uinteger":
		return null("uint32", "")

	case "ubigint":
		return null("uint64", "")

	case "hugeint", "uhugeint", "bignum":
		return "*big.Int"

	case "float", "float4", "real":
		return null("float32", "sql.NullFloat64")

	case "double", "float8":
		return null("float64", "sql.NullFloat64")

	case "decimal", "numeric":
		return driver("duckdb.Decimal", "")

	case "varchar", "char", "bpchar", "text", "string", "json", "bit", "bitstring":
		return null("string", "sql.NullString")

	case "blob", "bytea", "binary", "varbinary":
		return "[]byte"

	case "date",
		"timestamp", "datetime", "timestamp_s", "timestamp_ms", "timestamp_ns",
		"timestamp with time zone", "timestamptz", "timestamptz_ns",
		"time", "time with time zone", "timetz", "time_ns":
		return null("time.Time", "sql.NullTime")

	case "interval":
		return driver("duckdb.Interval", "")

	case "uuid":
		return driver("uuid.UUID", "uuid.NullUUID")

	default:
		return "any"
	}
}

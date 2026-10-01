package golang

import (
	"github.com/sqlc-dev/sqlc/internal/codegen/golang/opts"
	"github.com/sqlc-dev/sqlc/internal/plugin"
)

// clickhouseType maps a ClickHouse type to the Go type
// github.com/ClickHouse/clickhouse-go/v2 hands database/sql for it. The
// driver scans a Nullable value into a pointer, and the types with a
// database/sql wrapper take that instead unless the options ask for
// pointers. Wide integers are big.Int, decimals are shopspring decimals,
// and geo types are orb geometries, as the driver returns them. The
// nesting the type carries is rendered: Map(String, Array(UInt8)) is a
// map[string][]uint8, LowCardinality and SimpleAggregateFunction wrap the
// type they carry. A Tuple is a slice of any, or a map by field name when
// its elements are named.
//
// A parameter is typed for what the driver can send rather than what it
// scans. The std driver passes arguments through without calling their
// Value method, and a {name:Type} placeholder is sent as text formatted
// from the value itself, so a sql.Null wrapper would go out as its struct
// fields and a decimal, UUID or address quoted, which the server rejects.
// A nullable parameter is therefore a pointer, which the driver sends as
// NULL when nil, and the types without a plain Go form are strings, which
// the server parses as the placeholder's type.
func clickhouseType(req *plugin.GenerateRequest, options *opts.Options, col *plugin.Column, param bool) string {
	t, nullable := columnType(col)
	return clickhouseGoType(options, t, nullable, false, param)
}

// clickhouseGoType renders t. Inside an array or a map, where nested is
// set, the driver returns a nullable value as a pointer whatever the
// column's own nullability is scanned into, so the wrappers are not used.
// A parameter, where param is set, takes the pointer form too, and is
// typed at its top level for the text the driver sends.
func clickhouseGoType(options *opts.Options, t *plugin.TypeExpr, nullable, nested, param bool) string {
	param = param && !nested
	null := func(typ, nullTyp string) string {
		if nested || param {
			nullTyp = ""
		}
		return nullableGoType(options, nullable, typ, nullTyp)
	}
	// text is the type of a value the driver would send quoted: a string
	// for a parameter, typ otherwise.
	text := func(typ, nullTyp string) string {
		if param {
			return null("string", "")
		}
		return null(typ, nullTyp)
	}

	switch typeExprName(t) {
	case "nullable":
		if inner := typeExprArg(t, 0); inner != nil {
			return clickhouseGoType(options, inner, true, nested, param)
		}
		return "any"

	// The driver hands database/sql a Nullable value as a pointer, and
	// dereferences it only when Nullable is the column's outermost type.
	// Wrapped in LowCardinality or SimpleAggregateFunction the pointer
	// comes through as is, which a sql.Null wrapper cannot scan, so a
	// nullable one takes the pointer form the nested types do.
	case "lowcardinality":
		if inner := typeExprArg(t, 0); inner != nil {
			wrapped := nullable || t.Nullable || inner.Nullable
			return clickhouseGoType(options, inner, wrapped, nested || wrapped, param)
		}
		return "any"

	case "simpleaggregatefunction":
		if inner := typeExprLastArg(t); inner != nil {
			wrapped := nullable || t.Nullable || inner.Nullable
			return clickhouseGoType(options, inner, wrapped, nested || wrapped, param)
		}
		return "any"

	case "array":
		elem := typeExprArg(t, 0)
		if elem == nil {
			return "[]any"
		}
		return "[]" + clickhouseGoType(options, elem, elem.Nullable, true, false)

	case "map":
		key, value := typeExprArg(t, 0), typeExprArg(t, 1)
		if key == nil || value == nil {
			return "map[string]any"
		}
		return "map[" + clickhouseGoType(options, key, false, true, false) + "]" + clickhouseGoType(options, value, value.Nullable, true, false)

	case "tuple":
		for _, arg := range t.GetArgs() {
			if arg.GetLabel() != "" {
				return "map[string]any"
			}
		}
		return "[]any"

	case "bool":
		return null("bool", "sql.NullBool")

	case "int8":
		return null("int8", "sql.NullInt16")

	case "int16":
		return null("int16", "sql.NullInt16")

	case "int32":
		return null("int32", "sql.NullInt32")

	case "int64":
		return null("int64", "sql.NullInt64")

	case "uint8":
		return null("uint8", "sql.NullByte")

	case "uint16":
		return null("uint16", "")

	case "uint32":
		return null("uint32", "")

	case "uint64":
		return null("uint64", "")

	case "int128", "int256", "uint128", "uint256":
		if param {
			return null("string", "")
		}
		return "*big.Int"

	case "float32", "bfloat16":
		return null("float32", "sql.NullFloat64")

	case "float64":
		return null("float64", "sql.NullFloat64")

	case "decimal", "decimal32", "decimal64", "decimal128", "decimal256":
		return text("decimal.Decimal", "decimal.NullDecimal")

	case "string", "fixedstring", "enum", "enum8", "enum16":
		return null("string", "sql.NullString")

	case "date", "date32", "datetime", "datetime32", "datetime64":
		return null("time.Time", "sql.NullTime")

	case "time", "time64":
		return text("time.Duration", "")

	case "uuid":
		return text("uuid.UUID", "uuid.NullUUID")

	case "ipv4", "ipv6":
		return text("net.IP", "")

	case "point":
		return null("orb.Point", "")

	case "ring":
		return null("orb.Ring", "")

	case "linestring":
		return null("orb.LineString", "")

	case "multilinestring":
		return null("orb.MultiLineString", "")

	case "polygon":
		return null("orb.Polygon", "")

	case "multipolygon":
		return null("orb.MultiPolygon", "")

	default:
		return "any"
	}
}

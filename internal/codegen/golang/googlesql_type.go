package golang

import (
	"github.com/sqlc-dev/sqlc/internal/codegen/golang/opts"
	"github.com/sqlc-dev/sqlc/internal/plugin"
)

// googlesqlType maps a Spanner GoogleSQL type to the Go type
// github.com/googleapis/go-sql-spanner hands database/sql for it. The
// driver decodes the elements of an ARRAY into the Null-prefixed types of
// cloud.google.com/go/spanner unless it is told to decode them natively,
// so an array column is a slice of those. NUMERIC and JSON have no
// database/sql wrapper and take the spanner package's; a DATE is a civil
// date.
func googlesqlType(req *plugin.GenerateRequest, options *opts.Options, col *plugin.Column) string {
	t, nullable := columnType(col)
	if typeExprName(t) == "array" {
		elem := typeExprArg(t, 0)
		if elem == nil {
			return "[]any"
		}
		return "[]" + googlesqlElementType(elem)
	}
	null := func(typ, nullTyp string) string {
		return nullableGoType(options, nullable, typ, nullTyp)
	}

	switch typeExprName(t) {
	case "int64":
		return null("int64", "sql.NullInt64")

	case "float64":
		return null("float64", "sql.NullFloat64")

	case "float32":
		return null("float32", "spanner.NullFloat32")

	case "bool":
		return null("bool", "sql.NullBool")

	case "string":
		return null("string", "sql.NullString")

	case "bytes":
		return "[]byte"

	case "date":
		return null("civil.Date", "spanner.NullDate")

	case "timestamp":
		return null("time.Time", "sql.NullTime")

	case "numeric", "bignumeric":
		return null("big.Rat", "spanner.NullNumeric")

	case "json":
		return "spanner.NullJSON"

	case "uuid":
		return null("uuid.UUID", "spanner.NullUUID")

	case "interval":
		return null("spanner.Interval", "spanner.NullInterval")

	default:
		return "any"
	}
}

// googlesqlElementType is the type of an array's elements, as the driver
// decodes them: nullable whatever the schema says, in the spanner package's
// wrappers.
func googlesqlElementType(t *plugin.TypeExpr) string {
	switch typeExprName(t) {
	case "int64":
		return "spanner.NullInt64"
	case "float64":
		return "spanner.NullFloat64"
	case "float32":
		return "spanner.NullFloat32"
	case "bool":
		return "spanner.NullBool"
	case "string":
		return "spanner.NullString"
	case "bytes":
		return "[]byte"
	case "date":
		return "spanner.NullDate"
	case "timestamp":
		return "spanner.NullTime"
	case "numeric", "bignumeric":
		return "spanner.NullNumeric"
	case "json":
		return "spanner.NullJSON"
	case "uuid":
		return "spanner.NullUUID"
	case "interval":
		return "spanner.NullInterval"
	default:
		return "any"
	}
}

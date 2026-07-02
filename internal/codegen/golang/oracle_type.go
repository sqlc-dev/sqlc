package golang

import (
	"log"
	"strings"

	"github.com/sqlc-dev/sqlc/internal/codegen/golang/opts"
	"github.com/sqlc-dev/sqlc/internal/codegen/sdk"
	"github.com/sqlc-dev/sqlc/internal/debug"
	"github.com/sqlc-dev/sqlc/internal/plugin"
)

// oracleType maps an Oracle column type to a Go type. It targets the
// github.com/sijms/go-ora/v2 driver, which implements database/sql, so the
// nullable variants use the standard sql.Null* wrappers (or pointers when
// EmitPointersForNullTypes is set), mirroring the SQLite engine's mapping.
//
// Oracle type semantics (per go-ora):
//   - NUMBER with scale 0            -> integer  (int64)
//   - NUMBER with scale > 0          -> float64
//     sqlc does not track scale in the built-in analyzer, so bare NUMBER maps
//     to float64 to avoid silently truncating fractional values; INTEGER/INT and
//     the PL/SQL integer types map to int64.
//   - VARCHAR2/CHAR/CLOB/etc.        -> string
//   - DATE/TIMESTAMP*                -> time.Time
//   - RAW/BLOB/LONG RAW/BFILE        -> []byte
//   - BINARY_FLOAT/BINARY_DOUBLE     -> float64
func oracleType(req *plugin.GenerateRequest, options *opts.Options, col *plugin.Column) string {
	dt := strings.ToLower(sdk.DataType(col.Type))
	notNull := col.NotNull || col.IsArray
	emitPointersForNull := options.EmitPointersForNullTypes

	nullableString := func() string {
		if notNull {
			return "string"
		}
		if emitPointersForNull {
			return "*string"
		}
		return "sql.NullString"
	}
	nullableInt64 := func() string {
		if notNull {
			return "int64"
		}
		if emitPointersForNull {
			return "*int64"
		}
		return "sql.NullInt64"
	}
	nullableFloat64 := func() string {
		if notNull {
			return "float64"
		}
		if emitPointersForNull {
			return "*float64"
		}
		return "sql.NullFloat64"
	}
	nullableTime := func() string {
		if notNull {
			return "time.Time"
		}
		if emitPointersForNull {
			return "*time.Time"
		}
		return "sql.NullTime"
	}
	nullableBool := func() string {
		if notNull {
			return "bool"
		}
		if emitPointersForNull {
			return "*bool"
		}
		return "sql.NullBool"
	}

	switch dt {

	// Integer types.
	case "integer", "int", "smallint",
		"binary_integer", "pls_integer", "simple_integer",
		"natural", "naturaln", "positive", "positiven", "signtype":
		return nullableInt64()

	// Fixed / floating point numeric types.
	case "number", "numeric", "dec", "decimal",
		"float", "real", "double precision", "double",
		"binary_float", "binary_double":
		return nullableFloat64()

	// Boolean (Oracle 23c native BOOLEAN / PL/SQL BOOLEAN).
	case "boolean", "bool":
		return nullableBool()

	// Date / time types.
	case "date",
		"timestamp",
		"timestamp_unconstrained",
		"timestamp_tz_unconstrained",
		"timestamp_ltz_unconstrained":
		return nullableTime()

	// Binary types.
	case "raw", "long raw", "longraw", "blob", "bfile":
		return "[]byte"

	// XML.
	case "xmltype":
		return nullableString()

	case "any":
		return "interface{}"
	}

	// Prefix-based matches for parameterized types (e.g. VARCHAR2(100),
	// TIMESTAMP(6), TIMESTAMP WITH TIME ZONE, INTERVAL ...).
	switch {
	case strings.HasPrefix(dt, "varchar2"),
		strings.HasPrefix(dt, "nvarchar2"),
		strings.HasPrefix(dt, "varchar"),
		strings.HasPrefix(dt, "char"),
		strings.HasPrefix(dt, "nchar"),
		strings.HasPrefix(dt, "character"),
		strings.HasPrefix(dt, "clob"),
		strings.HasPrefix(dt, "nclob"),
		strings.HasPrefix(dt, "long"),
		strings.HasPrefix(dt, "string"),
		strings.HasPrefix(dt, "rowid"),
		strings.HasPrefix(dt, "urowid"):
		return nullableString()

	case strings.HasPrefix(dt, "timestamp"):
		return nullableTime()

	case strings.HasPrefix(dt, "interval"):
		// Oracle INTERVAL types are surfaced as strings by go-ora.
		return nullableString()

	case strings.HasPrefix(dt, "number"),
		strings.HasPrefix(dt, "numeric"),
		strings.HasPrefix(dt, "decimal"),
		strings.HasPrefix(dt, "dec"),
		strings.HasPrefix(dt, "float"):
		return nullableFloat64()

	case strings.HasPrefix(dt, "raw"):
		return "[]byte"
	}

	if debug.Active {
		log.Printf("unknown Oracle type: %s\n", dt)
	}
	return "interface{}"
}

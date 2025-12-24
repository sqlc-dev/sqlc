package golang

import (
	"log"
	"strings"

	"github.com/sqlc-dev/sqlc/internal/codegen/golang/opts"
	"github.com/sqlc-dev/sqlc/internal/codegen/sdk"
	"github.com/sqlc-dev/sqlc/internal/debug"
	"github.com/sqlc-dev/sqlc/internal/plugin"
)

func clickhouseType(req *plugin.GenerateRequest, options *opts.Options, col *plugin.Column) string {
	dt := strings.ToLower(sdk.DataType(col.Type))
	notNull := col.NotNull || col.IsArray
	emitPointersForNull := options.EmitPointersForNullTypes

	// Handle Nullable wrapper
	if strings.HasPrefix(dt, "nullable(") && strings.HasSuffix(dt, ")") {
		dt = dt[9 : len(dt)-1]
		notNull = false
	}

	// Handle LowCardinality wrapper
	if strings.HasPrefix(dt, "lowcardinality(") && strings.HasSuffix(dt, ")") {
		dt = dt[15 : len(dt)-1]
	}

	switch dt {
	// Integer types
	case "int8":
		if notNull {
			return "int8"
		}
		if emitPointersForNull {
			return "*int8"
		}
		return "sql.NullInt16" // No sql.NullInt8, use Int16

	case "int16":
		if notNull {
			return "int16"
		}
		if emitPointersForNull {
			return "*int16"
		}
		return "sql.NullInt16"

	case "int32":
		if notNull {
			return "int32"
		}
		if emitPointersForNull {
			return "*int32"
		}
		return "sql.NullInt32"

	case "int64":
		if notNull {
			return "int64"
		}
		if emitPointersForNull {
			return "*int64"
		}
		return "sql.NullInt64"

	case "uint8":
		if notNull {
			return "uint8"
		}
		if emitPointersForNull {
			return "*uint8"
		}
		return "sql.NullInt16" // No sql.NullUint8

	case "uint16":
		if notNull {
			return "uint16"
		}
		if emitPointersForNull {
			return "*uint16"
		}
		return "sql.NullInt32" // No sql.NullUint16

	case "uint32":
		if notNull {
			return "uint32"
		}
		if emitPointersForNull {
			return "*uint32"
		}
		return "sql.NullInt64" // No sql.NullUint32

	case "uint64":
		if notNull {
			return "uint64"
		}
		if emitPointersForNull {
			return "*uint64"
		}
		// Note: uint64 doesn't fit in sql.NullInt64 for large values
		return "sql.NullInt64"

	// Float types
	case "float32":
		if notNull {
			return "float32"
		}
		if emitPointersForNull {
			return "*float32"
		}
		return "sql.NullFloat64"

	case "float64":
		if notNull {
			return "float64"
		}
		if emitPointersForNull {
			return "*float64"
		}
		return "sql.NullFloat64"

	// String types
	case "string":
		if notNull {
			return "string"
		}
		if emitPointersForNull {
			return "*string"
		}
		return "sql.NullString"

	// Boolean type
	case "bool", "boolean":
		if notNull {
			return "bool"
		}
		if emitPointersForNull {
			return "*bool"
		}
		return "sql.NullBool"

	// Date and time types
	case "date", "date32":
		if notNull {
			return "time.Time"
		}
		if emitPointersForNull {
			return "*time.Time"
		}
		return "sql.NullTime"

	case "datetime", "datetime64":
		if notNull {
			return "time.Time"
		}
		if emitPointersForNull {
			return "*time.Time"
		}
		return "sql.NullTime"

	// UUID type
	case "uuid":
		if notNull {
			return "uuid.UUID"
		}
		if emitPointersForNull {
			return "*uuid.UUID"
		}
		return "uuid.NullUUID"

	// JSON type
	case "json":
		return "json.RawMessage"

	// Any type (for unknown types)
	case "any":
		return "interface{}"

	default:
		// Handle FixedString(N)
		if strings.HasPrefix(dt, "fixedstring") {
			if notNull {
				return "string"
			}
			if emitPointersForNull {
				return "*string"
			}
			return "sql.NullString"
		}

		// Handle Decimal types
		if strings.HasPrefix(dt, "decimal") {
			if notNull {
				return "float64"
			}
			if emitPointersForNull {
				return "*float64"
			}
			return "sql.NullFloat64"
		}

		// Handle Array types
		if strings.HasPrefix(dt, "array(") && strings.HasSuffix(dt, ")") {
			innerType := dt[6 : len(dt)-1]
			innerCol := &plugin.Column{
				Type:    &plugin.Identifier{Name: innerType},
				NotNull: true,
			}
			return "[]" + clickhouseType(req, options, innerCol)
		}

		// Handle Enum types
		if strings.HasPrefix(dt, "enum8") || strings.HasPrefix(dt, "enum16") {
			if notNull {
				return "string"
			}
			if emitPointersForNull {
				return "*string"
			}
			return "sql.NullString"
		}

		// Handle Map types
		if strings.HasPrefix(dt, "map(") {
			return "map[string]interface{}"
		}

		// Handle Tuple types
		if strings.HasPrefix(dt, "tuple(") {
			return "interface{}"
		}

		if debug.Active {
			log.Printf("unknown ClickHouse type: %s\n", dt)
		}

		return "interface{}"
	}
}

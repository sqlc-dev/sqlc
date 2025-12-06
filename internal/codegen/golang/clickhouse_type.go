package golang

import (
	"strings"

	"github.com/sqlc-dev/sqlc/internal/codegen/golang/opts"
	"github.com/sqlc-dev/sqlc/internal/codegen/sdk"
	"github.com/sqlc-dev/sqlc/internal/plugin"
)

func clickhouseType(req *plugin.GenerateRequest, options *opts.Options, col *plugin.Column) string {
	columnType := sdk.DataType(col.Type)
	notNull := col.NotNull || col.IsArray

	// Check if we're using the native ClickHouse driver
	driver := parseDriver(options.SqlPackage)
	useNativeDriver := driver.IsClickHouse()

	switch columnType {

	// String types
	case "string", "varchar", "text", "char", "fixedstring":
		if useNativeDriver {
			// Native driver uses *string for nullable
			if notNull {
				return "string"
			}
			if options.EmitPointersForNullTypes {
				return "*string"
			}
			return "sql.NullString"
		}
		if notNull {
			return "string"
		}
		return "sql.NullString"

	// Integer types - UInt variants (unsigned)
	case "uint8":
		if useNativeDriver {
			if notNull {
				return "uint8"
			}
			if options.EmitPointersForNullTypes {
				return "*uint8"
			}
			return "sql.NullInt16"
		}
		if notNull {
			return "uint8"
		}
		return "sql.NullInt16" // database/sql doesn't have NullUint8

	case "uint16":
		if useNativeDriver {
			if notNull {
				return "uint16"
			}
			if options.EmitPointersForNullTypes {
				return "*uint16"
			}
			return "sql.NullInt32"
		}
		if notNull {
			return "uint16"
		}
		return "sql.NullInt32" // database/sql doesn't have NullUint16

	case "uint32":
		if useNativeDriver {
			if notNull {
				return "uint32"
			}
			if options.EmitPointersForNullTypes {
				return "*uint32"
			}
			return "sql.NullInt64"
		}
		if notNull {
			return "uint32"
		}
		return "sql.NullInt64" // database/sql doesn't have NullUint32

	case "uint64":
		if useNativeDriver {
			if notNull {
				return "uint64"
			}
			if options.EmitPointersForNullTypes {
				return "*uint64"
			}
			return "sql.NullInt64"
		}
		if notNull {
			return "uint64"
		}
		return "string" // uint64 can overflow, use string for large values

	// Integer types - Int variants (signed)
	case "int8":
		if useNativeDriver {
			if notNull {
				return "int8"
			}
			if options.EmitPointersForNullTypes {
				return "*int8"
			}
			return "sql.NullInt16"
		}
		if notNull {
			return "int8"
		}
		return "sql.NullInt16"

	case "int16":
		if useNativeDriver {
			if notNull {
				return "int16"
			}
			if options.EmitPointersForNullTypes {
				return "*int16"
			}
			return "sql.NullInt16"
		}
		if notNull {
			return "int16"
		}
		return "sql.NullInt16"

	case "int32":
		if useNativeDriver {
			if notNull {
				return "int32"
			}
			if options.EmitPointersForNullTypes {
				return "*int32"
			}
			return "sql.NullInt32"
		}
		if notNull {
			return "int32"
		}
		return "sql.NullInt32"

	case "int64":
		if useNativeDriver {
			if notNull {
				return "int64"
			}
			if options.EmitPointersForNullTypes {
				return "*int64"
			}
			return "sql.NullInt64"
		}
		if notNull {
			return "int64"
		}
		return "sql.NullInt64"

	// Generic "integer" type (used for LIMIT/OFFSET parameters and other integer values)
	case "integer":
		if useNativeDriver {
			if notNull {
				return "int64"
			}
			if options.EmitPointersForNullTypes {
				return "*int64"
			}
			return "sql.NullInt64"
		}
		if notNull {
			return "int64"
		}
		return "sql.NullInt64"

	// Large integer types
	case "int128", "int256", "uint128", "uint256":
		// These are too large for standard Go integers, use string
		if notNull {
			return "string"
		}
		return "sql.NullString"

	// Floating point types
	case "float32", "real":
		if useNativeDriver {
			if notNull {
				return "float32"
			}
			if options.EmitPointersForNullTypes {
				return "*float32"
			}
			return "sql.NullFloat64"
		}
		if notNull {
			return "float32"
		}
		return "sql.NullFloat64" // database/sql doesn't have NullFloat32

	case "float64", "double precision", "double":
		if useNativeDriver {
			if notNull {
				return "float64"
			}
			if options.EmitPointersForNullTypes {
				return "*float64"
			}
			return "sql.NullFloat64"
		}
		if notNull {
			return "float64"
		}
		return "sql.NullFloat64"

	// Decimal types
	case "decimal":
		if notNull {
			return "string"
		}
		return "sql.NullString"

	// Date and time types
	case "date", "date32":
		if useNativeDriver {
			if notNull {
				return "time.Time"
			}
			if options.EmitPointersForNullTypes {
				return "*time.Time"
			}
			return "sql.NullTime"
		}
		if notNull {
			return "time.Time"
		}
		return "sql.NullTime"

	case "datetime", "datetime64", "timestamp":
		if useNativeDriver {
			if notNull {
				return "time.Time"
			}
			if options.EmitPointersForNullTypes {
				return "*time.Time"
			}
			return "sql.NullTime"
		}
		if notNull {
			return "time.Time"
		}
		return "sql.NullTime"

	// Boolean
	case "boolean", "bool":
		if useNativeDriver {
			if notNull {
				return "bool"
			}
			if options.EmitPointersForNullTypes {
				return "*bool"
			}
			return "sql.NullBool"
		}
		if notNull {
			return "bool"
		}
		return "sql.NullBool"

	// UUID
	case "uuid":
		if notNull {
			return "string"
		}
		return "sql.NullString"

	// IP address types
	case "ipv4", "ipv6":
		if notNull {
			return "netip.Addr"
		}
		if options.EmitPointersForNullTypes {
			return "*netip.Addr"
		}
		// Use a custom SQL null type for nullable IP addresses
		// For now, use pointer since netip.Addr doesn't have a nullable variant
		return "*netip.Addr"

	// JSON types
	case "json":
		return "json.RawMessage"

	// Arrays - ClickHouse array types
	case "array":
		if useNativeDriver {
			// Native driver has better array support
			// For now, still use generic until we have element type info
			return "[]interface{}"
		}
		return "[]interface{}" // Generic array type

	// Any/Unknown type
	case "any":
		return "interface{}"

	default:
		// Check if this is a map type (starts with "map[")
		// Map types come from the engine layer with full type information (e.g., "map[string]int64")
		if strings.HasPrefix(columnType, "map[") {
			if notNull {
				return columnType
			}
			// For nullable map types, wrap in pointer
			if options.EmitPointersForNullTypes {
				return "*" + columnType
			}
			// Otherwise treat as interface{} for nullable
			return "interface{}"
		}

		// Check for custom types (enums, etc.)
		for _, schema := range req.Catalog.Schemas {
			for _, enum := range schema.Enums {
				if enum.Name == columnType {
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
		}

		// Default fallback for unknown types
		return "interface{}"
	}
}

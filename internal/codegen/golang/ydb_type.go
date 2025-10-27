package golang

import (
	"log"
	"strings"

	"github.com/sqlc-dev/sqlc/internal/codegen/golang/opts"
	"github.com/sqlc-dev/sqlc/internal/codegen/sdk"
	"github.com/sqlc-dev/sqlc/internal/debug"
	"github.com/sqlc-dev/sqlc/internal/plugin"
)

func YDBType(req *plugin.GenerateRequest, options *opts.Options, col *plugin.Column) string {
	columnType := strings.ToLower(sdk.DataType(col.Type))
	isArray := col.IsArray || col.IsSqlcSlice
	notNull := (col.NotNull || isArray) && (!isNullableType(columnType) || isArray)
	emitPointersForNull := options.EmitPointersForNullTypes

	columnType = extractBaseType(columnType)

	// https://ydb.tech/docs/ru/yql/reference/types/
	// ydb-go-sdk doesn't support sql.Null* yet
	switch columnType {
	// decimal types
	case "bool":
		if notNull {
			return "bool"
		}
		if emitPointersForNull {
			return "*bool"
		}
		// return "sql.NullBool"
		return "*bool"

	case "int8":
		if notNull {
			return "int8"
		}
		if emitPointersForNull {
			return "*int8"
		}
		// // The database/sql package does not have a sql.NullInt8 type, so we
		// // use the smallest type they have which is NullInt16
		// return "sql.NullInt16"
		return "*int8"
	case "int16", "smallint":
		if notNull {
			return "int16"
		}
		if emitPointersForNull {
			return "*int16"
		}
		// return "sql.NullInt16"
		return "*int16"
	case "int", "int32", "integer":
		if notNull {
			return "int32"
		}
		if emitPointersForNull {
			return "*int32"
		}
		// return "sql.NullInt32"
		return "*int32"
	case "int64", "bigint":
		if notNull {
			return "int64"
		}
		if emitPointersForNull {
			return "*int64"
		}
		// return "sql.NullInt64"
		return "*int64"

	case "uint8":
		if notNull {
			return "uint8"
		}
		if emitPointersForNull {
			return "*uint8"
		}
		return "*uint8"
	case "uint16":
		if notNull {
			return "uint16"
		}
		if emitPointersForNull {
			return "*uint16"
		}
		return "*uint16"
	case "uint32":
		if notNull {
			return "uint32"
		}
		if emitPointersForNull {
			return "*uint32"
		}
		return "*uint32"
	case "uint64":
		if notNull {
			return "uint64"
		}
		if emitPointersForNull {
			return "*uint64"
		}
		return "*uint64"

	case "float":
		if notNull {
			return "float32"
		}
		if emitPointersForNull {
			return "*float32"
		}
		// The database/sql package does not have a sql.NullFloat32 type, so we
		// use the smallest type they have which is NullFloat64
		// return "sql.NullFloat64"
		return "*float32"
	case "double":
		if notNull {
			return "float64"
		}
		if emitPointersForNull {
			return "*float64"
		}
		// return "sql.NullFloat64"
		return "*float64"

	// string types
	case "utf8", "text":
		if notNull {
			return "string"
		}
		if emitPointersForNull {
			return "*string"
		}
		return "*string"

	// serial types
	case "smallserial", "serial2":
		if notNull {
			return "int16"
		}
		if emitPointersForNull {
			return "*int16"
		}
		// return "sql.NullInt16"
		return "*int16"

	case "serial", "serial4":
		if notNull {
			return "int32"
		}
		if emitPointersForNull {
			return "*int32"
		}
		// return "sql.NullInt32"
		return "*int32"

	case "bigserial", "serial8":
		if notNull {
			return "int64"
		}
		if emitPointersForNull {
			return "*int64"
		}
		// return "sql.NullInt64"
		return "*int64"

	case "json", "jsondocument":
		if notNull {
			return "string"
		}
		if emitPointersForNull {
			return "*string"
		}
		return "*string"

	case "date", "date32", "datetime", "timestamp", "tzdate", "tztimestamp", "tzdatetime", "datetime64", "timestamp64", "tzdate32", "tzdatetime64", "tztimestamp64":
		if notNull {
			return "time.Time"
		}
		if emitPointersForNull {
			return "*time.Time"
		}
		return "*time.Time"

	case "interval", "interval64":
		if notNull {
			return "time.Duration"
		}
		if emitPointersForNull {
			return "*time.Duration"
		}
		return "*time.Duration"

	case "uuid":
		if notNull {
			return "uuid.UUID"
		}
		if emitPointersForNull {
			return "*uuid.UUID"
		}
		return "*uuid.UUID"

	case "yson", "string":
		if notNull {
			return "[]byte"
		}
		if emitPointersForNull {
			return "*[]byte"
		}
		return "*[]byte"

	case "null":
		// return "sql.Null"
		return "interface{}"

	case "any":
		return "interface{}"

	default:
		if strings.HasPrefix(columnType, "decimal") {
			if notNull {
				return "types.Decimal"
			}

			if emitPointersForNull {
				return "*types.Decimal"
			}
			return "*types.Decimal"
		}

		if debug.Active {
			log.Printf("unknown YDB type: %s\n", columnType)
		}

		return "interface{}"
	}

}

// This function extracts the base type from optional types
func extractBaseType(typeStr string) string {
	if strings.HasPrefix(typeStr, "optional<") && strings.HasSuffix(typeStr, ">") {
		return strings.TrimSuffix(strings.TrimPrefix(typeStr, "optional<"), ">")
	}
	if strings.HasSuffix(typeStr, "?") {
		return strings.TrimSuffix(typeStr, "?")
	}
	return typeStr
}

func isNullableType(typeStr string) bool {
	return strings.HasPrefix(typeStr, "optional<") && strings.HasSuffix(typeStr, ">") || strings.HasSuffix(typeStr, "?")
}

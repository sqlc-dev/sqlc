package golang

import (
	"log"

	"github.com/sqlc-dev/sqlc/internal/codegen/golang/opts"
	"github.com/sqlc-dev/sqlc/internal/codegen/sdk"
	"github.com/sqlc-dev/sqlc/internal/debug"
	"github.com/sqlc-dev/sqlc/internal/plugin"
)

func mysqlType(req *plugin.GenerateRequest, options *opts.Options, col *plugin.Column) string {
	columnType := sdk.DataType(col.Type)
	notNull := col.NotNull || col.IsArray
	unsigned := col.Unsigned

	var genericNull *opts.ParsedGoType
	if options.GenericNullType != nil {
		genericNull = &options.GenericNullType.Parsed
	}

	switch columnType {

	case "varchar", "text", "char", "tinytext", "mediumtext", "longtext":
		if notNull {
			return "string"
		}
		if genericNull != nil {
			return genericNull.TypeName + "[string]"
		}
		return "sql.NullString"

	case "tinyint":
		if col.Length == 1 {
			if notNull {
				return "bool"
			}
			return "sql.NullBool"
		} else {
			baseType := "int32"
			if unsigned {
				baseType = "uint32"
			}

			if notNull {
				return baseType
			}
			if genericNull != nil {
				return genericNull.TypeName + "[" + baseType + "]"
			}
			return "sql.NullInt32"
		}

	case "smallint", "year":
		baseType := "int16"
		if unsigned {
			baseType = "uint16"
		}

		if notNull {
			return baseType
		}
		if genericNull != nil {
			return genericNull.TypeName + "[" + baseType + "]"
		}
		return "sql.NullInt16"

	case "int", "integer", "mediumint":
		baseType := "int32"
		if unsigned {
			baseType = "uint32"
		}

		if notNull {
			return baseType
		}
		if genericNull != nil {
			return genericNull.TypeName + "[" + baseType + "]"
		}
		return "sql.NullInt32"

	case "bigint":
		baseType := "int64"
		if unsigned {
			baseType = "uint64"
		}

		if notNull {
			return baseType
		}
		if genericNull != nil {
			return genericNull.TypeName + "[" + baseType + "]"
		}
		return "sql.NullInt64"

	case "blob", "binary", "varbinary", "tinyblob", "mediumblob", "longblob":
		if notNull {
			return "[]byte"
		}
		if genericNull != nil {
			return genericNull.TypeName + "[[]byte]"
		}
		return "sql.NullString"

	case "double", "double precision", "real", "float":
		if notNull {
			return "float64"
		}
		if genericNull != nil {
			return genericNull.TypeName + "[float64]"
		}
		return "sql.NullFloat64"

	case "decimal", "dec", "fixed":
		if notNull {
			return "string"
		}
		if genericNull != nil {
			return genericNull.TypeName + "[string]"
		}
		return "sql.NullString"

	case "enum":
		// TODO: Proper Enum support
		return "string"

	case "date", "timestamp", "datetime", "time":
		if notNull {
			return "time.Time"
		}
		if genericNull != nil {
			return genericNull.TypeName + "[time.Time]"
		}
		return "sql.NullTime"

	case "boolean", "bool":
		if notNull {
			return "bool"
		}
		if genericNull != nil {
			return genericNull.TypeName + "[bool]"
		}
		return "sql.NullBool"

	case "json":
		return "json.RawMessage"

	case "any":
		return "interface{}"

	default:
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
		if debug.Active {
			log.Printf("Unknown MySQL type: %s\n", columnType)
		}
		return "interface{}"

	}
}

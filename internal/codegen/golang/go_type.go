package golang

import (
	"log"

	"github.com/kyleconroy/sqlc/internal/compiler"
	"github.com/kyleconroy/sqlc/internal/config"
	"github.com/kyleconroy/sqlc/internal/sql/catalog"
)

func goType(r *compiler.Result, col *compiler.Column, settings config.CombinedSettings) string {
	// package overrides have a higher precedence
	for _, oride := range settings.Overrides {
		if oride.Column != "" && oride.ColumnName == col.Name && sameTableName(col.Table, oride.Table) {
			return oride.GoTypeName
		}
	}
	typ := goInnerType(r, col, settings)
	if col.IsArray {
		return "[]" + typ
	}
	return typ
}

func goInnerType(r *compiler.Result, col *compiler.Column, settings config.CombinedSettings) string {
	columnType := col.DataType
	notNull := col.NotNull || col.IsArray

	// package overrides have a higher precedence
	for _, oride := range settings.Overrides {
		if oride.DBType != "" && oride.DBType == columnType && oride.Null != notNull {
			return oride.GoTypeName
		}
	}

	// TODO: Extend the engine interface to handle types
	switch settings.Package.Engine {
	case config.EngineMySQL, config.EngineXDolphin:
		return mysqlType(r, col, settings)
	case config.EnginePostgreSQL:
		return postgresType(r, col, settings)
	case config.EngineXLemon:
		return postgresType(r, col, settings)
	default:
		return "interface{}"
	}
}

func postgresType(r *compiler.Result, col *compiler.Column, settings config.CombinedSettings) string {
	columnType := col.DataType
	notNull := col.NotNull || col.IsArray

	switch columnType {
	case "serial", "pg_catalog.serial4":
		if notNull {
			return "int32"
		}
		return "sql.NullInt32"

	case "bigserial", "pg_catalog.serial8":
		if notNull {
			return "int64"
		}
		return "sql.NullInt64"

	case "smallserial", "pg_catalog.serial2":
		return "int16"

	case "integer", "int", "int4", "pg_catalog.int4":
		if notNull {
			return "int32"
		}
		return "sql.NullInt32"

	case "bigint", "pg_catalog.int8":
		if notNull {
			return "int64"
		}
		return "sql.NullInt64"

	case "smallint", "pg_catalog.int2":
		return "int16"

	case "float", "double precision", "pg_catalog.float8":
		if notNull {
			return "float64"
		}
		return "sql.NullFloat64"

	case "real", "pg_catalog.float4":
		if notNull {
			return "float32"
		}
		return "sql.NullFloat64" // TODO: Change to sql.NullFloat32 after updating the go.mod file

	case "pg_catalog.numeric", "money":
		// Since the Go standard library does not have a decimal type, lib/pq
		// returns numerics as strings.
		//
		// https://github.com/lib/pq/issues/648
		if notNull {
			return "string"
		}
		return "sql.NullString"

	case "bool", "pg_catalog.bool":
		if notNull {
			return "bool"
		}
		return "sql.NullBool"

	case "json", "jsonb":
		return "json.RawMessage"

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
		return "uuid.UUID"

	case "inet":
		return "net.IP"

	case "macaddr", "macaddr8":
		return "net.HardwareAddr"

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

	case "void":
		// A void value always returns NULL. Since there is no built-in NULL
		// value into the SQL package, we'll use sql.NullBool
		return "sql.NullBool"

	case "any":
		return "interface{}"

	default:
		rel, err := compiler.ParseRelationString(columnType)
		if err != nil {
			// TODO: Should this actually return an error here?
			return "interface{}"
		}
		if rel.Schema == "" {
			rel.Schema = r.Catalog.DefaultSchema
		}

		for _, schema := range r.Catalog.Schemas {
			if schema.Name == "pg_catalog" {
				continue
			}
			for _, typ := range schema.Types {
				switch t := typ.(type) {
				case *catalog.Enum:
					if rel.Name == t.Name && rel.Schema == schema.Name {
						if schema.Name == r.Catalog.DefaultSchema {
							return StructName(t.Name, settings)
						}
						return StructName(schema.Name+"_"+t.Name, settings)
					}
				case *catalog.CompositeType:
					if notNull {
						return "string"
					}
					return "sql.NullString"
				}
			}
		}

		log.Printf("unknown PostgreSQL type: %s\n", columnType)
		return "interface{}"
	}
}

func mysqlType(r *compiler.Result, col *compiler.Column, settings config.CombinedSettings) string {
	columnType := col.DataType
	notNull := col.NotNull || col.IsArray

	switch columnType {

	case "varchar", "text", "char", "tinytext", "mediumtext", "longtext":
		if notNull {
			return "string"
		}
		return "sql.NullString"

	case "int", "integer", "smallint", "mediumint", "year":
		if notNull {
			return "int32"
		}
		return "sql.NullInt32"

	case "bigint":
		if notNull {
			return "int64"
		}
		return "sql.NullInt64"

	case "blob", "binary", "varbinary", "tinyblob", "mediumblob", "longblob":
		return "[]byte"

	case "double", "double precision", "real":
		if notNull {
			return "float64"
		}
		return "sql.NullFloat64"

	case "decimal", "dec", "fixed":
		if notNull {
			return "string"
		}
		return "sql.NullString"

	case "enum":
		// TODO: Proper Enum support
		return "string"

	case "date", "timestamp", "datetime", "time":
		if notNull {
			return "time.Time"
		}
		return "sql.NullTime"

	case "boolean", "bool", "tinyint":
		if notNull {
			return "bool"
		}
		return "sql.NullBool"

	default:
		log.Printf("unknown MySQL type: %s\n", columnType)
		return "interface{}"

	}
}

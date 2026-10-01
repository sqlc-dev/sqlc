package golang

import (
	"github.com/sqlc-dev/sqlc/internal/codegen/golang/opts"
	"github.com/sqlc-dev/sqlc/internal/plugin"
)

// mssqlType maps a SQL Server type to the Go type
// github.com/microsoft/go-mssqldb hands database/sql for it. Exact numerics
// come back as their decimal text, so they map to string the way
// PostgreSQL's numeric does under database/sql. A uniqueidentifier maps to
// the driver's own type, which keeps SQL Server's byte order straight. An
// alias type created with CREATE TYPE ... FROM is reported by its own name,
// which the catalog does not yet resolve to its base, so it maps to any
// until it does.
func mssqlType(req *plugin.GenerateRequest, options *opts.Options, col *plugin.Column) string {
	t, nullable := columnType(col)
	null := func(typ, nullTyp string) string {
		return nullableGoType(options, nullable, typ, nullTyp)
	}

	switch typeExprName(t) {
	case "bigint":
		return null("int64", "sql.NullInt64")

	case "int":
		return null("int32", "sql.NullInt32")

	case "smallint":
		return null("int16", "sql.NullInt16")

	case "tinyint":
		return null("uint8", "sql.NullByte")

	case "bit":
		return null("bool", "sql.NullBool")

	case "decimal", "numeric", "money", "smallmoney":
		return null("string", "sql.NullString")

	case "float":
		return null("float64", "sql.NullFloat64")

	case "real":
		return null("float32", "sql.NullFloat64")

	case "char", "varchar", "nchar", "nvarchar", "text", "ntext", "sysname", "xml", "json":
		return null("string", "sql.NullString")

	case "binary", "varbinary", "image", "rowversion", "timestamp":
		return "[]byte"

	case "date", "datetime", "datetime2", "smalldatetime", "datetimeoffset", "time":
		return null("time.Time", "sql.NullTime")

	case "uniqueidentifier":
		return null("mssql.UniqueIdentifier", "mssql.NullUniqueIdentifier")

	case "geography", "geometry", "hierarchyid", "vector":
		return "[]byte"

	default:
		return "any"
	}
}

package postgres

import "log"

type Schema struct {
	Tables []Table
}

type Table struct {
	GoName  string
	Name    string
	Columns []Column
}

type Column struct {
	GoName  string
	Name    string
	Type    string
	NotNull bool
}

func (c Column) GoType() string {
	// {{.GoName}} {{if .NotNull }}string{{else}}sql.NullString{{end}}
	switch c.Type {
	case "text":
		if c.NotNull {
			return "string"
		} else {
			return "sql.NullString"
		}
	case "serial":
		return "int"
	case "integer":
		return "int"
	case "bool":
		if c.NotNull {
			return "bool"
		} else {
			return "sql.NullBool"
		}
	case "pg_catalog.bool":
		if c.NotNull {
			return "bool"
		} else {
			return "sql.NullBool"
		}
	case "pg_catalog.int4":
		return "int"
	case "pg_catalog.int8":
		return "int"
	case "pg_catalog.timestamp":
		return "time.Time"
	case "pg_catalog.varchar":
		if c.NotNull {
			return "string"
		} else {
			return "sql.NullString"
		}
	default:
		log.Printf("unknown Postgres type: %s\n", c.Type)
		return "interface{}"
	}
}

package postgres

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
	case "pg_catalog.timestamp":
		return "time.Time"
	default:
		return "interface{}"
	}
}

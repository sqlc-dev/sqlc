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

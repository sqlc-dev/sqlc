package pg

func NewCatalog() Catalog {
	return Catalog{
		Schemas: map[string]Schema{
			"public": Schema{
				Tables: map[string]Table{},
				Enums:  map[string]Enum{},
			},
		},
	}
}

type FQN struct {
	Catalog string
	Schema  string
	Rel     string
}

type Catalog struct {
	Schemas map[string]Schema
}

type Schema struct {
	Name   string
	Tables map[string]Table
	Enums  map[string]Enum
}

type Table struct {
	Name    string
	Columns []Column
}

type Column struct {
	Name     string
	DataType string
	NotNull  bool
}

type Enum struct {
	Name string
	Vals []string
}

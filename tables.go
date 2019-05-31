package strongdb

// Other package
type Query int

const (
	Zero Query = iota
	One
	Many
)

type Column struct {
	Type string
	PK   bool
}

type Base struct {
	Table   string
	Columns map[string]Column
}

var Endpoint = Base{
	Table: "endpoint",
	Columns: map[string]Column{
		"id": {
			Type: "bytea(20",
			PK:   true,
		},
		"account_id": {
			// TODO: Add a type for types
			Type: "int",
		},
	},
}

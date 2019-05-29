package main

// Other package
type Query int

const (
	Zero Query = iota
	One
	Many
)

// Magic
// type EndpointConfig struct {
// 	ID        []byte          `sql:"id,pk"`
// 	AccountID int64           `sql:"account_id"`
// 	Settings  json.RawMessage `sql:"settings"`
// }
//
// func (ec *EndpointConfig) Table() string {
// 	return "endpoint_config"
// }
//
// func (ec *EndpointConfig) Fetch() (Query, string) {
// 	return Many, `SELECT * FROM {{.Table}} WHERE id = {} AND account_id = {acct}`
// }

type Column struct {
	Type string
	PK   bool
}

type Base struct {
	Table   string
	Columns map[string]Column
}

var EndpointConfig = Base{
	Table: "endpoint_config",
	Columns: map[string]Column{
		"id": {
			Type: "bytea(20",
			PK:   true,
		},
		"account_id": {
			Type: "Integer",
		},
		"settings": {
			Type: "jsonb",
		},
	},
}

func main() {
}

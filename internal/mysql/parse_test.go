package mysql

import (
	"encoding/json"
	"io/ioutil"
	"reflect"
	"testing"

	"github.com/davecgh/go-spew/spew"
	"github.com/kyleconroy/sqlc/internal/dinosql"
	"vitess.io/vitess/go/vt/sqlparser"
)

func init() {
	initMockSchema()
}

const filename = "test_data/queries.sql"
const configPath = "test_data/sqlc.json"

var mockSettings = dinosql.GenerateSettings{
	Version: "1",
	Packages: []dinosql.PackageSettings{
		dinosql.PackageSettings{
			Name: "db",
		},
	},
	Overrides: []dinosql.Override{},
}

func TestParseConfig(t *testing.T) {
	blob, err := ioutil.ReadFile(configPath)
	if err != nil {
		t.Fatal(err)
	}

	var settings dinosql.GenerateSettings
	if err := json.Unmarshal(blob, &settings); err != nil {
		t.Fatal(err)
	}
}

func TestGeneratePkg(t *testing.T) {
	_, err := GeneratePkg(mockSettings.Packages[0].Name, filename, filename, mockSettings)
	if err != nil {
		t.Fatal(err)
	}
}

func keep(interface{}) {}

var mockSchema *Schema

func initMockSchema() {
	var schemaMap = make(map[string][]*sqlparser.ColumnDefinition)
	mockSchema = &Schema{
		tables: schemaMap,
	}
	schemaMap["users"] = []*sqlparser.ColumnDefinition{
		&sqlparser.ColumnDefinition{
			Name: sqlparser.NewColIdent("first_name"),
			Type: sqlparser.ColumnType{
				Type:    "varchar",
				NotNull: true,
			},
		},
		&sqlparser.ColumnDefinition{
			Name: sqlparser.NewColIdent("last_name"),
			Type: sqlparser.ColumnType{
				Type:    "varchar",
				NotNull: false,
			},
		},
		&sqlparser.ColumnDefinition{
			Name: sqlparser.NewColIdent("id"),
			Type: sqlparser.ColumnType{
				Type:          "int",
				NotNull:       true,
				Autoincrement: true,
			},
		},
		&sqlparser.ColumnDefinition{
			Name: sqlparser.NewColIdent("age"),
			Type: sqlparser.ColumnType{
				Type:    "int",
				NotNull: true,
			},
		},
	}
	schemaMap["orders"] = []*sqlparser.ColumnDefinition{
		&sqlparser.ColumnDefinition{
			Name: sqlparser.NewColIdent("id"),
			Type: sqlparser.ColumnType{
				Type:          "int",
				NotNull:       true,
				Autoincrement: true,
			},
		},
		&sqlparser.ColumnDefinition{
			Name: sqlparser.NewColIdent("price"),
			Type: sqlparser.ColumnType{
				Type:          "DECIMAL(13, 4)",
				NotNull:       true,
				Autoincrement: true,
			},
		},
		&sqlparser.ColumnDefinition{
			Name: sqlparser.NewColIdent("user_id"),
			Type: sqlparser.ColumnType{
				Type:    "int",
				NotNull: true,
			},
		},
	}
}

func filterCols(allCols []*sqlparser.ColumnDefinition, tableNames map[string]struct{}) []*sqlparser.ColumnDefinition {
	filteredCols := []*sqlparser.ColumnDefinition{}
	for _, col := range allCols {
		if _, ok := tableNames[col.Name.String()]; ok {
			filteredCols = append(filteredCols, col)
		}
	}
	return filteredCols
}

func TestParseSelect(t *testing.T) {
	type expected struct {
		query  string
		schema *Schema
	}
	type testCase struct {
		input  expected
		output *Query
	}
	tests := []testCase{
		testCase{
			input: expected{
				query: `/* name: GetCount :one */
					SELECT id my_id, COUNT(id) id_count FROM users WHERE id > 4`,
				schema: mockSchema,
			},
			output: &Query{
				SQL: "select id as my_id, COUNT(id) as id_count from users where id > 4",
				Columns: []*sqlparser.ColumnDefinition{
					&sqlparser.ColumnDefinition{
						Name: sqlparser.NewColIdent("my_id"),
						Type: sqlparser.ColumnType{
							Type:          "int",
							NotNull:       true,
							Autoincrement: true,
						},
					},
					&sqlparser.ColumnDefinition{
						Name: sqlparser.NewColIdent("id_count"),
						Type: sqlparser.ColumnType{
							Type:    "int",
							NotNull: true,
						},
					},
				},
				Params:           []*Param{},
				Name:             "GetCount",
				Cmd:              ":one",
				defaultTableName: "users",
				schemaLookup:     mockSchema,
			},
		},
		testCase{
			input: expected{
				query: `/* name: GetNameByID :one */
									SELECT first_name, last_name FROM users WHERE id = ?`,
				schema: mockSchema,
			},
			output: &Query{
				SQL:     `select first_name, last_name from users where id = ?`,
				Columns: filterCols(mockSchema.tables["users"], map[string]struct{}{"first_name": struct{}{}, "last_name": struct{}{}}),
				Params: []*Param{
					&Param{
						originalName: ":v1",
						name:         "id",
						typ:          "int",
					}},
				Name:             "GetNameByID",
				Cmd:              ":one",
				defaultTableName: "users",
				schemaLookup:     mockSchema,
			},
		},
		testCase{
			input: expected{
				query: `/* name: GetAll :many */
				SELECT * FROM users;`,
				schema: mockSchema,
			},
			output: &Query{
				SQL:              "select first_name, last_name, id, age from users",
				Columns:          mockSchema.tables["users"],
				Params:           []*Param{},
				Name:             "GetAll",
				Cmd:              ":many",
				defaultTableName: "users",
				schemaLookup:     mockSchema,
			},
		},
		testCase{
			input: expected{
				query: `/* name: GetAllUsersOrders :many */
				SELECT u.id user_id, u.first_name, o.price, o.id order_id
				FROM orders o LEFT JOIN users u ON u.id = o.user_id`,
				schema: mockSchema,
			},
			output: &Query{
				SQL: "select u.id as user_id, u.first_name, o.price, o.id as order_id from orders as o left join users as u on u.id = o.user_id",
				Columns: []*sqlparser.ColumnDefinition{
					&sqlparser.ColumnDefinition{
						Name: sqlparser.NewColIdent("user_id"),
						Type: sqlparser.ColumnType{
							Type:          "int",
							Autoincrement: true,
							NotNull:       false, // beause of the left join
						},
					},
					&sqlparser.ColumnDefinition{
						Name: sqlparser.NewColIdent("first_name"),
						Type: sqlparser.ColumnType{
							Type:    "varchar",
							NotNull: false, // because of left join
						},
					},
					&sqlparser.ColumnDefinition{
						Name: sqlparser.NewColIdent("price"),
						Type: sqlparser.ColumnType{
							Type:          "DECIMAL(13, 4)",
							Autoincrement: true,
							NotNull:       true,
						},
					},
					&sqlparser.ColumnDefinition{
						Name: sqlparser.NewColIdent("order_id"),
						Type: sqlparser.ColumnType{
							Type:          "int",
							Autoincrement: true,
							NotNull:       true,
						},
					},
				},
				Params:           []*Param{},
				Name:             "GetAllUsersOrders",
				Cmd:              ":many",
				defaultTableName: "", // TODO: verify that this is desired behaviour
				schemaLookup:     mockSchema,
			},
		},
	}

	for ix, testCase := range tests {
		q, err := parseQueryString(testCase.input.query, testCase.input.schema, mockSettings)
		if err != nil {
			t.Errorf("Parsing failed with query: [%v]\n%v", testCase.input.query, err)
		}

		err = q.parseNameAndCmd()
		if err != nil {
			t.Errorf("Parsing failed with query: [%v]\n%v", testCase.input.query, err)
		}
		if !reflect.DeepEqual(testCase.output, q) {
			t.Errorf("Parsing query returned differently than expected. Index: %v", ix)
			t.Logf("Expected: %v\nResult: %v\n", spew.Sdump(testCase.output), spew.Sdump(q))
		}
	}
}

func TestParseLeadingComment(t *testing.T) {
	type expected struct {
		name string
		cmd  string
	}
	type testCase struct {
		input  string
		output expected
	}

	tests := []testCase{
		testCase{
			input:  "/* name: GetPeopleByID :many */",
			output: expected{name: "GetPeopleByID", cmd: ":many"},
		},
	}

	for _, tCase := range tests {
		qu := &Query{}
		err := qu.parseLeadingComment(tCase.input)
		if err != nil {
			t.Errorf("Failed to parse leading comment %v", err)
		}
		if qu.Name != tCase.output.name || qu.Cmd != tCase.output.cmd {
			t.Errorf("Leading comment parser returned unexpcted result: %v\n:\n Expected: [%v]\nRecieved:[%v]\n",
				err, spew.Sdump(tCase.output), spew.Sdump(qu))
		}

	}
}

func TestSchemaLookup(t *testing.T) {
	firstNameColDfn, err := mockSchema.schemaLookup("users", "first_name")
	if err != nil {
		t.Errorf("Failed to get column schema from mock schema: %v", err)
	}

	expected := filterCols(mockSchema.tables["users"], map[string]struct{}{"first_name": struct{}{}})
	if !reflect.DeepEqual(firstNameColDfn, expected[0]) {
		t.Errorf("Table schema lookup returned unexpected result")
	}
}

func TestParseInsertUpdate(t *testing.T) {
	type expected struct {
		query  string
		schema *Schema
	}
	type testCase struct {
		input  expected
		output *Query
	}
	query1 := "/* name: InsertNewUser :exec */\nINSERT INTO users (first_name, last_name) VALUES (?, ?)"
	query2 := "/* name: UpdateUserAt :exec */\nUPDATE users SET first_name = ?, last_name = ? WHERE id > ? AND first_name = ? LIMIT 3"
	tests := []testCase{
		testCase{
			input: expected{
				query:  query1,
				schema: mockSchema,
			},
			output: &Query{
				SQL:     "insert into users(first_name, last_name) values (?, ?)",
				Columns: nil,
				Params: []*Param{
					&Param{
						originalName: ":v1",
						name:         "first_name",
						typ:          "string",
					},
					&Param{
						originalName: ":v2",
						name:         "last_name",
						typ:          "sql.NullString",
					},
				},
				Name:             "InsertNewUser",
				Cmd:              ":exec",
				defaultTableName: "users",
				schemaLookup:     mockSchema,
			},
		},
		testCase{
			input: expected{
				query:  query2,
				schema: mockSchema,
			},
			output: &Query{
				SQL:     "update users set first_name = ?, last_name = ? where id > ? and first_name = ? limit 3",
				Columns: nil,
				Params: []*Param{
					&Param{
						originalName: ":v1",
						name:         "first_name",
						typ:          "string",
					},
					&Param{
						originalName: ":v2",
						name:         "last_name",
						typ:          "sql.NullString",
					},
					&Param{
						originalName: ":v3",
						name:         "id",
						typ:          "int",
					},
					&Param{
						originalName: ":v4",
						name:         "first_name",
						typ:          "string",
					},
				},
				Name:             "UpdateUserAt",
				Cmd:              ":exec",
				defaultTableName: "users",
				schemaLookup:     mockSchema,
			},
		},
	}

	for ix, testCase := range tests {
		q, err := parseQueryString(testCase.input.query, testCase.input.schema, mockSettings)
		if err != nil {
			t.Errorf("Parsing failed with query: [%v]\n", err)
			continue
		}

		err = q.parseNameAndCmd()
		if err != nil {
			t.Errorf("Parsing failed with query index: %d: [%v]\n", ix, testCase.input.query)
		}
		if !reflect.DeepEqual(testCase.output, q) {
			t.Errorf("Parsing query returned differently than expected.")
			t.Logf("Expected: %v\nResult: %v\n", spew.Sdump(testCase.output), spew.Sdump(q))
		}
	}
}

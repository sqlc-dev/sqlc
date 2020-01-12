package mysql

import (
	"encoding/json"
	"io/ioutil"
	"reflect"
	"testing"

	"github.com/davecgh/go-spew/spew"
	"github.com/google/go-cmp/cmp"
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

func filterCols(allCols []*sqlparser.ColumnDefinition, colNames map[string]string) []Column {
	cols := []Column{}
	for _, col := range allCols {
		if table, ok := colNames[col.Name.String()]; ok {
			cols = append(cols, Column{
				col,
				table,
			})
		}
	}
	return cols
}

func TestParseSelect(t *testing.T) {
	type expected struct {
		query  string
		schema *Schema
	}
	type testCase struct {
		name   string
		input  expected
		output *Query
	}
	tests := []testCase{
		testCase{
			name: "get_count",
			input: expected{
				query: `/* name: GetCount :one */
					SELECT id my_id, COUNT(id) id_count FROM users WHERE id > 4`,
				schema: mockSchema,
			},
			output: &Query{
				SQL: "select id as my_id, COUNT(id) as id_count from users where id > 4",
				Columns: []Column{
					Column{
						&sqlparser.ColumnDefinition{
							Name: sqlparser.NewColIdent("my_id"),
							Type: sqlparser.ColumnType{
								Type:          "int",
								NotNull:       true,
								Autoincrement: true,
							},
						},
						"users",
					},
					Column{
						&sqlparser.ColumnDefinition{
							Name: sqlparser.NewColIdent("id_count"),
							Type: sqlparser.ColumnType{
								Type:    "int",
								NotNull: true,
							},
						},
						"",
					},
				},
				Params:           []*Param{},
				Name:             "GetCount",
				Cmd:              ":one",
				DefaultTableName: "users",
				SchemaLookup:     mockSchema,
			},
		},
		testCase{
			name: "get_name_by_id",
			input: expected{
				query: `/* name: GetNameByID :one */
									SELECT first_name, last_name FROM users WHERE id = ?`,
				schema: mockSchema,
			},
			output: &Query{
				SQL:     `select first_name, last_name from users where id = ?`,
				Columns: filterCols(mockSchema.tables["users"], map[string]string{"first_name": "users", "last_name": "users"}),
				Params: []*Param{
					&Param{
						OriginalName: ":v1",
						Name:         "id",
						Typ:          "int",
					}},
				Name:             "GetNameByID",
				Cmd:              ":one",
				DefaultTableName: "users",
				SchemaLookup:     mockSchema,
			},
		},
		testCase{
			name: "get_all",
			input: expected{
				query: `/* name: GetAll :many */
				SELECT * FROM users;`,
				schema: mockSchema,
			},
			output: &Query{
				SQL:              "select first_name, last_name, id, age from users",
				Columns:          filterCols(mockSchema.tables["users"], map[string]string{"first_name": "users", "last_name": "users", "id": "users", "age": "users"}),
				Params:           []*Param{},
				Name:             "GetAll",
				Cmd:              ":many",
				DefaultTableName: "users",
				SchemaLookup:     mockSchema,
			},
		},
		testCase{
			name: "get_all_users_orders",
			input: expected{
				query: `/* name: GetAllUsersOrders :many */
				SELECT u.id user_id, u.first_name, o.price, o.id order_id
				FROM orders o LEFT JOIN users u ON u.id = o.user_id`,
				schema: mockSchema,
			},
			output: &Query{
				SQL: "select u.id as user_id, u.first_name, o.price, o.id as order_id from orders as o left join users as u on u.id = o.user_id",
				Columns: []Column{
					Column{
						&sqlparser.ColumnDefinition{
							Name: sqlparser.NewColIdent("user_id"),
							Type: sqlparser.ColumnType{
								Type:          "int",
								Autoincrement: true,
								NotNull:       false, // beause of the left join
							},
						},
						"users",
					},
					Column{
						&sqlparser.ColumnDefinition{
							Name: sqlparser.NewColIdent("first_name"),
							Type: sqlparser.ColumnType{
								Type:    "varchar",
								NotNull: false, // because of left join
							},
						},
						"users",
					},
					Column{
						&sqlparser.ColumnDefinition{
							Name: sqlparser.NewColIdent("price"),
							Type: sqlparser.ColumnType{
								Type:          "DECIMAL(13, 4)",
								Autoincrement: true,
								NotNull:       true,
							},
						},
						"orders",
					},
					Column{
						&sqlparser.ColumnDefinition{
							Name: sqlparser.NewColIdent("order_id"),
							Type: sqlparser.ColumnType{
								Type:          "int",
								Autoincrement: true,
								NotNull:       true,
							},
						},
						"orders",
					},
				},
				Params:           []*Param{},
				Name:             "GetAllUsersOrders",
				Cmd:              ":many",
				DefaultTableName: "orders",
				SchemaLookup:     mockSchema,
			},
		},
	}

	for _, tt := range tests {
		testCase := tt
		t.Run(tt.name, func(t *testing.T) {
			qs, err := parseContents("example.sql", testCase.input.query, testCase.input.schema, mockSettings)
			if err != nil {
				t.Fatalf("Parsing failed with query: [%v]\n", err)
			}
			if len(qs) != 1 {
				t.Fatalf("Expected one query, not %d", len(qs))
			}
			q := qs[0]
			q.SchemaLookup = nil
			q.Filename = ""
			testCase.output.SchemaLookup = nil
			if diff := cmp.Diff(testCase.output, q); diff != "" {
				t.Errorf("parsed query differs: \n%s", diff)
			}
		})
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

	expected := filterCols(mockSchema.tables["users"], map[string]string{"first_name": "users"})
	if !reflect.DeepEqual(Column{firstNameColDfn, "users"}, expected[0]) {
		t.Errorf("Table schema lookup returned unexpected result")
	}
}

func TestParseInsertUpdate(t *testing.T) {
	type expected struct {
		query  string
		schema *Schema
	}
	type testCase struct {
		name   string
		input  expected
		output *Query
	}

	tests := []testCase{
		testCase{
			name: "insert_users",
			input: expected{
				query:  "/* name: InsertNewUser :exec */\nINSERT INTO users (first_name, last_name) VALUES (?, ?)",
				schema: mockSchema,
			},
			output: &Query{
				SQL:     "insert into users(first_name, last_name) values (?, ?)",
				Columns: nil,
				Params: []*Param{
					&Param{
						OriginalName: ":v1",
						Name:         "first_name",
						Typ:          "string",
					},
					&Param{
						OriginalName: ":v2",
						Name:         "last_name",
						Typ:          "sql.NullString",
					},
				},
				Name:             "InsertNewUser",
				Cmd:              ":exec",
				DefaultTableName: "users",
				SchemaLookup:     mockSchema,
			},
		},
		testCase{
			name: "update_without_where",
			input: expected{
				query:  "/* name: UpdateAllUsers :exec */ update users set first_name = 'Bob'",
				schema: mockSchema,
			},
			output: &Query{
				SQL:              "update users set first_name = 'Bob'",
				Columns:          nil,
				Params:           []*Param{},
				Name:             "UpdateAllUsers",
				Cmd:              ":exec",
				DefaultTableName: "users",
				SchemaLookup:     mockSchema,
			},
		},
		testCase{
			name: "update_users",
			input: expected{
				query:  "/* name: UpdateUserAt :exec */\nUPDATE users SET first_name = ?, last_name = ? WHERE id > ? AND first_name = ? LIMIT 3",
				schema: mockSchema,
			},
			output: &Query{
				SQL:     "update users set first_name = ?, last_name = ? where id > ? and first_name = ? limit 3",
				Columns: nil,
				Params: []*Param{
					&Param{
						OriginalName: ":v1",
						Name:         "first_name",
						Typ:          "string",
					},
					&Param{
						OriginalName: ":v2",
						Name:         "last_name",
						Typ:          "sql.NullString",
					},
					&Param{
						OriginalName: ":v3",
						Name:         "id",
						Typ:          "int",
					},
					&Param{
						OriginalName: ":v4",
						Name:         "first_name",
						Typ:          "string",
					},
				},
				Name:             "UpdateUserAt",
				Cmd:              ":exec",
				DefaultTableName: "users",
				SchemaLookup:     mockSchema,
			},
		},
		testCase{
			name: "insert_users_from_orders",
			input: expected{
				query:  "/* name: InsertUsersFromOrders :exec */\ninsert into users ( first_name ) select user_id from orders where id = ?;",
				schema: mockSchema,
			},
			output: &Query{
				SQL:     "insert into users(first_name) select user_id from orders where id = ?",
				Columns: nil,
				Params: []*Param{
					&Param{
						OriginalName: ":v1",
						Name:         "id",
						Typ:          "int",
					},
				},
				Name:             "InsertUsersFromOrders",
				Cmd:              ":exec",
				DefaultTableName: "users",
				SchemaLookup:     mockSchema,
			},
		},
	}

	for _, tt := range tests {
		testCase := tt
		t.Run(tt.name, func(t *testing.T) {
			qs, err := parseContents("example.sql", testCase.input.query, testCase.input.schema, mockSettings)
			if err != nil {
				t.Fatalf("Parsing failed with query: [%v]\n", err)
			}
			if len(qs) != 1 {
				t.Fatalf("Expected one query, not %d", len(qs))
			}
			q := qs[0]
			testCase.output.SchemaLookup = nil
			q.SchemaLookup = nil
			q.Filename = ""
			if diff := cmp.Diff(testCase.output, q); diff != "" {
				t.Errorf("parsed query differs: \n%s", diff)
			}
		})
	}
}

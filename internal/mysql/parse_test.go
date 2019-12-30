package mysql

import (
	"reflect"
	"testing"

	"github.com/davecgh/go-spew/spew"
	"github.com/kyleconroy/sqlc/internal/dinosql"
	"vitess.io/vitess/go/vt/sqlparser"
)

func init() {
	initMockSchema()
}

const query = `
/* name: GetAllStudents :many */
SELECT school_id, id FROM students WHERE id = :id + ?
`

const create = `
	CREATE TABLE students (
		id int,
		school_id VARCHAR(255),
		school_lat VARCHAR(255),
		PRIMARY KEY (ID)
	);`

const filename = "test.sql"

func TestParseFile(t *testing.T) {
	// s := NewSchema()
	// _, err := parseFile(filename, s)
	// keep(err)
	tree, _ := sqlparser.Parse("SELECT id, first_name FROM users WHERE age < ?")
	p := sqlparser.NewParsedQuery(tree)
	// spew.Dump(p)
	// for k, _ :=
	result := sqlparser.GetBindvars(tree)
	newVars := make(map[string]string)
	for k := range result {
		newVars[k] = "?"
	}
	// spew.Dump(newVars)
	keep(p)
	// p.GenerateQuery(newVars)
	// r, _ := p.MarshalJSON()
	// spew.Dump(string(r))
	// spew.Dump(p.GenerateQuery())
}

var mockSettings = dinosql.GenerateSettings{
	Version: "1",
	Packages: []dinosql.PackageSettings{
		dinosql.PackageSettings{
			Name: "db",
		},
	},
	Overrides: []dinosql.Override{},
}

func TestGenerate(t *testing.T) {
	// t.Skip()
	s := NewSchema()
	result, _ := parseFile(filename, "db", s, mockSettings)
	output, err := dinosql.Generate(result, mockSettings)
	if err != nil {
		t.Errorf("Failed to generate output: %v", err)
	}
	keep(output)
	// for k, v := range output {
	// 	fmt.Println(k)
	// 	fmt.Println(v)
	// 	fmt.Println("")
	// }
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
				// could add more here later if needed
			},
		},
		&sqlparser.ColumnDefinition{
			Name: sqlparser.NewColIdent("last_name"),
			Type: sqlparser.ColumnType{
				Type:    "varchar",
				NotNull: false,
				// could add more here later if needed
			},
		},
		&sqlparser.ColumnDefinition{
			Name: sqlparser.NewColIdent("id"),
			Type: sqlparser.ColumnType{
				Type:          "int",
				NotNull:       true,
				Autoincrement: true,
				// could add more here later if needed
			},
		},
		&sqlparser.ColumnDefinition{
			Name: sqlparser.NewColIdent("age"),
			Type: sqlparser.ColumnType{
				Type:    "int",
				NotNull: true,
				// could add more here later if needed
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
				// could add more here later if needed
			},
		},
		&sqlparser.ColumnDefinition{
			Name: sqlparser.NewColIdent("price"),
			Type: sqlparser.ColumnType{
				Type:          "DECIMAL(13, 4)",
				NotNull:       true,
				Autoincrement: true,
				// could add more here later if needed
			},
		},
		&sqlparser.ColumnDefinition{
			Name: sqlparser.NewColIdent("user_id"),
			Type: sqlparser.ColumnType{
				Type:    "int",
				NotNull: true,
				// could add more here later if needed
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
	query2 := `/* name: GetAll :many */
						SELECT * FROM users;`
	tests := []testCase{
		testCase{
			input: expected{
				query: `/* name: GetNameByID :one */
								SELECT first_name, last_name FROM users WHERE id = ?`,
				schema: mockSchema,
			},
			output: &Query{
				SQL:     `select first_name, last_name from users where id = :v1`,
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
				query:  query2,
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
	}

	for _, testCase := range tests {
		q, err := parseQueryString(testCase.input.query, testCase.input.schema, mockSettings)
		if err != nil {
			t.Errorf("Parsing failed withe query: [%v]\n:schema: %v", query, spew.Sdump(testCase.input.schema))
		}

		err = q.parseNameAndCmd()
		if err != nil {
			t.Errorf("Parsing failed withe query: [%v]\n:schema: %v", query, spew.Sdump(testCase.input.schema))
		}
		if !reflect.DeepEqual(testCase.output, q) {
			t.Errorf("Parsing query returned differently than expected.")
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

func TestParseInsert(t *testing.T) {
	type expected struct {
		query  string
		schema *Schema
	}
	type testCase struct {
		input  expected
		output *Query
	}
	query1 := `/* name: InsertNewUser :exec */
	INSERT INTO users (first_name, last_name) VALUES (?, ?)`
	query2 := `/* name: UpdateUserAt :exec */
	UPDATE users SET first_name = ?, last_name = ? WHERE id > ? AND first_name = ? LIMIT 3`
	tests := []testCase{
		testCase{
			input: expected{
				query:  query1,
				schema: mockSchema,
			},
			output: &Query{
				SQL:     query1,
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
				SQL:     query2,
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
			t.Errorf("Parsing failed with query index: %d: [%v]\n", ix, query)
		}
		if !reflect.DeepEqual(testCase.output, q) {
			t.Errorf("Parsing query returned differently than expected.")
			t.Logf("Expected: %v\nResult: %v\n", spew.Sdump(testCase.output), spew.Sdump(q))
		}
	}
}

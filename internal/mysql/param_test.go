package mysql

import (
	"reflect"
	"testing"

	"github.com/davecgh/go-spew/spew"
	"github.com/kyleconroy/sqlc/internal/dinosql"
	"vitess.io/vitess/go/vt/sqlparser"
)

func TestSelectParamSearcher(t *testing.T) {
	type testCase struct {
		input  string
		output []*Param
	}

	tests := []testCase{
		testCase{
			input: "SELECT first_name, id, last_name FROM users WHERE id < ?",
			output: []*Param{&Param{
				OriginalName: ":v1",
				Name:         "id",
				Typ:          "int",
			},
			},
		},
		testCase{
			input: `SELECT
								users.id,
								users.first_name,
								orders.price
							FROM
								orders
							LEFT JOIN users ON orders.user_id = users.id
							WHERE orders.price > :minPrice`,
			output: []*Param{
				&Param{
					OriginalName: ":minPrice",
					Name:         "minPrice",
					Typ:          "float64",
				},
			},
		},
		testCase{
			input: "SELECT first_name, id, last_name FROM users WHERE id = :targetID",
			output: []*Param{&Param{
				OriginalName: ":targetID",
				Name:         "targetID",
				Typ:          "int",
			},
			},
		},
		testCase{
			input: "SELECT first_name, last_name FROM users WHERE age < :maxAge AND last_name = :inFamily",
			output: []*Param{
				&Param{
					OriginalName: ":maxAge",
					Name:         "maxAge",
					Typ:          "int",
				},
				&Param{
					OriginalName: ":inFamily",
					Name:         "inFamily",
					Typ:          "sql.NullString",
				},
			},
		},
		testCase{
			input: "SELECT first_name, last_name FROM users LIMIT ?",
			output: []*Param{
				&Param{
					OriginalName: ":v1",
					Name:         "limit",
					Typ:          "uint32",
				},
			},
		},
		{
			input: "select first_name, id FROM users LIMIT sqlc.arg(UsersLimit)",
			output: []*Param{
				&Param{
					OriginalName: "sqlc.arg(UsersLimit)",
					Name:         "UsersLimit",
					Typ:          "uint32",
				},
			},
		},
	}
	settings := dinosql.Combine(dinosql.GenerateSettings{}, dinosql.PackageSettings{})
	for _, tCase := range tests {
		generator := PackageGenerator{
			Schema:           mockSchema,
			CombinedSettings: settings,
			packageName:      "db",
		}
		tree, err := sqlparser.Parse(tCase.input)
		if err != nil {
			t.Errorf("Failed to parse input query")
		}
		selectStm, ok := tree.(*sqlparser.Select)

		tableAliasMap, _, err := parseFrom(selectStm.From, false)
		if err != nil {
			t.Errorf("Failed to parse table name alias's: %v", err)
		}

		limitParams, err := generator.paramsInLimitExpr(selectStm.Limit, tableAliasMap)
		if err != nil {
			t.Errorf("Failed to parse limit expression params: %v", err)
		}
		whereParams, err := generator.paramsInWhereExpr(selectStm.Where, tableAliasMap, "users")
		if err != nil {
			t.Errorf("Failed to parse where expression params: %v", err)
		}

		params := append(limitParams, whereParams...)
		if !ok {
			t.Errorf("Test case is not SELECT statement as expected")
		}

		if !reflect.DeepEqual(params, tCase.output) {
			t.Errorf("Param searcher returned unexpected result\nResult: %v\nExpected: %v",
				spew.Sdump(params), spew.Sdump(tCase.output))
		}
	}
}

func TestInsertParamSearcher(t *testing.T) {
	type testCase struct {
		input         string
		output        []*Param
		expectedNames []string
	}

	tests := []testCase{
		testCase{
			input: "/* name: InsertNewUser :exec */\nINSERT INTO users (first_name, last_name) VALUES (?, sqlc.arg(user_last_name))",
			output: []*Param{
				&Param{
					OriginalName: ":v1",
					Name:         "first_name",
					Typ:          "string",
				},
				&Param{
					OriginalName: "sqlc.arg(user_last_name)",
					Name:         "user_last_name",
					Typ:          "sql.NullString",
				},
			},
			expectedNames: []string{"first_name", "user_last_name"},
		},
	}
	settings := dinosql.Combine(dinosql.GenerateSettings{}, dinosql.PackageSettings{})
	for _, tCase := range tests {
		generator := PackageGenerator{
			Schema:           mockSchema,
			CombinedSettings: settings,
			packageName:      "db",
		}
		tree, err := sqlparser.Parse(tCase.input)
		if err != nil {
			t.Errorf("Failed to parse input query")
		}
		insertStm, ok := tree.(*sqlparser.Insert)
		if !ok {
			t.Errorf("Test case is not SELECT statement as expected")
		}
		result, err := generator.parseInsert(insertStm, tCase.input)

		if err != nil {
			t.Errorf("Failed to parse insert statement.")
		}

		if !reflect.DeepEqual(result.Params, tCase.output) {
			t.Errorf("Param searcher returned unexpected result\nResult: %v\nExpected: %v\nQuery: %s",
				spew.Sdump(result.Params), spew.Sdump(tCase.output), tCase.input)
		}
		if len(result.Params) != len(tCase.expectedNames) {
			t.Errorf("Insufficient test cases. Mismatch in length of expected param names and parsed params")
		}
		for ix, p := range result.Params {
			if p.Name != tCase.expectedNames[ix] {
				t.Errorf("Derived param does not match expected output.\nResult: %v\nExpected: %v",
					p.Name, tCase.expectedNames[ix])
			}
		}
	}
}

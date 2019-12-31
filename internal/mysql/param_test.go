package mysql

import (
	"reflect"
	"testing"

	"github.com/davecgh/go-spew/spew"
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
				originalName: ":v1",
				name:         "id",
				typ:          "int",
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
					originalName: ":minPrice",
					name:         "minPrice",
					typ:          "float64",
				},
			},
		},
		testCase{
			input: "SELECT first_name, id, last_name FROM users WHERE id = :targetID",
			output: []*Param{&Param{
				originalName: ":targetID",
				name:         "targetID",
				typ:          "int",
			},
			},
		},
		testCase{
			input: "SELECT first_name, last_name FROM users WHERE age < :maxAge AND last_name = :inFamily",
			output: []*Param{
				&Param{
					originalName: ":maxAge",
					name:         "maxAge",
					typ:          "int",
				},
				&Param{
					originalName: ":inFamily",
					name:         "inFamily",
					typ:          "sql.NullString",
				},
			},
		},
		testCase{
			input: "SELECT first_name, last_name FROM users LIMIT ?",
			output: []*Param{
				&Param{
					originalName: ":v1",
					name:         "limit",
					typ:          "uint32",
				},
			},
		},
	}
	for _, tCase := range tests {
		tree, err := sqlparser.Parse(tCase.input)
		if err != nil {
			t.Errorf("Failed to parse input query")
		}
		selectStm, ok := tree.(*sqlparser.Select)

		limitParams, err := paramsInLimitExpr(selectStm.Limit, mockSchema, mockSettings)
		if err != nil {
			t.Errorf("Failed to parse limit expression params: %v", err)
		}
		whereParams, err := paramsInWhereExpr(selectStm.Where, mockSchema, "users", mockSettings)
		if err != nil {
			t.Errorf("Failed to parse where expression params: %v", err)
		}

		params := append(limitParams, whereParams...)
		if !ok {
			t.Errorf("Test case is not SELECT statement as expected")
		}

		// TODO: get this out of the unit test and/or deprecate defaultTable
		defaultTable := getDefaultTable(&selectStm.From)
		keep(defaultTable)

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
			input: "INSERT INTO users (first_name, last_name) VALUES (?, ?)",
			output: []*Param{
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
			expectedNames: []string{"first_name", "last_name"},
		},
	}
	for _, tCase := range tests {
		tree, err := sqlparser.Parse(tCase.input)
		if err != nil {
			t.Errorf("Failed to parse input query")
		}
		insertStm, ok := tree.(*sqlparser.Insert)
		if !ok {
			t.Errorf("Test case is not SELECT statement as expected")
		}
		result, err := parseInsert(insertStm, tCase.input, mockSchema, mockSettings)
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
			if p.name != tCase.expectedNames[ix] {
				t.Errorf("Derived param does not match expected output.\nResult: %v\nExpected: %v",
					p.name, tCase.expectedNames[ix])
			}
		}
	}
}

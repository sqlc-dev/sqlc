// Temp Test For Debugging. all the test will be finally removed and migrate to the end to end test .
package clickhouse

import (
	"bytes"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/sqlc-dev/sqlc/internal/sql/ast"
)

type testCase struct {
	name     string
	sql      string
	expected []ast.Statement
}

func (c *testCase) evaluate(t *testing.T) {
	parser := Parser{}
	parsed, err := parser.Parse(bytes.NewBuffer([]byte(c.sql)))
	if err != nil {
		t.Error(err)
	}
	diff := cmp.Diff(parsed, c.expected, cmpopts.EquateEmpty())
	if diff != "" {
		t.Error(diff)
	}
}

func TestCreateTable(t *testing.T) {
	case1 := &testCase{
		name: "simple",
		sql:  `CREATE TABLE foo(a String, b Int64);`,
		expected: []ast.Statement{
			{
				Raw: &ast.RawStmt{
					Stmt: &ast.CreateTableStmt{
						IfNotExists: false,
						Name:        &ast.TableName{Name: "foo"},
						Cols: []*ast.ColumnDef{
							{Colname: "a", TypeName: &ast.TypeName{Name: "String"}},
							{Colname: "b", TypeName: &ast.TypeName{Name: "Int64"}},
						},
					},
				},
			},
		},
	}
	t.Run(case1.name, case1.evaluate)
}

func TestSelectQuery(t *testing.T) {
	case1 := []*testCase{
		{
			name: "select all",
			sql:  `SELECT * FROM foo;`,
			expected: []ast.Statement{
				{
					Raw: &ast.RawStmt{
						Stmt: &ast.SelectStmt{
							FromClause: &ast.List{Items: []ast.Node{&ast.TableName{Name: "foo"}}},
							All:        true,
						},
						StmtLen: 17,
					},
				},
			},
		},
		{
			name: "select with where",
			sql:  `SELECT a,b FROM foo WHERE a > 2;`,
			expected: []ast.Statement{
				{
					Raw: &ast.RawStmt{
						Stmt: &ast.SelectStmt{
							FromClause: &ast.List{
								Items: []ast.Node{&ast.TableName{Name: "foo"}}},
							WhereClause: &ast.BoolExpr{Args: &ast.List{Items: []ast.Node{
								&ast.String{Str: "a"},
								&ast.A_Const{Val: &ast.Integer{Ival: 2}},
							}}},
						},
						StmtLen: 31,
					},
				},
			},
		},
	}

	for _, tc := range case1 {
		t.Run(tc.name, tc.evaluate)
	}

}

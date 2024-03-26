package clickhouse

import (
	"bytes"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/sqlc-dev/sqlc/internal/sql/ast"
	"github.com/sqlc-dev/sqlc/internal/sql/catalog"
)

func newCatalog() *catalog.Catalog {
	c := catalog.New("public")
	return c
}

func TestCreateTable(t *testing.T) {

	sql := `CREATE TABLE foo(a String, b Int64);`
	parser := Parser{}
	parsed, err := parser.Parse(bytes.NewBuffer([]byte(sql)))
	if err != nil {
		t.Error(err)
	}

	diff := cmp.Diff(parsed, []ast.Statement{
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
	}, cmpopts.EquateEmpty())
	if diff != "" {
		t.Error(diff)
	}
}

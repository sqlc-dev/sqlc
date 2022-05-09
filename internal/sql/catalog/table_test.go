package catalog

import (
	"testing"

	"github.com/kyleconroy/sqlc/internal/sql/ast"
)

func TestGetTable(t *testing.T) {

	testCases := []struct {
		name      string
		tableName *ast.TableName
		outStub   func(t *testing.T, outSchema *Schema, outTable *Table, err error)
	}{
		{
			name: "Use Catalog Default Schema",
			tableName: &ast.TableName{
				Catalog: "catalog1",
				Schema:  "",
				Name:    "table1",
			},
			outStub: func(t *testing.T, outSchema *Schema, outTable *Table, err error) {

			},
		},
	}

	newCatalog := New("default")

	for i := range testCases {

		tc := testCases[i]

		outSchema, outTable, err := newCatalog.getTable(tc.tableName)

		tc.outStub(t, outSchema, outTable, err)

	}
}

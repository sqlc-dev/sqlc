package duckdb

import (
	"log"
	"strings"

	dw "github.com/sqlc-dev/darkwing/ast"

	"github.com/sqlc-dev/sqlc/internal/debug"
	"github.com/sqlc-dev/sqlc/internal/sql/ast"
)

func todo(n dw.Node) *ast.TODO {
	if debug.Active {
		log.Printf("duckdb.convert: Unknown node type %T\n", n)
	}
	return &ast.TODO{}
}

// identifier normalizes an identifier. DuckDB identifiers are
// case-insensitive (though case-preserving); the catalog matches them
// lowercased.
func identifier(id string) string {
	return strings.ToLower(id)
}

func NewIdentifier(t string) *ast.String {
	return &ast.String{Str: identifier(t)}
}

// schemaName normalizes a schema qualifier. DuckDB's default "main" schema
// maps to the catalog's default namespace, so it is treated as unqualified.
func schemaName(s string) string {
	s = identifier(s)
	if s == "main" {
		return ""
	}
	return s
}

func parseTableName(catalog, schema, name string) *ast.TableName {
	return &ast.TableName{
		Catalog: identifier(catalog),
		Schema:  schemaName(schema),
		Name:    identifier(name),
	}
}

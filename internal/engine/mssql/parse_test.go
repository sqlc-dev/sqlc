package mssql

import (
	"strings"
	"testing"

	"github.com/sqlc-dev/sqlc/internal/sql/ast"
)

// The fmt CLI does not support MSSQL, so exercise its AST printer directly.
func TestLoweredFunctionFormatting(t *testing.T) {
	parser := NewParser()
	for _, name := range []string{"isnull", "coalesce"} {
		t.Run(name, func(t *testing.T) {
			stmts, err := parser.Parse(strings.NewReader("SELECT " + name + "(NULL, 1);"))
			if err != nil {
				t.Fatal(err)
			}
			// This literal-only expression needs no dialect formatting hooks.
			formatted := ast.Format(stmts[0].Raw, nil)
			if !strings.Contains(formatted, name+"(") {
				t.Fatalf("function identity changed: %s", formatted)
			}
			if _, err := parser.Parse(strings.NewReader(formatted)); err != nil {
				t.Fatalf("formatted SQL cannot be parsed: %v", err)
			}
		})
	}
}

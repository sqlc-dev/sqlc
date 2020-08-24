package postgresql

import (
	"strings"
	"testing"

	"github.com/kyleconroy/sqlc/internal/sql/ast"
	"github.com/kyleconroy/sqlc/internal/sql/astutils"

	"github.com/google/go-cmp/cmp"
)

func TestApply(t *testing.T) {
	p := NewParser()

	input, err := p.Parse(strings.NewReader("SELECT sqlc.arg(name)"))
	if err != nil {
		t.Fatal(err)
	}
	output, err := p.Parse(strings.NewReader("SELECT $1"))
	if err != nil {
		t.Fatal(err)
	}

	expect := &output[0]
	actual := astutils.Apply(&input[0], func(cr *astutils.Cursor) bool {
		fun, ok := cr.Node().(*ast.FuncCall)
		if !ok {
			return true
		}
		if astutils.Join(fun.Funcname, ".") == "sqlc.arg" {
			cr.Replace(&ast.ParamRef{
				Number:   1,
				Location: fun.Location,
			})
			return false
		}
		return true
	}, nil)

	if diff := cmp.Diff(expect, actual); diff != "" {
		t.Errorf("rewrite mismatch:\n%s", diff)
	}
}

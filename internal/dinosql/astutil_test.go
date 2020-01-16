package dinosql

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	pg "github.com/lfittl/pg_query_go"
	nodes "github.com/lfittl/pg_query_go/nodes"
)

func TestApply(t *testing.T) {
	input, err := pg.Parse("SELECT sqlc.arg(name)")
	if err != nil {
		t.Fatal(err)
	}
	output, err := pg.Parse("SELECT $1")
	if err != nil {
		t.Fatal(err)
	}

	// spew.Dump(input.Statements[0])

	expect := output.Statements[0]
	actual := Apply(input.Statements[0], func(cr *Cursor) bool {
		fun, ok := cr.Node().(nodes.FuncCall)
		if !ok {
			return true
		}
		if join(fun.Funcname, ".") == "sqlc.arg" {
			cr.Replace(nodes.ParamRef{Number: 1})
			return false
		}
		return true
	}, nil)

	if diff := cmp.Diff(expect, actual); diff != "" {
		t.Errorf("rewrite mismatch:\n%s", diff)
	}
}

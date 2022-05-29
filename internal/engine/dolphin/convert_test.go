package dolphin_test

import (
	"strings"
	"testing"

	"github.com/kyleconroy/sqlc/internal/engine/dolphin"
	"github.com/kyleconroy/sqlc/internal/sql/ast"
	"github.com/kyleconroy/sqlc/internal/sql/astutils"
)

func Test_TupleComparison(t *testing.T) {
	p := dolphin.NewParser()
	stmts, err := p.Parse(strings.NewReader("SELECT * WHERE (a, b) > (?, ?)"))
	if err != nil {
		t.Fatal(err)
	}

	if l := len(stmts); l != 1 {
		t.Fatalf("expected 1 statement, got %d", l)
	}

	// Right now all this test does is make sure we noticed the two ParamRefs.
	// This ensures that the Go code is generated correctly.
	e := &refExtractor{}
	astutils.Walk(e, stmts[0].Raw.Stmt)
	if l := len(e.params); l != 2 {
		t.Fatalf("expected to extract 2 params, got %d", l)
	}
}

type refExtractor struct {
	params []*ast.ParamRef
}

func (e *refExtractor) Visit(n ast.Node) astutils.Visitor {
	switch t := n.(type) {
	case *ast.ParamRef:
		e.params = append(e.params, t)
	}
	return e
}

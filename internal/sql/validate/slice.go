package validate

import (
	"fmt"

	"github.com/sqlc-dev/sqlc/internal/sql/ast"
	"github.com/sqlc-dev/sqlc/internal/sql/astutils"
	"github.com/sqlc-dev/sqlc/internal/sql/sqlerr"
)

type sliceVisitor struct {
	// slices maps the location of a placeholder in the query text to the name
	// of the sqlc.slice() parameter it was rewritten from.
	slices map[int]string
	err    error
}

func (v *sliceVisitor) Visit(node ast.Node) astutils.Visitor {
	if v.err != nil {
		return nil
	}

	in, ok := node.(*ast.In)
	if !ok {
		return v
	}

	// A slice expands to a comma-separated list at query time, so it has to be
	// the only element of the IN expression:
	//    id IN (sqlc.slice("ids"))       -- GOOD
	//    id in (0, 1, sqlc.slice("ids")) -- BAD
	if len(in.List) <= 1 {
		return v
	}

	for _, n := range in.List {
		ref, ok := n.(*ast.ParamRef)
		if !ok {
			continue
		}
		name, ok := v.slices[ref.Location]
		if !ok {
			continue
		}

		inExpr := "..."
		if col, ok := in.Expr.(*ast.ColumnRef); ok {
			inExpr = col.Name
		}
		v.err = &sqlerr.Error{
			Message:  fmt.Sprintf("expected '%s IN' expr to consist only of sqlc.slice(%s); eg ", inExpr, name),
			Location: ref.Location,
		}
		return nil
	}

	return v
}

// Slice reports whether every sqlc.slice() parameter is the only element of the
// IN expression that contains it. The slices map is keyed by the location the
// preprocessor recorded for each rewritten slice placeholder.
func Slice(n ast.Node, slices map[int]string) error {
	if len(slices) == 0 {
		return nil
	}
	visitor := sliceVisitor{slices: slices}
	astutils.Walk(&visitor, n)
	return visitor.err
}

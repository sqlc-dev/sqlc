package validate

import (
	"fmt"

	"github.com/kyleconroy/sqlc/internal/sql/ast"
	"github.com/kyleconroy/sqlc/internal/sql/astutils"
	"github.com/kyleconroy/sqlc/internal/sql/sqlerr"
)

func ParamRef(n ast.Node) (map[int]bool, error) {
	var allrefs []*ast.ParamRef

	// Find all parameter references
	astutils.Walk(astutils.VisitorFunc(func(node ast.Node) {
		switch n := node.(type) {
		case *ast.ParamRef:
			allrefs = append(allrefs, n)
		}
	}), n)

	seen := map[int]bool{}
	for _, r := range allrefs {
		if r.Number > 0 {
			seen[r.Number] = true
		}
	}
	for i := 1; i <= len(seen); i += 1 {
		if _, ok := seen[i]; !ok {
			return nil, &sqlerr.Error{
				Code:    "42P18",
				Message: fmt.Sprintf("could not determine data type of parameter $%d", i),
			}
		}
	}
	return seen, nil
}

package validate

import (
	"fmt"

	"github.com/kyleconroy/sqlc/internal/sql/ast"
	"github.com/kyleconroy/sqlc/internal/sql/astutils"
	"github.com/kyleconroy/sqlc/internal/sql/sqlerr"
)

func ParamRef(n ast.Node) (int, error) {
	var allrefs []*ast.ParamRef

	// Find all parameter references
	astutils.Walk(astutils.VisitorFunc(func(node ast.Node) {
		switch n := node.(type) {
		case *ast.ParamRef:
			allrefs = append(allrefs, n)
		}
	}), n)

	seen := map[int]struct{}{}
	for _, r := range allrefs {
		seen[r.Number] = struct{}{}
	}
	var max int
	for i := 1; i <= len(seen); i += 1 {
		if i > max {
			max = i
		}
		if _, ok := seen[i]; !ok {
			return 0, &sqlerr.Error{
				Code:    "42P18",
				Message: fmt.Sprintf("could not determine data type of parameter $%d", i),
			}
		}
	}
	return max, nil
}

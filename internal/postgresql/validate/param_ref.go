package validate

import (
	"fmt"

	nodes "github.com/lfittl/pg_query_go/nodes"

	"github.com/kyleconroy/sqlc/internal/pg"
	"github.com/kyleconroy/sqlc/internal/postgresql/ast"
)

func ParamRef(n nodes.Node) error {
	var allrefs []nodes.ParamRef

	// Find all parameter references
	ast.Walk(ast.VisitorFunc(func(node nodes.Node) {
		switch n := node.(type) {
		case nodes.ParamRef:
			allrefs = append(allrefs, n)
		}
	}), n)

	seen := map[int]struct{}{}
	for _, r := range allrefs {
		seen[r.Number] = struct{}{}
	}

	for i := 1; i <= len(seen); i += 1 {
		if _, ok := seen[i]; !ok {
			return pg.Error{
				Code:    "42P18",
				Message: fmt.Sprintf("could not determine data type of parameter $%d", i),
			}
		}
	}
	return nil
}

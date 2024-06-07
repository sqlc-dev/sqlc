package validate

import (
	"errors"
	"fmt"

	"github.com/sqlc-dev/sqlc/internal/sql/ast"
	"github.com/sqlc-dev/sqlc/internal/sql/astutils"
	"github.com/sqlc-dev/sqlc/internal/sql/sqlerr"
)

func ParamRef(n ast.Node) (map[int]bool, bool, error) {
	seen := map[int]bool{}
	var dollar, nodollar bool
	// Find all parameter references
	astutils.Walk(astutils.VisitorFunc(func(node ast.Node) {
		switch n := node.(type) {
		case *ast.ParamRef:
			if n.Dollar {
				dollar = true
			} else {
				nodollar = true
			}
			if n.Number > 0 {
				seen[n.Number] = true
			}
		}
	}), n)
	if dollar && nodollar {
		return nil, false, errors.New("can not mix $1 format with ? format")
	}

	for i := 1; i <= len(seen); i += 1 {
		if _, ok := seen[i]; !ok {
			return seen, !nodollar, &sqlerr.Error{
				Code:    "42P18",
				Message: fmt.Sprintf("could not determine data type of parameter $%d", i),
			}
		}
	}
	return seen, !nodollar, nil
}

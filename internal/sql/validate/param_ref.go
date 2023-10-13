package validate

import (
	"errors"
	"fmt"

	"github.com/sqlc-dev/sqlc/internal/sql/ast"
	"github.com/sqlc-dev/sqlc/internal/sql/astutils"
	"github.com/sqlc-dev/sqlc/internal/sql/sqlerr"
)

func ParamRef(n ast.Node) (map[int]bool, bool, error) {
	var allrefs []*ast.ParamRef
	var dollar bool
	var nodollar bool
	// Find all parameter references
	astutils.Walk(astutils.VisitorFunc(func(node ast.Node) {
		switch n := node.(type) {
		case *ast.ParamRef:
			ref := node.(*ast.ParamRef)
			if ref.Dollar {
				dollar = true
			} else {
				nodollar = true
			}
			allrefs = append(allrefs, n)
		}
	}), n)
	if dollar && nodollar {
		return nil, false, errors.New("can not mix $1 format with ? format")
	}

	seen := map[int]bool{}
	for _, r := range allrefs {
		if r.Number > 0 {
			seen[r.Number] = true
		}
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

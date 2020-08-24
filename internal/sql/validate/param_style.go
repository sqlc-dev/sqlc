package validate

import (
	"github.com/kyleconroy/sqlc/internal/sql/ast"
	"github.com/kyleconroy/sqlc/internal/sql/astutils"
	"github.com/kyleconroy/sqlc/internal/sql/named"
	"github.com/kyleconroy/sqlc/internal/sql/sqlerr"
)

// A query can use one (and only one) of the following formats:
// - positional parameters           $1
// - named parameter operator        @param
// - named parameter function calls  sqlc.arg(param)
func ParamStyle(n ast.Node) error {
	positional := astutils.Search(n, func(node ast.Node) bool {
		_, ok := node.(*ast.ParamRef)
		return ok
	})
	namedFunc := astutils.Search(n, named.IsParamFunc)
	namedSign := astutils.Search(n, named.IsParamSign)
	for _, check := range []bool{
		len(positional.Items) > 0 && len(namedSign.Items)+len(namedFunc.Items) > 0,
		len(namedFunc.Items) > 0 && len(namedSign.Items)+len(positional.Items) > 0,
		len(namedSign.Items) > 0 && len(positional.Items)+len(namedFunc.Items) > 0,
	} {
		if check {
			return &sqlerr.Error{
				Code:    "", // TODO: Pick a new error code
				Message: "query mixes positional parameters ($1) and named parameters (sqlc.arg or @arg)",
			}
		}
	}
	return nil
}

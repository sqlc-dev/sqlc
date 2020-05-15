package validate

import (
	nodes "github.com/lfittl/pg_query_go/nodes"

	"github.com/kyleconroy/sqlc/internal/pg"
	"github.com/kyleconroy/sqlc/internal/postgresql"
	"github.com/kyleconroy/sqlc/internal/postgresql/ast"
)

// A query can use one (and only one) of the following formats:
// - positional parameters           $1
// - named parameter operator        @param
// - named parameter function calls  sqlc.arg(param)
func ParamStyle(n nodes.Node) error {
	positional := ast.Search(n, func(node nodes.Node) bool {
		_, ok := node.(nodes.ParamRef)
		return ok
	})
	namedFunc := ast.Search(n, postgresql.IsNamedParamFunc)
	namedSign := ast.Search(n, postgresql.IsNamedParamSign)
	for _, check := range []bool{
		len(positional.Items) > 0 && len(namedSign.Items)+len(namedFunc.Items) > 0,
		len(namedFunc.Items) > 0 && len(namedSign.Items)+len(positional.Items) > 0,
		len(namedSign.Items) > 0 && len(positional.Items)+len(namedFunc.Items) > 0,
	} {
		if check {
			return pg.Error{
				Code:    "", // TODO: Pick a new error code
				Message: "query mixes positional parameters ($1) and named parameters (sqlc.arg or @arg)",
			}
		}
	}
	return nil
}

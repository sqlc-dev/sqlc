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
	namedFunc := astutils.Search(n, named.IsParamFunc)
	for _, f := range namedFunc.Items {
		fc, ok := f.(*ast.FuncCall)
		if ok {
			/*
				if len(fc.Args.Items) != 1 {
					return &sqlerr.Error{
						Code:    "", // TODO: Pick a new error code
						Message: "Wrong number of arguments to sqlc.arg()",
					}
				}
			*/
			switch fc.Args.Items[0].(type) {
			case *ast.FuncCall:
				return &sqlerr.Error{
					Code:    "", // TODO: Pick a new error code
					Message: "expected parameter to sqlc.arg to be string or reference; got *ast.FuncCall",
				}
			case *ast.ParamRef:
				return &sqlerr.Error{
					Code:    "", // TODO: Pick a new error code
					Message: "query mixes positional parameters ($1) and named parameters (sqlc.arg or @arg)",
				}
			case *ast.A_Const, *ast.ColumnRef:
			default:
				return &sqlerr.Error{
					Code:    "", // TODO: Pick a new error code
					Message: "Invalid argument to sqlc.arg()",
				}

			}
		}
	}
	return nil
}

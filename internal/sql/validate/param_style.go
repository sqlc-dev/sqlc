package validate

import (
	"fmt"
	"strconv"

	"github.com/sqlc-dev/sqlc/internal/sql/ast"
	"github.com/sqlc-dev/sqlc/internal/sql/astutils"
	"github.com/sqlc-dev/sqlc/internal/sql/sqlerr"
)

type sqlcFuncVisitor struct {
	err error
}

func (v *sqlcFuncVisitor) Visit(node ast.Node) astutils.Visitor {
	if v.err != nil {
		return nil
	}

	call, ok := node.(*ast.FuncCall)
	if !ok {
		return v
	}
	fn := call.Func
	if fn == nil {
		return v
	}

	// Custom validation for sqlc.arg, sqlc.narg and sqlc.slice
	// TODO: Replace this once type-checking is implemented
	if fn.Schema == "sqlc" {
		if !(fn.Name == "arg" || fn.Name == "narg" || fn.Name == "slice" || fn.Name == "embed" || fn.Name == "sort") {
			v.err = sqlerr.FunctionNotFound("sqlc." + fn.Name)
			return nil
		}

		minArgs := 1
		maxArgs := 1
		if fn.Name == "sort" {
			maxArgs = 4
		}
		if len(call.Args.Items) > maxArgs || len(call.Args.Items) < minArgs {
			expectedNumArgs := strconv.Itoa(minArgs)
			if maxArgs != minArgs {
				expectedNumArgs += "-" + strconv.Itoa(maxArgs)
			}
			expectedNumArgs += " parameter"
			if maxArgs != minArgs {
				expectedNumArgs += "s"
			}
			v.err = &sqlerr.Error{
				Message:  fmt.Sprintf("expected %s to sqlc.%s; got %d", expectedNumArgs, fn.Name, len(call.Args.Items)),
				Location: call.Pos(),
			}
			return nil
		}

		switch n := call.Args.Items[0].(type) {
		case *ast.A_Const:
		case *ast.ColumnRef:
		default:
			v.err = &sqlerr.Error{
				Message:  fmt.Sprintf("expected parameter to sqlc.%s to be string or reference; got %T", fn.Name, n),
				Location: call.Pos(),
			}
			return nil
		}

		// If we have sqlc.arg or sqlc.narg, there is no need to resolve the function call.
		// It won't resolve anyway, sinc it is not a real function.
		return nil
	}

	return nil
}

func SqlcFunctions(n ast.Node) error {
	visitor := sqlcFuncVisitor{}
	astutils.Walk(&visitor, n)
	return visitor.err
}

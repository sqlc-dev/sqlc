package validate

import (
	"fmt"
	"strings"

	"github.com/kyleconroy/sqlc/internal/sql/ast"
	"github.com/kyleconroy/sqlc/internal/sql/astutils"
	"github.com/kyleconroy/sqlc/internal/sql/catalog"
	"github.com/kyleconroy/sqlc/internal/sql/sqlerr"
)

type funcCallVisitor struct {
	catalog *catalog.Catalog
	err     error
}

func (v *funcCallVisitor) Visit(node ast.Node) astutils.Visitor {
	if v.err != nil {
		return nil
	}

	call, ok := node.(*ast.FuncCall)
	if !ok {
		return v
	}
	if call.Func == nil {
		return v
	}

	// Do not validate unknown functions
	funs, err := v.catalog.ListFuncsByName(call.Func)
	if err != nil || len(funs) == 0 {
		return v
	}

	var args int
	if call.Args != nil {
		args = len(call.Args.Items)
	}
	for _, fun := range funs {
		if len(fun.InArgs()) == args {
			return v
		}
	}

	var sig []string
	for range call.Args.Items {
		sig = append(sig, "unknown")
	}

	v.err = &sqlerr.Error{
		Code:     "42883",
		Message:  fmt.Sprintf("function %s(%s) does not exist", call.Func.Name, strings.Join(sig, ", ")),
		Location: call.Pos(),
		// Hint: "No function matches the given name and argument types. You might need to add explicit type casts.",
	}

	return nil
}

func FuncCall(c *catalog.Catalog, n ast.Node) error {
	visitor := funcCallVisitor{catalog: c}
	astutils.Walk(&visitor, n)
	return visitor.err
}

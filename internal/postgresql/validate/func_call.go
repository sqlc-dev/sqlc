package validate

import (
	"fmt"
	"strings"

	nodes "github.com/lfittl/pg_query_go/nodes"

	"github.com/kyleconroy/sqlc/internal/catalog"
	"github.com/kyleconroy/sqlc/internal/pg"
	"github.com/kyleconroy/sqlc/internal/postgresql/ast"
)

type funcCallVisitor struct {
	catalog *pg.Catalog
	err     error
}

func (v *funcCallVisitor) Visit(node nodes.Node) ast.Visitor {
	if v.err != nil {
		return nil
	}

	funcCall, ok := node.(nodes.FuncCall)
	if !ok {
		return v
	}

	fqn, err := catalog.ParseList(funcCall.Funcname)
	if err != nil {
		v.err = err
		return v
	}

	// Do not validate unknown functions
	funs, err := v.catalog.LookupFunctions(fqn)
	if err != nil {
		return v
	}

	args := len(funcCall.Args.Items)
	for _, fun := range funs {
		arity := fun.ArgN
		if fun.Arguments != nil {
			arity = len(fun.Arguments)
		}
		if arity == args {
			return v
		}
	}

	var sig []string
	for range funcCall.Args.Items {
		sig = append(sig, "unknown")
	}

	v.err = pg.Error{
		Code:     "42883",
		Message:  fmt.Sprintf("function %s(%s) does not exist", fqn.Rel, strings.Join(sig, ", ")),
		Hint:     "No function matches the given name and argument types. You might need to add explicit type casts.",
		Location: funcCall.Location,
	}

	return nil
}

func FuncCall(c *pg.Catalog, n nodes.Node) error {
	visitor := funcCallVisitor{catalog: c}
	ast.Walk(&visitor, n)
	return visitor.err
}

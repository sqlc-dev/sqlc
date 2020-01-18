package dinosql

import (
	"fmt"
	"strings"

	nodes "github.com/lfittl/pg_query_go/nodes"

	"github.com/kyleconroy/sqlc/internal/catalog"
	"github.com/kyleconroy/sqlc/internal/pg"
	"github.com/kyleconroy/sqlc/internal/postgresql/ast"
)

func validateParamRef(n nodes.Node) error {
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

func validateFuncCall(c *pg.Catalog, n nodes.Node) error {
	visitor := funcCallVisitor{catalog: c}
	ast.Walk(&visitor, n)
	return visitor.err
}

func validateInsertStmt(stmt nodes.InsertStmt) error {
	sel, ok := stmt.SelectStmt.(nodes.SelectStmt)
	if !ok {
		return nil
	}
	if len(sel.ValuesLists) != 1 {
		return nil
	}

	colsLen := len(stmt.Cols.Items)
	valsLen := len(sel.ValuesLists[0])
	switch {
	case colsLen > valsLen:
		return pg.Error{
			Code:    "42601",
			Message: "INSERT has more target columns than expressions",
		}
	case colsLen < valsLen:
		return pg.Error{
			Code:    "42601",
			Message: "INSERT has more expressions than target columns",
		}
	}
	return nil
}

// A query can use one (and only one) of the following formats:
// - positional parameters           $1
// - named parameter operator        @param
// - named parameter function calls  sqlc.arg(param)
func validateParamStyle(n nodes.Node) error {
	positional := search(n, func(node nodes.Node) bool {
		_, ok := node.(nodes.ParamRef)
		return ok
	})
	namedFunc := search(n, isNamedParamFunc)
	namedSign := search(n, isNamedParamSign)
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

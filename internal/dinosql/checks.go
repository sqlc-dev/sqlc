package dinosql

import (
	"fmt"
	"strings"

	"github.com/kyleconroy/dinosql/internal/catalog"
	"github.com/kyleconroy/dinosql/internal/pg"
	nodes "github.com/lfittl/pg_query_go/nodes"
)

type Error struct {
	Message string
	Code    string
	Hint    string
}

func (e Error) Error() string {
	return e.Message
}

func validateParamRef(n nodes.Node) error {
	var allrefs []nodes.ParamRef

	// Find all parameter references
	Walk(VisitorFunc(func(node nodes.Node) {
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
			return Error{
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

func (v *funcCallVisitor) Visit(node nodes.Node) Visitor {
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
		if fun.ArgN == args {
			return v
		}
	}

	var sig []string
	for _, _ = range funcCall.Args.Items {
		sig = append(sig, "unknown")
	}

	v.err = Error{
		Code:    "42883",
		Message: fmt.Sprintf("function %s(%s) does not exist", fqn.Rel, strings.Join(sig, ", ")),
		Hint:    "No function matches the given name and argument types. You might need to add explicit type casts.",
	}

	return nil
}

func validateFuncCall(c *pg.Catalog, n nodes.Node) error {
	visitor := funcCallVisitor{catalog: c}
	Walk(&visitor, n)
	return visitor.err
}

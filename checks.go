package dinosql

import (
	"fmt"
	"strings"

	"github.com/kyleconroy/dinosql/postgres"
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
	err error
}

func (v *funcCallVisitor) Visit(node nodes.Node) Visitor {
	if v.err != nil {
		return nil
	}

	funcCall, ok := node.(nodes.FuncCall)
	if !ok {
		return v
	}

	// Do not validate unknown functions
	name := join(funcCall.Funcname, ".")
	if _, ok := postgres.Functions[name]; !ok {
		return v
	}

	args := len(funcCall.Args.Items)
	if _, ok := postgres.Functions[name][args]; ok {
		return v
	}

	var sig []string
	for _, _ = range funcCall.Args.Items {
		sig = append(sig, "unknown")
	}

	v.err = Error{
		Code:    "42883",
		Message: fmt.Sprintf("function %s(%s) does not exist", name, strings.Join(sig, ", ")),
		Hint:    "No function matches the given name and argument types. You might need to add explicit type casts.",
	}

	return nil
}

func validateFuncCall(n nodes.Node) error {
	visitor := funcCallVisitor{}
	Walk(&visitor, n)
	return visitor.err
}

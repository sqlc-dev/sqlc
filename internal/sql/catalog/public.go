package catalog

import (
	"fmt"
	"strings"

	"github.com/sqlc-dev/sqlc/internal/sql/ast"
	"github.com/sqlc-dev/sqlc/internal/sql/sqlerr"
)

func (c *Catalog) schemasToSearch(ns string) []string {
	if ns == "" {
		ns = c.DefaultSchema
	}
	return append(c.SearchPath, ns)
}

func (c *Catalog) ListFuncsByName(rel *ast.FuncName) ([]Function, error) {
	var funcs []Function
	lowered := strings.ToLower(rel.Name)
	for _, ns := range c.schemasToSearch(rel.Schema) {
		s, err := c.getSchema(ns)
		if err != nil {
			return nil, err
		}
		for i := range s.Funcs {
			if strings.ToLower(s.Funcs[i].Name) == lowered {
				funcs = append(funcs, *s.Funcs[i])
			}
		}
	}
	return funcs, nil
}

func (c *Catalog) ResolveFuncCall(call *ast.FuncCall) (*Function, error) {
	// Do not validate unknown functions
	funs, err := c.ListFuncsByName(call.Func)
	if err != nil || len(funs) == 0 {
		return nil, sqlerr.FunctionNotFound(call.Func.Name)
	}

	// https://www.postgresql.org/docs/current/sql-syntax-calling-funcs.html
	var positional []ast.Node
	var named []*ast.NamedArgExpr

	if call.Args != nil {
		for _, arg := range call.Args.Items {
			if narg, ok := arg.(*ast.NamedArgExpr); ok {
				named = append(named, narg)
			} else {
				// The mixed notation combines positional and named notation.
				// However, as already mentioned, named arguments cannot precede
				// positional arguments.
				if len(named) > 0 {
					return nil, &sqlerr.Error{
						Code:     "",
						Message:  "positional argument cannot follow named argument",
						Location: call.Pos(),
					}
				}
				positional = append(positional, arg)
			}
		}
	}

	for _, fun := range funs {
		args := fun.InArgs()
		var defaults int
		var variadic bool
		known := map[string]struct{}{}
		for _, arg := range args {
			if arg.HasDefault {
				defaults += 1
			}
			if arg.Mode == ast.FuncParamVariadic {
				variadic = true
				defaults += 1
			}
			if arg.Name != "" {
				known[arg.Name] = struct{}{}
			}
		}

		if variadic {
			if (len(named) + len(positional)) < (len(args) - defaults) {
				continue
			}
		} else {
			if (len(named) + len(positional)) > len(args) {
				continue
			}
			if (len(named) + len(positional)) < (len(args) - defaults) {
				continue
			}
		}

		// Validate that the provided named arguments exist in the function
		var unknownArgName bool
		for _, expr := range named {
			if expr.Name != nil {
				if _, found := known[*expr.Name]; !found {
					unknownArgName = true
				}
			}
		}
		if unknownArgName {
			continue
		}

		return &fun, nil
	}

	var sig []string
	for range call.Args.Items {
		sig = append(sig, "unknown")
	}

	return nil, &sqlerr.Error{
		Code:     "42883",
		Message:  fmt.Sprintf("function %s(%s) does not exist", call.Func.Name, strings.Join(sig, ", ")),
		Location: call.Pos(),
		// Hint: "No function matches the given name and argument types. You might need to add explicit type casts.",
	}
}

func (c *Catalog) GetTable(rel *ast.TableName) (Table, error) {
	_, table, err := c.getTable(rel)
	if table == nil {
		return Table{}, err
	} else {
		return *table, err
	}
}

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
		// Separate input and output args from the function signature
		inArgs := fun.InArgs()
		outArgs := fun.OutArgs()

		// Build known argument names from all parameters (IN/OUT/INOUT/etc.)
		known := map[string]struct{}{}
		for _, a := range fun.Args {
			if a.Name != "" {
				known[a.Name] = struct{}{}
			}
		}

		// Count defaults and whether the last IN arg is variadic
		var defaultsIn int
		var variadic bool
		for _, arg := range inArgs {
			if arg.HasDefault {
				defaultsIn += 1
			}
			if arg.Mode == ast.FuncParamVariadic {
				variadic = true
				// Treat the variadic parameter like having a default for count checks
				defaultsIn += 1
			}
		}

		// Tally named arguments provided by the call: which refer to IN vs OUT names
		var namedIn, namedOut int
		var unknownArgName bool
		for _, expr := range named {
			if expr.Name != nil {
				name := *expr.Name
				if _, ok := known[name]; !ok {
					unknownArgName = true
					continue
				}
				// Classify whether the provided named arg matches an IN or OUT param
				var isIn bool
				for _, a := range inArgs {
					if a.Name == name {
						isIn = true
						break
					}
				}
				if isIn {
					namedIn += 1
				} else {
					// If not IN, treat it as an OUT placeholder/name
					namedOut += 1
				}
			}
		}
		if unknownArgName {
			// Provided a named argument that doesn't exist in the signature
			continue
		}

		// Positional arguments always come first (we validated above that
		// positional cannot follow named). They fill IN parameters first; any
		// excess positional arguments are treated as placeholders for OUT params
		var posFillIn = len(positional)
		if posFillIn > len(inArgs) {
			posFillIn = len(inArgs)
		}
		// Count how many IN arguments are provided (positional for IN + named for IN)
		inProvided := posFillIn + namedIn

		// Validate IN argument counts against the signature considering defaults/variadic
		if variadic {
			if inProvided < (len(inArgs) - defaultsIn) {
				continue
			}
		} else {
			if inProvided > len(inArgs) {
				continue
			}
			if inProvided < (len(inArgs) - defaultsIn) {
				continue
			}
		}

		// Validate OUT placeholders. These are only valid in procedure calls.
		// For normal function invocation, callers cannot pass values for OUT params.
		posOut := 0
		if len(positional) > len(inArgs) {
			posOut = len(positional) - len(inArgs)
		}
		outProvided := posOut + namedOut
		if fun.ReturnType == nil {
			// Procedure: allow passing placeholders for OUT params, but not more than available
			if outProvided > len(outArgs) {
				continue
			}
		} else {
			// Function: do not allow any OUT placeholders
			if outProvided > 0 {
				continue
			}
		}

		// All checks passed for this candidate
		return &fun, nil
	}

	var sig []string
	for range call.Args.Items {
		sig = append(sig, "unknown")
	}

	return nil, &sqlerr.Error{
		Code:     "42883",
		Message:  fmt.Sprintf("CODE 42883: function %s(%s) does not exist", call.Func.Name, strings.Join(sig, ", ")),
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

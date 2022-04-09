package rewrite

import (
	"fmt"

	"github.com/kyleconroy/sqlc/internal/config"
	"github.com/kyleconroy/sqlc/internal/source"
	"github.com/kyleconroy/sqlc/internal/sql/ast"
	"github.com/kyleconroy/sqlc/internal/sql/astutils"
	"github.com/kyleconroy/sqlc/internal/sql/named"
)

// Given an AST node, return the string representation of names
func flatten(root ast.Node) (string, bool) {
	sw := &stringWalker{}
	astutils.Walk(sw, root)
	return sw.String, sw.IsConst
}

type stringWalker struct {
	String  string
	IsConst bool
}

func (s *stringWalker) Visit(node ast.Node) astutils.Visitor {
	if _, ok := node.(*ast.A_Const); ok {
		s.IsConst = true
	}
	if n, ok := node.(*ast.String); ok {
		s.String += n.Str
	}
	return s
}

func isNamedParamSignCast(node ast.Node) bool {
	expr, ok := node.(*ast.A_Expr)
	if !ok {
		return false
	}
	_, cast := expr.Rexpr.(*ast.TypeCast)
	return astutils.Join(expr.Name, ".") == "@" && cast
}

type namedParameter struct {
	param named.Param
	locs  []int
}

// Add a new instance of this parameter of the same name
func (n namedParameter) AddInstance(loc int, p named.Param) namedParameter {
	param := named.Combine(n.param, p)
	locs := append(n.locs, loc)
	return namedParameter{
		param: param,
		locs:  locs,
	}
}

// paramFromName takes a user-defined parameter name and builds the appropiate parameter
func paramFromName(name string) named.Param {
	return named.NewUnspecifiedParam(name)
}

func NamedParameters(engine config.Engine, raw *ast.RawStmt, numbs map[int]bool, dollar bool) (*ast.RawStmt, map[int]named.Param, []source.Edit) {
	foundFunc := astutils.Search(raw, named.IsParamFunc)
	foundSign := astutils.Search(raw, named.IsParamSign)
	if len(foundFunc.Items)+len(foundSign.Items) == 0 {
		return raw, map[int]named.Param{}, nil
	}

	hasNamedParameterSupport := engine != config.EngineMySQL

	args := map[string]namedParameter{}
	argn := 0
	var edits []source.Edit
	node := astutils.Apply(raw, func(cr *astutils.Cursor) bool {
		node := cr.Node()
		switch {
		case named.IsParamFunc(node):
			fun := node.(*ast.FuncCall)
			paramName, isConst := flatten(fun.Args)
			param := paramFromName(paramName)

			if namedP, ok := args[param.Name()]; ok && hasNamedParameterSupport {
				cr.Replace(&ast.ParamRef{
					Number:   namedP.locs[0],
					Location: fun.Location,
				})
			} else {
				// Find the arg number that has not yet been used
				argn++
				for numbs[argn] {
					argn++
				}

				args[param.Name()] = args[param.Name()].AddInstance(argn, param)
				cr.Replace(&ast.ParamRef{
					Number:   argn,
					Location: fun.Location,
				})
			}
			// TODO: This code assumes that sqlc.arg(name) is on a single line
			var old, replace string
			if isConst {
				old = fmt.Sprintf("sqlc.arg('%s')", paramName)
			} else {
				old = fmt.Sprintf("sqlc.arg(%s)", paramName)
			}
			if engine == config.EngineMySQL || !dollar {
				replace = "?"
			} else {
				replace = fmt.Sprintf("$%d", args[param.Name()].locs[0])
			}
			edits = append(edits, source.Edit{
				Location: fun.Location - raw.StmtLocation,
				Old:      old,
				New:      replace,
			})
			return false

		case isNamedParamSignCast(node):
			expr := node.(*ast.A_Expr)
			cast := expr.Rexpr.(*ast.TypeCast)
			paramName, _ := flatten(cast.Arg)
			param := paramFromName(paramName)
			if p, ok := args[param.Name()]; ok {
				cast.Arg = &ast.ParamRef{
					Number:   p.locs[0],
					Location: expr.Location,
				}
				cr.Replace(cast)
			} else {
				argn++
				for numbs[argn] {
					argn++
				}

				args[param.Name()] = args[param.Name()].AddInstance(argn, param)
				cast.Arg = &ast.ParamRef{
					Number:   argn,
					Location: expr.Location,
				}
				cr.Replace(cast)
			}
			// TODO: This code assumes that @foo::bool is on a single line
			var replace string
			if engine == config.EngineMySQL || !dollar {
				replace = "?"
			} else {
				replace = fmt.Sprintf("$%d", args[param.Name()].locs[0])
			}
			edits = append(edits, source.Edit{
				Location: expr.Location - raw.StmtLocation,
				Old:      fmt.Sprintf("@%s", paramName),
				New:      replace,
			})
			return false

		case named.IsParamSign(node):
			expr := node.(*ast.A_Expr)
			paramName, _ := flatten(expr.Rexpr)
			param := paramFromName(paramName)
			if p, ok := args[param.Name()]; ok {
				cr.Replace(&ast.ParamRef{
					Number:   p.locs[0],
					Location: expr.Location,
				})
			} else {
				argn++
				for numbs[argn] {
					argn++
				}

				args[param.Name()] = args[param.Name()].AddInstance(argn, param)
				cr.Replace(&ast.ParamRef{
					Number:   argn,
					Location: expr.Location,
				})
			}
			// TODO: This code assumes that @foo is on a single line
			var replace string
			if engine == config.EngineMySQL || !dollar {
				replace = "?"
			} else {
				replace = fmt.Sprintf("$%d", args[param.Name()].locs[0])
			}

			edits = append(edits, source.Edit{
				Location: expr.Location - raw.StmtLocation,
				Old:      fmt.Sprintf("@%s", paramName),
				New:      replace,
			})
			return false

		default:
			return true
		}
	}, nil)

	paramByLoc := map[int]named.Param{}
	for _, namedParam := range args {
		for _, loc := range namedParam.locs {
			paramByLoc[loc] = namedParam.param
		}
	}

	return node.(*ast.RawStmt), paramByLoc, edits
}

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

func NamedParameters(engine config.Engine, raw *ast.RawStmt, numbs map[int]bool, dollar bool) (*ast.RawStmt, map[int]string, []source.Edit) {
	foundFunc := astutils.Search(raw, named.IsParamFunc)
	foundSign := astutils.Search(raw, named.IsParamSign)
	if len(foundFunc.Items)+len(foundSign.Items) == 0 {
		return raw, map[int]string{}, nil
	}

	hasNamedParameterSupport := engine != config.EngineMySQL

	args := map[string][]int{}
	argn := 0
	var edits []source.Edit
	node := astutils.Apply(raw, func(cr *astutils.Cursor) bool {
		node := cr.Node()
		switch {
		case named.IsParamFunc(node):
			fun := node.(*ast.FuncCall)
			param, isConst := flatten(fun.Args)
			if nums, ok := args[param]; ok && hasNamedParameterSupport {
				cr.Replace(&ast.ParamRef{
					Number:   nums[0],
					Location: fun.Location,
				})
			} else {
				argn++
				for numbs[argn] {
					argn++
				}
				if _, found := args[param]; !found {
					args[param] = []int{argn}
				} else {
					args[param] = append(args[param], argn)
				}
				cr.Replace(&ast.ParamRef{
					Number:   argn,
					Location: fun.Location,
				})
			}
			// TODO: This code assumes that sqlc.arg(name) is on a single line
			var old, replace string
			if isConst {
				old = fmt.Sprintf("sqlc.arg('%s')", param)
			} else {
				old = fmt.Sprintf("sqlc.arg(%s)", param)
			}
			if engine == config.EngineMySQL || !dollar {
				replace = "?"
			} else {
				replace = fmt.Sprintf("$%d", args[param][0])
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
			param, _ := flatten(cast.Arg)
			if nums, ok := args[param]; ok {
				cast.Arg = &ast.ParamRef{
					Number:   nums[0],
					Location: expr.Location,
				}
				cr.Replace(cast)
			} else {
				argn++
				for numbs[argn] {
					argn++
				}
				if _, found := args[param]; !found {
					args[param] = []int{argn}
				} else {
					args[param] = append(args[param], argn)
				}
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
				replace = fmt.Sprintf("$%d", args[param][0])
			}
			edits = append(edits, source.Edit{
				Location: expr.Location - raw.StmtLocation,
				Old:      fmt.Sprintf("@%s", param),
				New:      replace,
			})
			return false

		case named.IsParamSign(node):
			expr := node.(*ast.A_Expr)
			param, _ := flatten(expr.Rexpr)
			if nums, ok := args[param]; ok {
				cr.Replace(&ast.ParamRef{
					Number:   nums[0],
					Location: expr.Location,
				})
			} else {
				argn++
				for numbs[argn] {
					argn++
				}
				if _, found := args[param]; !found {
					args[param] = []int{argn}
				} else {
					args[param] = append(args[param], argn)
				}
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
				replace = fmt.Sprintf("$%d", args[param][0])
			}
			edits = append(edits, source.Edit{
				Location: expr.Location - raw.StmtLocation,
				Old:      fmt.Sprintf("@%s", param),
				New:      replace,
			})
			return false

		default:
			return true
		}
	}, nil)

	named := map[int]string{}
	for k, vs := range args {
		for _, v := range vs {
			named[v] = k
		}
	}
	return node.(*ast.RawStmt), named, edits
}

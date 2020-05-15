package rewrite

import (
	"fmt"

	nodes "github.com/lfittl/pg_query_go/nodes"

	"github.com/kyleconroy/sqlc/internal/postgresql"
	"github.com/kyleconroy/sqlc/internal/postgresql/ast"
	"github.com/kyleconroy/sqlc/internal/source"
)

// Given an AST node, return the string representation of names
func flatten(root nodes.Node) (string, bool) {
	sw := &stringWalker{}
	ast.Walk(sw, root)
	return sw.String, sw.IsConst
}

type stringWalker struct {
	String  string
	IsConst bool
}

func (s *stringWalker) Visit(node nodes.Node) ast.Visitor {
	if _, ok := node.(nodes.A_Const); ok {
		s.IsConst = true
	}
	if n, ok := node.(nodes.String); ok {
		s.String += n.Str
	}
	return s
}

func isNamedParamSignCast(node nodes.Node) bool {
	expr, ok := node.(nodes.A_Expr)
	if !ok {
		return false
	}
	_, cast := expr.Rexpr.(nodes.TypeCast)
	return ast.Join(expr.Name, ".") == "@" && cast
}

func NamedParameters(raw nodes.RawStmt) (nodes.RawStmt, map[int]string, []source.Edit) {
	foundFunc := ast.Search(raw, postgresql.IsNamedParamFunc)
	foundSign := ast.Search(raw, postgresql.IsNamedParamSign)
	if len(foundFunc.Items)+len(foundSign.Items) == 0 {
		return raw, map[int]string{}, nil
	}

	args := map[string]int{}
	argn := 0
	var edits []source.Edit
	node := ast.Apply(raw, func(cr *ast.Cursor) bool {
		node := cr.Node()
		switch {

		case postgresql.IsNamedParamFunc(node):
			fun := node.(nodes.FuncCall)
			param, isConst := flatten(fun.Args)
			if num, ok := args[param]; ok {
				cr.Replace(nodes.ParamRef{
					Number:   num,
					Location: fun.Location,
				})
			} else {
				argn += 1
				args[param] = argn
				cr.Replace(nodes.ParamRef{
					Number:   argn,
					Location: fun.Location,
				})
			}
			// TODO: This code assumes that sqlc.arg(name) is on a single line
			var old string
			if isConst {
				old = fmt.Sprintf("sqlc.arg('%s')", param)
			} else {
				old = fmt.Sprintf("sqlc.arg(%s)", param)
			}
			edits = append(edits, source.Edit{
				Location: fun.Location - raw.StmtLocation,
				Old:      old,
				New:      fmt.Sprintf("$%d", args[param]),
			})
			return false

		case isNamedParamSignCast(node):
			expr := node.(nodes.A_Expr)
			cast := expr.Rexpr.(nodes.TypeCast)
			param, _ := flatten(cast.Arg)
			if num, ok := args[param]; ok {
				cast.Arg = nodes.ParamRef{
					Number:   num,
					Location: expr.Location,
				}
				cr.Replace(cast)
			} else {
				argn += 1
				args[param] = argn
				cast.Arg = nodes.ParamRef{
					Number:   argn,
					Location: expr.Location,
				}
				cr.Replace(cast)
			}
			// TODO: This code assumes that @foo::bool is on a single line
			edits = append(edits, source.Edit{
				Location: expr.Location - raw.StmtLocation,
				Old:      fmt.Sprintf("@%s", param),
				New:      fmt.Sprintf("$%d", args[param]),
			})
			return false

		case postgresql.IsNamedParamSign(node):
			expr := node.(nodes.A_Expr)
			param, _ := flatten(expr.Rexpr)
			if num, ok := args[param]; ok {
				cr.Replace(nodes.ParamRef{
					Number:   num,
					Location: expr.Location,
				})
			} else {
				argn += 1
				args[param] = argn
				cr.Replace(nodes.ParamRef{
					Number:   argn,
					Location: expr.Location,
				})
			}
			// TODO: This code assumes that @foo is on a single line
			edits = append(edits, source.Edit{
				Location: expr.Location - raw.StmtLocation,
				Old:      fmt.Sprintf("@%s", param),
				New:      fmt.Sprintf("$%d", args[param]),
			})
			return false

		default:
			return true
		}
	}, nil)

	named := map[int]string{}
	for k, v := range args {
		named[v] = k
	}
	return node.(nodes.RawStmt), named, edits
}

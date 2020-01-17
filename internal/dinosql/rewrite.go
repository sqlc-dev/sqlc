package dinosql

import (
	"fmt"

	nodes "github.com/lfittl/pg_query_go/nodes"

	"github.com/kyleconroy/sqlc/internal/postgresql/ast"
)

// Given an AST node, return the string representation of names
func flatten(root nodes.Node) string {
	sw := &stringWalker{}
	ast.Walk(sw, root)
	return sw.String
}

type stringWalker struct {
	String string
}

func (s *stringWalker) Visit(node nodes.Node) ast.Visitor {
	if n, ok := node.(nodes.String); ok {
		s.String += n.Str
	}
	return s
}

func rewriteNamedParameters(raw nodes.RawStmt) (nodes.RawStmt, map[int]string, []edit) {
	found := search(raw, func(node nodes.Node) bool {
		fun, ok := node.(nodes.FuncCall)
		return ok && ast.Join(fun.Funcname, ".") == "sqlc.arg"
	})
	if len(found.Items) == 0 {
		return raw, map[int]string{}, nil
	}

	args := map[string]int{}
	argn := 0
	var edits []edit
	node := ast.Apply(raw, func(cr *ast.Cursor) bool {
		fun, ok := cr.Node().(nodes.FuncCall)
		if !ok {
			return true
		}
		if ast.Join(fun.Funcname, ".") == "sqlc.arg" {
			param := flatten(fun.Args)
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
			edits = append(edits, edit{
				Location: fun.Location - raw.StmtLocation,
				Old:      fmt.Sprintf("sqlc.arg(%s)", param),
				New:      fmt.Sprintf("$%d", args[param]),
			})

			return false
		}
		return true
	}, nil)

	named := map[int]string{}
	for k, v := range args {
		named[v] = k
	}
	return node.(nodes.RawStmt), named, edits
}

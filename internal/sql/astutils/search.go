package astutils

import "github.com/sqlc-dev/sqlc/internal/sql/ast"

type nodeSearch struct {
	list  *ast.List
	check func(ast.Node) bool
}

func (s *nodeSearch) Visit(node ast.Node) Visitor {
	if s.check(node) {
		s.list.Items = append(s.list.Items, node)
	}
	return s
}

func Search(root ast.Node, f func(ast.Node) bool) *ast.List {
	ns := &nodeSearch{check: f, list: &ast.List{}}
	Walk(ns, root)
	return ns.list
}

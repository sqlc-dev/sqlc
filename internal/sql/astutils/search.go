package astutils

import "github.com/kyleconroy/sqlc/internal/sql/ast"

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

func IsChildOfNodes(parents []*ast.Node, node *ast.Node) bool {
	for _, v := range parents {
		if IsChildOfNode(v, node) {
			return true
		}
	}
	return false
}

func IsChildOfNode(parent *ast.Node, node *ast.Node) bool {
	res := Search(*parent, func(n ast.Node) bool {
		if n == *node {
			return true
		}
		return false
	})
	if len(res.Items) > 0 {
		return true
	}
	return false
}

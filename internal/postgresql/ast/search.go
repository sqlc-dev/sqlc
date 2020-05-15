package ast

import (
	nodes "github.com/lfittl/pg_query_go/nodes"
)

type nodeSearch struct {
	list  nodes.List
	check func(nodes.Node) bool
}

func (s *nodeSearch) Visit(node nodes.Node) Visitor {
	if s.check(node) {
		s.list.Items = append(s.list.Items, node)
	}
	return s
}

func Search(root nodes.Node, f func(nodes.Node) bool) nodes.List {
	ns := &nodeSearch{check: f}
	Walk(ns, root)
	return ns.list
}

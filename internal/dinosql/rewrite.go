package dinosql

import (
	nodes "github.com/lfittl/pg_query_go/nodes"
)

type stringWalker struct {
	String string
}

func (s *stringWalker) Visit(node nodes.Node) Visitor {
	if n, ok := node.(nodes.String); ok {
		s.String += n.Str
	}
	return s
}

func flatten(root nodes.Node) string {
	sw := &stringWalker{}
	Walk(sw, root)
	return sw.String
}

func rewriteNamedParameters(raw nodes.RawStmt) (nodes.RawStmt, error) {
	named := search(raw, func(node nodes.Node) bool {
		fun, ok := node.(nodes.FuncCall)
		return ok && join(fun.Funcname, ".") == "sqlc.arg"
	})
	for _, np := range named.Items {
		fun := np.(nodes.FuncCall)
		flatten(fun.Args)
	}
	return raw, nil
}

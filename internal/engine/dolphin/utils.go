package dolphin

import (
	pcast "github.com/pingcap/tidb/parser/ast"

	"github.com/kyleconroy/sqlc/internal/sql/ast"
)

type nodeSearch struct {
	list  []pcast.Node
	check func(pcast.Node) bool
}

func (s *nodeSearch) Enter(n pcast.Node) (pcast.Node, bool) {
	if s.check(n) {
		s.list = append(s.list, n)
	}
	return n, false // skipChildren
}

func (s *nodeSearch) Leave(n pcast.Node) (pcast.Node, bool) {
	return n, true // ok
}

func collect(root pcast.Node, f func(pcast.Node) bool) []pcast.Node {
	if root == nil {
		return nil
	}
	ns := &nodeSearch{check: f}
	root.Accept(ns)
	return ns.list
}

type nodeVisit struct {
	fn func(pcast.Node)
}

func (s *nodeVisit) Enter(n pcast.Node) (pcast.Node, bool) {
	s.fn(n)
	return n, false // skipChildren
}

func (s *nodeVisit) Leave(n pcast.Node) (pcast.Node, bool) {
	return n, true // ok
}

func visit(root pcast.Node, f func(pcast.Node)) {
	if root == nil {
		return
	}
	ns := &nodeVisit{fn: f}
	root.Accept(ns)
}

// Maybe not useful?
func text(nodes []pcast.Node) []string {
	str := make([]string, len(nodes))
	for i := range nodes {
		if nodes[i] == nil {
			continue
		}
		str[i] = nodes[i].Text()
	}
	return str
}

func parseTableName(n *pcast.TableName) *ast.TableName {
	return &ast.TableName{
		Schema: identifier(n.Schema.String()),
		Name:   identifier(n.Name.String()),
	}
}

func toList(node pcast.Node) *ast.List {
	var items []ast.Node
	switch n := node.(type) {
	case *pcast.TableName:
		if schema := n.Schema.String(); schema != "" {
			items = append(items, NewIdentifier(schema))
		}
		items = append(items, NewIdentifier(n.Name.String()))
	default:
		return nil
	}
	return &ast.List{Items: items}
}

func isNotNull(n *pcast.ColumnDef) bool {
	for i := range n.Options {
		if n.Options[i].Tp == pcast.ColumnOptionNotNull {
			return true
		}
		if n.Options[i].Tp == pcast.ColumnOptionPrimaryKey {
			return true
		}
	}
	return false
}

func convertToRangeVarList(list *ast.List, result *ast.List) {
	if len(list.Items) == 0 {
		return
	}
	switch rel := list.Items[0].(type) {

	// Special case for joins in updates
	case *ast.JoinExpr:
		left, ok := rel.Larg.(*ast.RangeVar)
		if !ok {
			if list, check := rel.Larg.(*ast.List); check {
				convertToRangeVarList(list, result)
			} else {
				panic("expected range var")
			}
		}
		if left != nil {
			result.Items = append(result.Items, left)
		}

		right, ok := rel.Rarg.(*ast.RangeVar)
		if !ok {
			if list, check := rel.Rarg.(*ast.List); check {
				convertToRangeVarList(list, result)
			} else {
				panic("expected range var")
			}
		}
		if right != nil {
			result.Items = append(result.Items, right)
		}

	case *ast.RangeVar:
		result.Items = append(result.Items, rel)

	default:
		panic("expected range var")
	}
}

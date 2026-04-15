package dolphin

import (
	pcast "github.com/pingcap/tidb/pkg/parser/ast"
	"github.com/pingcap/tidb/pkg/parser/mysql"

	"github.com/sqlc-dev/sqlc/internal/sql/ast"
)

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
			} else if subselect, check := rel.Larg.(*ast.RangeSubselect); check {
				// Handle subqueries in JOIN clauses
				result.Items = append(result.Items, subselect)
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
			} else if subselect, check := rel.Rarg.(*ast.RangeSubselect); check {
				// Handle subqueries in JOIN clauses
				result.Items = append(result.Items, subselect)
			} else {
				panic("expected range var")
			}
		}
		if right != nil {
			result.Items = append(result.Items, right)
		}

	case *ast.RangeVar:
		result.Items = append(result.Items, rel)

	case *ast.RangeSubselect:
		result.Items = append(result.Items, rel)

	default:
		panic("expected range var")
	}
}

func isUnsigned(n *pcast.ColumnDef) bool {
	return mysql.HasUnsignedFlag(n.Tp.GetFlag())
}

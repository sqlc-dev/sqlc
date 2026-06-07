package compiler

import (
	"sort"

	"github.com/sqlc-dev/sqlc/internal/sql/ast"
	"github.com/sqlc-dev/sqlc/internal/sql/astutils"
)

// sourceTableNames returns the sorted, deduplicated names of the base tables a
// statement reads from. It covers every table referenced in a FROM, a JOIN, or
// a subquery in any clause, including the bodies of common table expressions.
// The names of common table expressions and the target relations of INSERT,
// UPDATE, DELETE, and TRUNCATE statements are not reads and are excluded.
func sourceTableNames(root ast.Node) []string {
	cteNames := map[string]struct{}{}
	writeTargets := map[*ast.RangeVar]struct{}{}

	collect := astutils.VisitorFunc(func(node ast.Node) {
		switch n := node.(type) {
		case *ast.CommonTableExpr:
			if n.Ctename != nil {
				cteNames[*n.Ctename] = struct{}{}
			}
		case *ast.InsertStmt:
			if n.Relation != nil {
				markRangeVars(writeTargets, n.Relation)
			}
		case *ast.UpdateStmt:
			if n.Relations != nil {
				markRangeVars(writeTargets, n.Relations)
			}
		case *ast.DeleteStmt:
			if n.Relations != nil {
				markRangeVars(writeTargets, n.Relations)
			}
		case *ast.TruncateStmt:
			if n.Relations != nil {
				markRangeVars(writeTargets, n.Relations)
			}
		}
	})
	astutils.Walk(collect, root)

	seen := map[string]struct{}{}
	var names []string
	for _, rv := range rangeVars(root) {
		if _, ok := writeTargets[rv]; ok {
			continue
		}
		table, err := ParseTableName(rv)
		if err != nil {
			continue
		}
		if _, ok := cteNames[table.Name]; ok {
			continue
		}
		if _, ok := seen[table.Name]; ok {
			continue
		}
		seen[table.Name] = struct{}{}
		names = append(names, table.Name)
	}
	sort.Strings(names)
	return names
}

func markRangeVars(set map[*ast.RangeVar]struct{}, node ast.Node) {
	for _, rv := range rangeVars(node) {
		set[rv] = struct{}{}
	}
}

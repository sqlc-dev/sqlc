package compiler

import (
	"github.com/sqlc-dev/sqlc/internal/sql/ast"
	"github.com/sqlc-dev/sqlc/internal/sql/astutils"
)

// UsingInfo tracks USING columns for a specific join to avoid duplication
// when expanding SELECT * across multiple tables
type UsingInfo struct {
	// columns is a set of column names that appear in the USING clause
	columns map[string]bool
}

// NewUsingInfo creates a new UsingInfo from a JoinExpr's UsingClause
func NewUsingInfo(joinExpr *ast.JoinExpr) *UsingInfo {
	ui := &UsingInfo{
		columns: make(map[string]bool),
	}

	if joinExpr == nil || joinExpr.UsingClause == nil {
		return ui
	}

	// Extract column names from the USING clause
	for _, item := range joinExpr.UsingClause.Items {
		if str, ok := item.(*ast.String); ok {
			ui.columns[str.Str] = true
		}
	}

	return ui
}

// HasColumn checks if a column name is in the USING clause
func (ui *UsingInfo) HasColumn(colName string) bool {
	return ui.columns[colName]
}

// getJoinUsingMap builds a map of table names to their USING columns
// This helps identify which columns should not be duplicated when expanding *
func getJoinUsingMap(node ast.Node) map[string]*UsingInfo {
	usingMap := make(map[string]*UsingInfo)

	// Find all JoinExpr nodes in the query and extract USING information
	visitor := &joinVisitor{
		usingMap: usingMap,
	}
	astutils.Walk(visitor, node)

	return usingMap
}

// joinVisitor traverses the AST to find and track USING information in joins
type joinVisitor struct {
	usingMap map[string]*UsingInfo
}

func (v *joinVisitor) Visit(node ast.Node) astutils.Visitor {
	if join, ok := node.(*ast.JoinExpr); ok {
		// Create UsingInfo for this join
		// The right argument of the join is the table being joined in
		if rarg, ok := join.Rarg.(*ast.RangeVar); ok {
			if rarg.Relname != nil {
				usingInfo := NewUsingInfo(join)
				if len(usingInfo.columns) > 0 {
					v.usingMap[*rarg.Relname] = usingInfo
				}
			}
		}
	}
	return v
}

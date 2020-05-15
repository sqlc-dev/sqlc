package pg

import (
	"github.com/kyleconroy/sqlc/internal/sql/ast"
)

type SelectStmt struct {
	DistinctClause *ast.List
	IntoClause     *IntoClause
	TargetList     *ast.List
	FromClause     *ast.List
	WhereClause    ast.Node
	GroupClause    *ast.List
	HavingClause   ast.Node
	WindowClause   *ast.List
	ValuesLists    [][]ast.Node
	SortClause     *ast.List
	LimitOffset    ast.Node
	LimitCount     ast.Node
	LockingClause  *ast.List
	WithClause     *WithClause
	Op             SetOperation
	All            bool
	Larg           *SelectStmt
	Rarg           *SelectStmt
}

func (n *SelectStmt) Pos() int {
	return 0
}

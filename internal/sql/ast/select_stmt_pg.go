package ast

import (
	"github.com/kyleconroy/sqlc/internal/sql/ast"
)

type SelectStmt struct {
	DistinctClause *List
	IntoClause     *IntoClause
	TargetList     *List
	FromClause     *List
	WhereClause    ast.Node
	GroupClause    *List
	HavingClause   ast.Node
	WindowClause   *List
	ValuesLists    *List
	SortClause     *List
	LimitOffset    ast.Node
	LimitCount     ast.Node
	LockingClause  *List
	WithClause     *WithClause
	Op             SetOperation
	All            bool
	Larg           *SelectStmt
	Rarg           *SelectStmt
}

func (n *SelectStmt) Pos() int {
	return 0
}

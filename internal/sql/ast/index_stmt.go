package ast

import (
	"github.com/kyleconroy/sqlc/internal/sql/ast"
)

type IndexStmt struct {
	Idxname        *string
	Relation       *RangeVar
	AccessMethod   *string
	TableSpace     *string
	IndexParams    *List
	Options        *List
	WhereClause    ast.Node
	ExcludeOpNames *List
	Idxcomment     *string
	IndexOid       Oid
	OldNode        Oid
	Unique         bool
	Primary        bool
	Isconstraint   bool
	Deferrable     bool
	Initdeferred   bool
	Transformed    bool
	Concurrent     bool
	IfNotExists    bool
}

func (n *IndexStmt) Pos() int {
	return 0
}

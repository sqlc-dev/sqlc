package pg

import (
	"github.com/kyleconroy/sqlc/internal/sql/ast"
)

type IndexStmt struct {
	Idxname        *string
	Relation       *RangeVar
	AccessMethod   *string
	TableSpace     *string
	IndexParams    *ast.List
	Options        *ast.List
	WhereClause    ast.Node
	ExcludeOpNames *ast.List
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

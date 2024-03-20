package ast

type IndexStmt struct {
	Idxname        *string
	Relation       *RangeVar
	AccessMethod   *string
	TableSpace     *string
	IndexParams    *List
	Options        *List
	WhereClause    Node
	ExcludeOpNames *List
	Idxcomment     *string
	IndexOid       Oid
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

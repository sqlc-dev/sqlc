package ast

type IndexStmt struct {
	Tag NodeTag[IndexStmt] `json:"tag"`

	Idxname        *string   `json:",omitempty"`
	Relation       *RangeVar `json:",omitempty"`
	AccessMethod   *string   `json:",omitempty"`
	TableSpace     *string   `json:",omitempty"`
	IndexParams    *List     `json:",omitempty"`
	Options        *List     `json:",omitempty"`
	WhereClause    Node      `json:",omitempty"`
	ExcludeOpNames *List     `json:",omitempty"`
	Idxcomment     *string   `json:",omitempty"`
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

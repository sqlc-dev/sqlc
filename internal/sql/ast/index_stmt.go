package ast

type IndexStmt struct {
	Tag NodeTag[IndexStmt] `json:"tag"`

	Idxname        *string   `json:"idxname,omitempty"`
	Relation       *RangeVar `json:"relation,omitempty"`
	AccessMethod   *string   `json:"access_method,omitempty"`
	TableSpace     *string   `json:"table_space,omitempty"`
	IndexParams    *List     `json:"index_params,omitempty"`
	Options        *List     `json:"options,omitempty"`
	WhereClause    Node      `json:"where_clause,omitempty"`
	ExcludeOpNames *List     `json:"exclude_op_names,omitempty"`
	Idxcomment     *string   `json:"idxcomment,omitempty"`
	IndexOid       Oid       `json:"index_oid"`
	Unique         bool      `json:"unique"`
	Primary        bool      `json:"primary"`
	Isconstraint   bool      `json:"isconstraint"`
	Deferrable     bool      `json:"deferrable"`
	Initdeferred   bool      `json:"initdeferred"`
	Transformed    bool      `json:"transformed"`
	Concurrent     bool      `json:"concurrent"`
	IfNotExists    bool      `json:"if_not_exists"`
}

func (n *IndexStmt) Pos() int {
	return 0
}

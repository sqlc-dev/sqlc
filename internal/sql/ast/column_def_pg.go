package ast

type ColumnDef_PG struct {
	Colname       *string
	TypeName      *TypeName
	Inhcount      int
	IsLocal       bool
	IsNotNull     bool
	IsFromType    bool
	IsFromParent  bool
	Storage       byte
	RawDefault    Node
	CookedDefault Node
	Identity      byte
	CollClause    *CollateClause
	CollOid       Oid
	Constraints   *List
	Fdwoptions    *List
	Location      int
}

func (n *ColumnDef_PG) Pos() int {
	return n.Location
}

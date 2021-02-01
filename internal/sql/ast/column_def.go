package ast

type ColumnDef struct {
	Colname   string
	TypeName  *TypeName
	IsNotNull bool
	IsArray   bool
	Vals      *List
	Length    *int

	// From pg.ColumnDef
	Inhcount      int
	IsLocal       bool
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
	Comment       string
}

func (n *ColumnDef) Pos() int {
	return n.Location
}

package ast

type ColumnDef struct {
	Colname    string
	TypeName   *TypeName
	IsNotNull  bool
	IsUnsigned bool
	IsArray    bool
	ArrayDims  int
	Vals       *List
	Length     *int
	PrimaryKey bool

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

func (n *ColumnDef) Format(buf *TrackedBuffer) {
	if n == nil {
		return
	}
	buf.WriteString(n.Colname)
	buf.WriteString(" ")
	buf.astFormat(n.TypeName)
	if n.PrimaryKey {
		buf.WriteString(" PRIMARY KEY")
	} else if n.IsNotNull {
		buf.WriteString(" NOT NULL")
	}
	buf.astFormat(n.Constraints)
}

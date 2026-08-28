package ast

import "github.com/sqlc-dev/sqlc/internal/sql/format"

type ColumnDef struct {
	Tag NodeTag[ColumnDef] `json:"tag"`

	Colname string
	// TypeName is the column's type for the catalog. When Typeless is set
	// the author wrote no type at all — SQLite allows it — and the column
	// prints without one, whatever TypeName carries.
	TypeName   *TypeName `json:",omitempty"`
	Typeless   bool
	IsNotNull  bool
	IsUnsigned bool
	IsArray    bool
	ArrayDims  int
	Vals       *List `json:",omitempty"`
	Length     *int  `json:",omitempty"`
	PrimaryKey bool
	// IsHidden marks a column a relation offers by name without listing it,
	// like the column an sqlite fts5 table names after itself. The legacy
	// catalog drops hidden columns; the core catalog keeps them out of star
	// expansions and models.
	IsHidden bool

	// From pg.ColumnDef
	Inhcount      int
	IsLocal       bool
	IsFromType    bool
	IsFromParent  bool
	Storage       byte
	RawDefault    Node `json:",omitempty"`
	CookedDefault Node `json:",omitempty"`
	Identity      byte
	CollClause    *CollateClause `json:",omitempty"`
	CollOid       Oid
	Constraints   *List `json:",omitempty"`
	Fdwoptions    *List `json:",omitempty"`
	Location      int
	Comment       string
}

func (n *ColumnDef) Pos() int {
	return n.Location
}

func (n *ColumnDef) Format(buf *TrackedBuffer, d format.Dialect) {
	if n == nil {
		return
	}
	buf.WriteString(n.Colname)
	if !n.Typeless {
		buf.WriteString(" ")
		buf.astFormat(n.TypeName, d)
	}
	// Use IsArray from ColumnDef since TypeName.ArrayBounds may not be set
	// (for type resolution compatibility)
	if n.IsArray && !items(n.TypeName.ArrayBounds) {
		buf.WriteString("[]")
	}
	if n.PrimaryKey {
		buf.WriteString(" PRIMARY KEY")
	} else if n.IsNotNull {
		buf.WriteString(" NOT NULL")
	}
	buf.astFormat(n.Constraints, d)
}

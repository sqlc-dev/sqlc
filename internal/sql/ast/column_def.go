package ast

import "github.com/sqlc-dev/sqlc/internal/sql/format"

type ColumnDef struct {
	Tag NodeTag[ColumnDef] `json:"tag"`

	Colname string
	// TypeName is the column's type for the catalog. When Typeless is set
	// the author wrote no type at all — SQLite allows it — and the column
	// prints without one, whatever TypeName carries.
	TypeName   *TypeName
	Typeless   bool
	IsNotNull  bool
	IsUnsigned bool
	IsArray    bool
	ArrayDims  int
	Vals       *List
	Length     *int
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

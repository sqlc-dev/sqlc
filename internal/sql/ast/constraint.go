package ast

type Constraint struct {
	Tag NodeTag[Constraint] `json:"tag"`

	Contype        ConstrType
	Conname        *string `json:",omitempty"`
	Deferrable     bool
	Initdeferred   bool
	Location       int
	IsNoInherit    bool
	RawExpr        Node    `json:",omitempty"`
	CookedExpr     *string `json:",omitempty"`
	GeneratedWhen  byte
	Keys           *List     `json:",omitempty"`
	Exclusions     *List     `json:",omitempty"`
	Options        *List     `json:",omitempty"`
	Indexname      *string   `json:",omitempty"`
	Indexspace     *string   `json:",omitempty"`
	AccessMethod   *string   `json:",omitempty"`
	WhereClause    Node      `json:",omitempty"`
	Pktable        *RangeVar `json:",omitempty"`
	FkAttrs        *List     `json:",omitempty"`
	PkAttrs        *List     `json:",omitempty"`
	FkMatchtype    byte
	FkUpdAction    byte
	FkDelAction    byte
	OldConpfeqop   *List `json:",omitempty"`
	OldPktableOid  Oid
	SkipValidation bool
	InitiallyValid bool
}

func (n *Constraint) Pos() int {
	return n.Location
}

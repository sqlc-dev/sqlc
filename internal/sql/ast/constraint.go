package ast

type Constraint struct {
	Contype        ConstrType
	Conname        *string
	Deferrable     bool
	Initdeferred   bool
	Location       int
	IsNoInherit    bool
	RawExpr        Node
	CookedExpr     *string
	GeneratedWhen  byte
	Keys           *List
	Exclusions     *List
	Options        *List
	Indexname      *string
	Indexspace     *string
	AccessMethod   *string
	WhereClause    Node
	Pktable        *RangeVar
	FkAttrs        *List
	PkAttrs        *List
	FkMatchtype    byte
	FkUpdAction    byte
	FkDelAction    byte
	OldConpfeqop   *List
	OldPktableOid  Oid
	SkipValidation bool
	InitiallyValid bool
}

func (n *Constraint) Pos() int {
	return n.Location
}

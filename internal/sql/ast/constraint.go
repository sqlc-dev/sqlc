package ast

type Constraint struct {
	Tag NodeTag[Constraint] `json:"tag"`

	Contype        ConstrType `json:"contype"`
	Conname        *string    `json:"conname,omitempty"`
	Deferrable     bool       `json:"deferrable"`
	Initdeferred   bool       `json:"initdeferred"`
	Location       int        `json:"location"`
	IsNoInherit    bool       `json:"is_no_inherit"`
	RawExpr        Node       `json:"raw_expr,omitempty"`
	CookedExpr     *string    `json:"cooked_expr,omitempty"`
	GeneratedWhen  byte       `json:"generated_when"`
	Keys           *List      `json:"keys,omitempty"`
	Exclusions     *List      `json:"exclusions,omitempty"`
	Options        *List      `json:"options,omitempty"`
	Indexname      *string    `json:"indexname,omitempty"`
	Indexspace     *string    `json:"indexspace,omitempty"`
	AccessMethod   *string    `json:"access_method,omitempty"`
	WhereClause    Node       `json:"where_clause,omitempty"`
	Pktable        *RangeVar  `json:"pktable,omitempty"`
	FkAttrs        *List      `json:"fk_attrs,omitempty"`
	PkAttrs        *List      `json:"pk_attrs,omitempty"`
	FkMatchtype    byte       `json:"fk_matchtype"`
	FkUpdAction    byte       `json:"fk_upd_action"`
	FkDelAction    byte       `json:"fk_del_action"`
	OldConpfeqop   *List      `json:"old_conpfeqop,omitempty"`
	OldPktableOid  Oid        `json:"old_pktable_oid"`
	SkipValidation bool       `json:"skip_validation"`
	InitiallyValid bool       `json:"initially_valid"`
}

func (n *Constraint) Pos() int {
	return n.Location
}

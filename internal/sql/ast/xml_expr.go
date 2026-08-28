package ast

type XmlExpr struct {
	Tag NodeTag[XmlExpr] `json:"tag"`

	Xpr       Node `json:",omitempty"`
	Op        XmlExprOp
	Name      *string `json:",omitempty"`
	NamedArgs *List   `json:",omitempty"`
	ArgNames  *List   `json:",omitempty"`
	Args      *List   `json:",omitempty"`
	Xmloption XmlOptionType
	Type      Oid
	Typmod    int32
	Location  int
}

func (n *XmlExpr) Pos() int {
	return n.Location
}

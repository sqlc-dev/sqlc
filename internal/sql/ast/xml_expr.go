package ast

type XmlExpr struct {
	Tag NodeTag[XmlExpr] `json:"tag"`

	Xpr       Node          `json:"xpr,omitempty"`
	Op        XmlExprOp     `json:"op"`
	Name      *string       `json:"name,omitempty"`
	NamedArgs *List         `json:"named_args,omitempty"`
	ArgNames  *List         `json:"arg_names,omitempty"`
	Args      *List         `json:"args,omitempty"`
	Xmloption XmlOptionType `json:"xmloption"`
	Type      Oid           `json:"type"`
	Typmod    int32         `json:"typmod"`
	Location  int           `json:"location"`
}

func (n *XmlExpr) Pos() int {
	return n.Location
}

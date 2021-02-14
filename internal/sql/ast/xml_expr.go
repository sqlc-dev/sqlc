package ast

type XmlExpr struct {
	Xpr       Node
	Op        XmlExprOp
	Name      *string
	NamedArgs *List
	ArgNames  *List
	Args      *List
	Xmloption XmlOptionType
	Type      Oid
	Typmod    int32
	Location  int
}

func (n *XmlExpr) Pos() int {
	return n.Location
}

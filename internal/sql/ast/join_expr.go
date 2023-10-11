package ast

type JoinExpr struct {
	Jointype    JoinType
	IsNatural   bool
	Larg        Node
	Rarg        Node
	UsingClause *List
	Quals       Node
	Alias       *Alias
	Rtindex     int
}

func (n *JoinExpr) Pos() int {
	return 0
}

func (n *JoinExpr) Format(buf *TrackedBuffer) {
	if n == nil {
		return
	}
	buf.astFormat(n.Larg)
	buf.WriteString(" JOIN ")
	buf.astFormat(n.Rarg)
	buf.WriteString(" ON ")
	buf.astFormat(n.Quals)
}

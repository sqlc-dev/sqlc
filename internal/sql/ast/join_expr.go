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
	switch n.Jointype {
	case JoinTypeLeft:
		buf.WriteString(" LEFT JOIN ")
	case JoinTypeInner:
		buf.WriteString(" INNER JOIN ")
	default:
		buf.WriteString(" JOIN ")
	}
	buf.astFormat(n.Rarg)
	buf.WriteString(" ON ")
	if n.Jointype == JoinTypeInner {
		if set(n.Quals) {
			buf.astFormat(n.Quals)
		} else {
			buf.WriteString("TRUE")
		}
	} else {
		buf.astFormat(n.Quals)
	}
}

package ast

import "github.com/sqlc-dev/sqlc/internal/sql/format"

type JoinExpr struct {
	Tag NodeTag[JoinExpr] `json:"tag"`

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

func (n *JoinExpr) Format(buf *TrackedBuffer, d format.Dialect) {
	if n == nil {
		return
	}
	buf.astFormat(n.Larg, d)
	if n.Jointype == JoinTypeComma {
		buf.WriteString(",")
		buf.beforeClause(n.Rarg, d)
		buf.Line()
	} else {
		buf.beforeClause(n.Rarg, d)
		buf.Line()
		if n.IsNatural {
			buf.WriteString("NATURAL ")
		}
		switch n.Jointype {
		case JoinTypeLeft:
			buf.WriteString("LEFT JOIN ")
		case JoinTypeRight:
			buf.WriteString("RIGHT JOIN ")
		case JoinTypeFull:
			buf.WriteString("FULL JOIN ")
		case JoinTypeCross:
			buf.WriteString("CROSS JOIN ")
		default:
			buf.WriteString("JOIN ")
		}
	}
	buf.astFormat(n.Rarg, d)
	if items(n.UsingClause) {
		buf.WriteString(" USING (")
		buf.join(n.UsingClause, d, ", ")
		buf.WriteString(")")
	} else if set(n.Quals) {
		buf.WriteString(" ON ")
		buf.condition(n.Quals, d)
	}
}

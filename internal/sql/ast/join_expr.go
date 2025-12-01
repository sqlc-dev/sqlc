package ast

import "github.com/sqlc-dev/sqlc/internal/sql/format"

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

func (n *JoinExpr) Format(buf *TrackedBuffer, d format.Dialect) {
	if n == nil {
		return
	}
	buf.astFormat(n.Larg, d)
	if n.IsNatural {
		buf.WriteString(" NATURAL")
	}
	switch n.Jointype {
	case JoinTypeLeft:
		buf.WriteString(" LEFT JOIN ")
	case JoinTypeRight:
		buf.WriteString(" RIGHT JOIN ")
	case JoinTypeFull:
		buf.WriteString(" FULL JOIN ")
	case JoinTypeInner:
		// CROSS JOIN has no ON or USING clause
		if !items(n.UsingClause) && !set(n.Quals) {
			buf.WriteString(" CROSS JOIN ")
		} else {
			buf.WriteString(" JOIN ")
		}
	default:
		buf.WriteString(" JOIN ")
	}
	buf.astFormat(n.Rarg, d)
	if items(n.UsingClause) {
		buf.WriteString(" USING (")
		buf.join(n.UsingClause, d, ", ")
		buf.WriteString(")")
	} else if set(n.Quals) {
		buf.WriteString(" ON ")
		buf.astFormat(n.Quals, d)
	}
}

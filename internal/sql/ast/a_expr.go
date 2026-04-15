package ast

import "github.com/sqlc-dev/sqlc/internal/sql/format"

type A_Expr struct {
	Kind     A_Expr_Kind
	Name     *List
	Lexpr    Node
	Rexpr    Node
	Location int
}

func (n *A_Expr) Pos() int {
	return n.Location
}

// isNamedParam returns true if this A_Expr represents a named parameter (@name)
// and extracts the parameter name if so.
func (n *A_Expr) isNamedParam() (string, bool) {
	if n.Name == nil || len(n.Name.Items) != 1 {
		return "", false
	}
	s, ok := n.Name.Items[0].(*String)
	if !ok || s.Str != "@" {
		return "", false
	}
	if set(n.Lexpr) || !set(n.Rexpr) {
		return "", false
	}
	if nameStr, ok := n.Rexpr.(*String); ok {
		return nameStr.Str, true
	}
	return "", false
}

func (n *A_Expr) Format(buf *TrackedBuffer, d format.Dialect) {
	if n == nil {
		return
	}

	// Check for named parameter first (works regardless of Kind)
	if name, ok := n.isNamedParam(); ok {
		buf.WriteString(d.NamedParam(name))
		return
	}

	switch n.Kind {
	case A_Expr_Kind_IN:
		buf.astFormat(n.Lexpr, d)
		buf.WriteString(" IN (")
		buf.astFormat(n.Rexpr, d)
		buf.WriteString(")")
	case A_Expr_Kind_LIKE:
		buf.astFormat(n.Lexpr, d)
		buf.WriteString(" LIKE ")
		buf.astFormat(n.Rexpr, d)
	case A_Expr_Kind_ILIKE:
		buf.astFormat(n.Lexpr, d)
		buf.WriteString(" ILIKE ")
		buf.astFormat(n.Rexpr, d)
	case A_Expr_Kind_SIMILAR:
		buf.astFormat(n.Lexpr, d)
		buf.WriteString(" SIMILAR TO ")
		buf.astFormat(n.Rexpr, d)
	case A_Expr_Kind_BETWEEN:
		buf.astFormat(n.Lexpr, d)
		buf.WriteString(" BETWEEN ")
		if l, ok := n.Rexpr.(*List); ok && len(l.Items) == 2 {
			buf.astFormat(l.Items[0], d)
			buf.WriteString(" AND ")
			buf.astFormat(l.Items[1], d)
		}
	case A_Expr_Kind_NOT_BETWEEN:
		buf.astFormat(n.Lexpr, d)
		buf.WriteString(" NOT BETWEEN ")
		if l, ok := n.Rexpr.(*List); ok && len(l.Items) == 2 {
			buf.astFormat(l.Items[0], d)
			buf.WriteString(" AND ")
			buf.astFormat(l.Items[1], d)
		}
	case A_Expr_Kind_DISTINCT:
		buf.astFormat(n.Lexpr, d)
		buf.WriteString(" IS DISTINCT FROM ")
		buf.astFormat(n.Rexpr, d)
	case A_Expr_Kind_NOT_DISTINCT:
		buf.astFormat(n.Lexpr, d)
		buf.WriteString(" IS NOT DISTINCT FROM ")
		buf.astFormat(n.Rexpr, d)
	case A_Expr_Kind_NULLIF:
		buf.WriteString("NULLIF(")
		buf.astFormat(n.Lexpr, d)
		buf.WriteString(", ")
		buf.astFormat(n.Rexpr, d)
		buf.WriteString(")")
	default:
		// Standard operator (including A_Expr_Kind_OP)
		if set(n.Lexpr) {
			buf.astFormat(n.Lexpr, d)
			buf.WriteString(" ")
		}
		buf.astFormat(n.Name, d)
		if set(n.Rexpr) {
			buf.WriteString(" ")
			buf.astFormat(n.Rexpr, d)
		}
	}
}

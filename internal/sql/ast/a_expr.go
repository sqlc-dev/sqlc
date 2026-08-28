package ast

import (
	"strings"

	"github.com/sqlc-dev/sqlc/internal/sql/format"
)

type A_Expr struct {
	Tag NodeTag[A_Expr] `json:"tag"`

	Kind A_Expr_Kind `json:"kind"`
	// Name is the operator, as a list of String nodes: PostgreSQL's
	// operator space is open (user-defined and schema-qualified operators),
	// so the shared tree keeps pg_query's shape rather than an enum. Each
	// engine's converter writes the operator as its parser saw it — for an
	// engine that accepts more than one spelling of the same operator
	// (SQLite's != for <>, == for =), that is the author's spelling, which
	// is how the printer preserves it.
	Name     *List `json:"name,omitempty"`
	Lexpr    Node  `json:"lexpr,omitempty"`
	Rexpr    Node  `json:"rexpr,omitempty"`
	Location int   `json:"location"`
}

func (n *A_Expr) Pos() int {
	return n.Location
}

// isNamedParam reports whether this A_Expr represents a named parameter,
// returning its sigil and name. Engines encode a named parameter as a
// prefix pseudo-operator carrying the sigil the author wrote (@name for
// most engines; SQLite also spells :name and $name). The spellings mean
// different things to sqlc, so the printer emits the one it was given.
func (n *A_Expr) isNamedParam() (sigil, name string, ok bool) {
	if n.Name == nil || len(n.Name.Items) != 1 {
		return "", "", false
	}
	s, sok := n.Name.Items[0].(*String)
	if !sok || (s.Str != "@" && s.Str != ":" && s.Str != "$") {
		return "", "", false
	}
	if set(n.Lexpr) || !set(n.Rexpr) {
		return "", "", false
	}
	if nameStr, sok := n.Rexpr.(*String); sok {
		return s.Str, nameStr.Str, true
	}
	// Before the compiler rewrites named parameters, @name parses as the @
	// operator applied to a bare column reference.
	if col, sok := n.Rexpr.(*ColumnRef); sok && col.Fields != nil && len(col.Fields.Items) == 1 {
		if f, sok := col.Fields.Items[0].(*String); sok {
			return s.Str, f.Str, true
		}
	}
	return "", "", false
}

// negated returns true when the expression's operator carries the negated
// spelling of a pattern-match operator (e.g. "!~~" for NOT LIKE).
func (n *A_Expr) negated() bool {
	if n.Name == nil || len(n.Name.Items) != 1 {
		return false
	}
	s, ok := n.Name.Items[0].(*String)
	return ok && strings.HasPrefix(s.Str, "!")
}

func (n *A_Expr) Format(buf *TrackedBuffer, d format.Dialect) {
	if n == nil {
		return
	}

	// Check for named parameter first (works regardless of Kind)
	if sigil, name, ok := n.isNamedParam(); ok {
		buf.WriteString(sigil)
		buf.WriteString(name)
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
		if n.negated() {
			buf.WriteString(" NOT LIKE ")
		} else {
			buf.WriteString(" LIKE ")
		}
		buf.astFormat(n.Rexpr, d)
	case A_Expr_Kind_ILIKE:
		buf.astFormat(n.Lexpr, d)
		if n.negated() {
			buf.WriteString(" NOT ILIKE ")
		} else {
			buf.WriteString(" ILIKE ")
		}
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

package ast

import (
	"github.com/sqlc-dev/sqlc/internal/debug"
	"github.com/sqlc-dev/sqlc/internal/sql/ast/printer"
	"github.com/sqlc-dev/sqlc/internal/sql/format"
)

type nodeFormatter interface {
	Format(*TrackedBuffer, format.Dialect)
}

// TrackedBuffer layers AST formatting over the document renderer in
// ast/printer: node dispatch, list joins, and comment emission. The
// renderer itself — tokens, groups, width fitting — lives in the printer
// package; TrackedBuffer stays here because every node's Format method
// names it, which an ast → printer → ast import cycle would forbid.
type TrackedBuffer struct {
	*printer.Buffer
	// ct, when set, carries the statement's interior comments attached to
	// their anchor nodes; the buffer emits each comment when it reaches the
	// comment's anchor. Positions were consulted once, in AttachComments —
	// never here — so edited and synthetic trees print comments correctly.
	ct *CommentTable
	// anchors, when set, records every positioned node and emission point
	// in print order instead of producing output; AttachComments uses this
	// dry run to classify comments against the same order printing will
	// use.
	anchors *[]anchor
}

// NewTrackedBuffer creates a new TrackedBuffer.
func NewTrackedBuffer() *TrackedBuffer {
	return &TrackedBuffer{Buffer: &printer.Buffer{}}
}

func (t *TrackedBuffer) astFormat(n Node, d format.Dialect) {
	if t.anchors != nil && n != nil && n.Pos() > 0 {
		*t.anchors = append(*t.anchors, anchor{node: n, pos: n.Pos()})
	}
	if t.ct != nil && n != nil {
		t.emitComments(t.ct.take(n))
	}
	if ft, ok := n.(nodeFormatter); ok {
		ft.Format(t, d)
	} else {
		debug.Dump(n)
	}
}

// emitComments prints attached comments at the boundary before their anchor
// node. A trailing comment continues the line before it (a line comment
// then forces the enclosing groups to break); an own-line or line comment
// goes on its own line; an inline block comment stays in the flow.
func (t *TrackedBuffer) emitComments(recs []commentRec) {
	for _, rec := range recs {
		switch {
		case rec.trailing:
			t.WriteString(" ")
			t.WriteString(rec.c.Text)
			if rec.c.Line() {
				t.Breaker()
			}
		case rec.c.OwnLine || rec.c.Line():
			t.Hardline()
			t.WriteString(rec.c.Text)
			t.Hardline()
		default:
			t.WriteString(rec.c.Text)
			t.WriteString(" ")
		}
	}
}

// beforeClause is the emission point for comments that sit between the
// previous clause and the one about to be printed, so they land before the
// clause keyword. Call it before the clause's line break. During the
// AttachComments dry run it records itself as a boundary marker, which is
// how classification and emission are guaranteed to pick the same spot —
// and how the author's own line break before the clause is observed and
// kept.
func (t *TrackedBuffer) beforeClause(n Node, d format.Dialect) {
	t.boundary(n)
}

// boundary is a separator-adjacent emission point (before a clause, after a
// comma, before an AND/OR): comments classified to the upcoming node print
// here, before the separator's line break, and a line break the author
// wrote at this boundary forces the break to stay. On the dry run it
// records a marker instead, so classification and emission agree by
// construction.
func (t *TrackedBuffer) boundary(next Node) {
	if next == nil {
		return
	}
	if t.anchors != nil {
		*t.anchors = append(*t.anchors, anchor{node: next, marker: true})
		return
	}
	if t.ct != nil {
		t.emitComments(t.ct.take(next))
		if t.ct.breakAt(next) {
			t.Breaker()
		}
	}
}

// flushRemaining prints every comment not yet printed; the statement is
// over, so everything left trails it.
func (t *TrackedBuffer) flushRemaining() {
	if t.ct == nil {
		return
	}
	for _, rec := range t.ct.takeRemaining() {
		if rec.trailing {
			t.WriteString(" ")
			t.WriteString(rec.c.Text)
			if rec.c.Line() {
				t.Breaker()
			}
		} else {
			t.Hardline()
			t.WriteString(rec.c.Text)
		}
	}
}

func (t *TrackedBuffer) join(n *List, d format.Dialect, sep string) {
	if n == nil {
		return
	}
	first := true
	for _, item := range n.Items {
		if _, ok := item.(*TODO); ok {
			continue
		}
		if !first {
			t.WriteString(sep)
		}
		first = false
		t.astFormat(item, d)
	}
}

// joinComma writes the list items separated by a comma and a break
// opportunity, so a list that fits stays on one line and a list that does
// not puts every item on its own line. An item's trailing comment prints
// after the comma, gofmt-style: `total, -- computed`.
func (t *TrackedBuffer) joinComma(n *List, d format.Dialect) {
	if n == nil {
		return
	}
	items := make([]Node, 0, len(n.Items))
	for _, item := range n.Items {
		if _, ok := item.(*TODO); ok {
			continue
		}
		items = append(items, item)
	}
	for i, item := range items {
		if i > 0 {
			t.WriteString(",")
			t.boundary(item)
			t.Line()
		}
		t.astFormat(item, d)
	}
}

// condition formats a clause-level boolean condition (WHERE, HAVING, ON).
// A top-level AND/OR chain is printed without the enclosing parentheses a
// nested BoolExpr gets — the clause keyword already delimits it — and each
// branch goes on its own line when the chain does not fit:
//
//	WHERE name LIKE $1
//	  AND id > $2
func (t *TrackedBuffer) condition(n Node, d format.Dialect) {
	be, ok := n.(*BoolExpr)
	if !ok || (be.Boolop != BoolExprTypeAnd && be.Boolop != BoolExprTypeOr) {
		t.astFormat(n, d)
		return
	}
	op := "AND "
	if be.Boolop == BoolExprTypeOr {
		op = "OR "
	}
	t.Group()
	t.Indent()
	first := true
	t.conditionArgs(be, d, op, &first)
	t.EndIndent()
	t.EndGroup()
}

// conditionArgs writes the branches of an AND/OR chain, flattening nested
// chains of the same operator (some parsers produce `a AND b AND c` as a
// left-nested tree) so no redundant parentheses appear.
func (t *TrackedBuffer) conditionArgs(be *BoolExpr, d format.Dialect, op string, first *bool) {
	if be.Args == nil {
		return
	}
	for _, item := range be.Args.Items {
		if _, ok := item.(*TODO); ok {
			continue
		}
		if nested, ok := item.(*BoolExpr); ok && nested.Boolop == be.Boolop {
			t.conditionArgs(nested, d, op, first)
			continue
		}
		if !*first {
			t.boundary(item)
			t.Line()
			t.WriteString(op)
		}
		*first = false
		t.astFormat(item, d)
	}
}

// Format renders the node as SQL on a single line.
func Format(n Node, d format.Dialect) string {
	return Pretty(n, d, -1)
}

// Pretty renders the node as SQL, breaking lines so the output fits within
// width columns where the statement's structure allows it. A negative width
// renders everything on a single line.
func Pretty(n Node, d format.Dialect, width int) string {
	return PrettyWithComments(n, d, width, nil)
}

// PrettyWithComments renders the node as SQL like Pretty, emitting the
// comments a prior AttachComments call anchored to its nodes and keeping
// the line breaks the author wrote at the boundaries the printer models.
// Each comment prints at the boundary before its anchor, a trailing comment
// continues the line of the code it followed, and every comment in the
// table is printed — anything left when the statement ends trails it.
func PrettyWithComments(n Node, d format.Dialect, width int, ct *CommentTable) string {
	tb := NewTrackedBuffer()
	tb.ct = ct
	if ft, ok := n.(nodeFormatter); ok {
		ft.Format(tb, d)
	}
	return tb.Print(width)
}

func set(n Node) bool {
	if n == nil {
		return false
	}
	_, ok := n.(*TODO)
	return !ok
}

func items(n *List) bool {
	if n == nil {
		return false
	}
	return len(n.Items) > 0
}

func todo(n *List) bool {
	for _, item := range n.Items {
		if _, ok := item.(*TODO); !ok {
			return false
		}
	}
	return true
}

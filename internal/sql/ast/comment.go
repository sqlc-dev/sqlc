package ast

import (
	"sort"

	"github.com/sqlc-dev/sqlc/internal/sql/format"
)

// File is a parsed query file: its statements together with the comments
// the parser's lexer saw. Engines whose parsers surface their trivia — the
// lexical channel of whitespace and comments the grammar never sees, in
// Roslyn's terminology — return it from ParseFile, so the formatter gets
// statements and comments from one lexer pass.
type File struct {
	Stmts    []Statement
	Comments []Comment
}

// Comment is a single SQL comment, positioned by byte offsets into the
// source the statement was parsed from — the same coordinates the engine
// parsers stamp on nodes.
type Comment struct {
	// Text is the comment as written, marker included ("-- x", "/* x */").
	Text  string
	Start int
	End   int
	// OwnLine reports that only blank space preceded the comment on its
	// line, so it leads the code after it rather than trailing the code
	// before it.
	OwnLine bool
}

// Line reports whether the comment runs to the end of its line (-- or #),
// so nothing may be printed after it on the same line.
func (c Comment) Line() bool {
	return len(c.Text) < 2 || c.Text[0] != '/' || c.Text[1] != '*'
}

// commentRec is one attached comment: the anchor nodes on either side of it
// in print order, and whether it trails the code before it (same source
// line) or leads the code after it (its own line).
//
// The next anchor is where the printer emits the comment; the prev anchor
// is the node the comment belongs to when the tree is edited — a trailing
// comment travels with the node it annotates, the way dave/dst attaches
// decorations in the Go ecosystem.
type commentRec struct {
	c          Comment
	prev, next Node
	trailing   bool
}

// CommentTable holds a statement's comments attached to its nodes. It is
// built once by AttachComments, before any printing (or, later, editing):
// each comment is classified against the statement's anchor nodes by source
// position and line, and from then on positions are never consulted again —
// the printer emits comments by node identity, which is what lets an edited
// or synthetic tree print its comments correctly.
type CommentTable struct {
	recs []commentRec
	// byNext indexes recs by their next anchor, in source order; end holds
	// the comments with no anchor after them, which trail the statement.
	byNext map[Node][]int
	end    []int
	taken  []bool
	nTaken int
}

// Exhausted reports whether every attached comment was printed.
func (t *CommentTable) Exhausted() bool {
	return t == nil || t.nTaken == len(t.recs)
}

// take returns and consumes the comments anchored to n, in order.
func (t *CommentTable) take(n Node) []commentRec {
	if t == nil || n == nil {
		return nil
	}
	idxs := t.byNext[n]
	if len(idxs) == 0 {
		return nil
	}
	out := make([]commentRec, 0, len(idxs))
	for _, i := range idxs {
		if t.taken[i] {
			continue
		}
		t.taken[i] = true
		t.nTaken++
		out = append(out, t.recs[i])
	}
	return out
}

// takeRemaining returns and consumes every comment not yet printed; the
// statement is over, so everything left trails it.
func (t *CommentTable) takeRemaining() []commentRec {
	if t == nil {
		return nil
	}
	var out []commentRec
	for i := range t.recs {
		if t.taken[i] {
			continue
		}
		t.taken[i] = true
		t.nTaken++
		out = append(out, t.recs[i])
	}
	return out
}

// AttachComments classifies a statement's interior comments against its
// nodes, producing the table the printer (and any future rewriting tool)
// works from. Placement follows gofmt's rules, decided here once from
// source positions and lines: a comment on the same line as the code before
// it trails that code; any other comment leads the first node printed after
// it; a comment after the last node trails the statement.
func AttachComments(raw *RawStmt, d format.Dialect, comments []Comment, src string) *CommentTable {
	anchors := collectAnchors(raw, d)
	lines := []int{0}
	for i := 0; i < len(src); i++ {
		if src[i] == '\n' {
			lines = append(lines, i+1)
		}
	}
	lineOf := func(pos int) int {
		return sort.SearchInts(lines, pos+1) - 1
	}

	table := &CommentTable{byNext: make(map[Node][]int, len(anchors))}
	for _, c := range comments {
		// prev: the printed node with the greatest position before the
		// comment (for the trailing/leading call). next: the emission point
		// for the first printed node after the comment — the node itself,
		// or, when boundary markers immediately precede it in print order,
		// the earliest of those markers, which is where the printer will
		// look for this comment first.
		var prev, next Node
		prevPos := -1
		for _, a := range anchors {
			if !a.marker && a.pos < c.Start && a.pos > prevPos {
				prev, prevPos = a.node, a.pos
			}
		}
		for i, a := range anchors {
			if a.marker || a.pos <= c.Start {
				continue
			}
			j := i
			for j > 0 && anchors[j-1].marker {
				j--
			}
			next = anchors[j].node
			break
		}
		rec := commentRec{
			c:        c,
			prev:     prev,
			next:     next,
			trailing: !c.OwnLine && prev != nil && lineOf(prevPos) == lineOf(c.Start),
		}
		table.recs = append(table.recs, rec)
		i := len(table.recs) - 1
		if next == nil {
			table.end = append(table.end, i)
		} else {
			table.byNext[next] = append(table.byNext[next], i)
		}
	}
	table.taken = make([]bool, len(table.recs))
	return table
}

type anchor struct {
	node Node
	pos  int
	// marker anchors are emission points (beforeClause, list boundaries)
	// rather than printed nodes; they carry no position of their own.
	marker bool
}

// collectAnchors renders the statement once, flat, recording every
// positioned node in the order the printer visits them. Print order is what
// comment emission is defined against, so classifying against it keeps the
// attach-time decision and the print-time emission point identical.
func collectAnchors(n Node, d format.Dialect) (out []anchor) {
	defer func() {
		// A formatter panic here surfaces later, on the real print; anchors
		// collected so far still place most comments.
		recover()
	}()
	tb := NewTrackedBuffer()
	tb.anchors = &out
	if ft, ok := n.(nodeFormatter); ok {
		ft.Format(tb, d)
	}
	return out
}


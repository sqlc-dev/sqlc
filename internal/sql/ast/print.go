package ast

import (
	"strings"
	"unicode/utf8"

	"github.com/sqlc-dev/sqlc/internal/debug"
	"github.com/sqlc-dev/sqlc/internal/sql/format"
)

type nodeFormatter interface {
	Format(*TrackedBuffer, format.Dialect)
}

// The printer is a small Wadler-style document renderer, the model behind
// Prettier and ruff (via Biome's printer). Format methods emit a stream of
// tokens: literal text, break opportunities (line, softline), and markers
// that open and close groups and indented regions. The renderer then lays
// out each group on a single line when its flat width fits within the line
// width, and otherwise breaks the group: its line and softline tokens
// become newlines followed by the current indentation, and nested groups
// are measured again against the space left on their own lines.
type tokenKind uint8

const (
	tokenText tokenKind = iota
	// tokenLine renders as a single space when its group fits on one line
	// and as a line break when the group is broken.
	tokenLine
	// tokenSoftline renders as nothing when its group fits on one line and
	// as a line break when the group is broken.
	tokenSoftline
	// tokenHardline always renders a line break and forces every group
	// containing it to break (it measures as infinitely wide). Emitted
	// around comments, which pin the author's layout open.
	tokenHardline
	// tokenBreaker renders nothing but forces every group containing it to
	// break. Emitted after a trailing line comment, whose text swallows the
	// rest of the line: the following separator must become a real break.
	tokenBreaker
	tokenOpenGroup
	tokenCloseGroup
	tokenOpenIndent
	tokenCloseIndent
)

// breakWidth is the measured width of tokens that force a break: wider than
// any real line, so no group containing one ever fits.
const breakWidth = 1 << 24

// indentWidth is the number of spaces per indentation level in broken
// (multi-line) output.
const indentWidth = 2

type token struct {
	kind tokenKind
	text string
}

type TrackedBuffer struct {
	tokens []token
	// cs, when set, carries the statement's interior comments; the buffer
	// flushes them by position as nodes are formatted, gofmt-style.
	cs *CommentSet
}

// NewTrackedBuffer creates a new TrackedBuffer.
func NewTrackedBuffer() *TrackedBuffer {
	return &TrackedBuffer{}
}

func (t *TrackedBuffer) WriteString(s string) (int, error) {
	if s == "" {
		return 0, nil
	}
	t.tokens = append(t.tokens, token{kind: tokenText, text: s})
	return len(s), nil
}

// Write implements io.Writer so Format methods can use fmt.Fprintf.
func (t *TrackedBuffer) Write(p []byte) (int, error) {
	return t.WriteString(string(p))
}

func (t *TrackedBuffer) WriteRune(r rune) (int, error) {
	return t.WriteString(string(r))
}

// String renders the buffer on a single line.
func (t *TrackedBuffer) String() string {
	return t.print(-1)
}

// line marks a break opportunity between two tokens: a space when the
// surrounding group fits on one line, a line break when it does not.
func (t *TrackedBuffer) line() {
	t.tokens = append(t.tokens, token{kind: tokenLine})
}

// softline marks a break opportunity that disappears entirely when the
// surrounding group fits on one line.
func (t *TrackedBuffer) softline() {
	t.tokens = append(t.tokens, token{kind: tokenSoftline})
}

func (t *TrackedBuffer) hardline() {
	t.tokens = append(t.tokens, token{kind: tokenHardline})
}

func (t *TrackedBuffer) breaker() {
	t.tokens = append(t.tokens, token{kind: tokenBreaker})
}

// group opens a region the renderer tries to lay out on a single line,
// breaking it only when its flat form does not fit. Must be paired with
// endGroup.
func (t *TrackedBuffer) group() {
	t.tokens = append(t.tokens, token{kind: tokenOpenGroup})
}

func (t *TrackedBuffer) endGroup() {
	t.tokens = append(t.tokens, token{kind: tokenCloseGroup})
}

// indent opens a region printed one indentation level deeper when line
// breaks occur inside it. Must be paired with endIndent.
func (t *TrackedBuffer) indent() {
	t.tokens = append(t.tokens, token{kind: tokenOpenIndent})
}

func (t *TrackedBuffer) endIndent() {
	t.tokens = append(t.tokens, token{kind: tokenCloseIndent})
}

func (t *TrackedBuffer) astFormat(n Node, d format.Dialect) {
	if t.cs != nil && n != nil {
		if pos := n.Pos(); pos > 0 {
			t.flushLeading(pos)
			t.cs.advance(pos)
		}
	}
	if ft, ok := n.(nodeFormatter); ok {
		ft.Format(t, d)
	} else {
		debug.Dump(n)
	}
}

// flushLeading prints every unprinted comment positioned before pos, in
// leading style: own-line comments (and line comments, which cannot share
// their line with anything after them) go on their own line; an inline
// block comment stays in the flow.
func (t *TrackedBuffer) flushLeading(pos int) {
	for {
		c, ok := t.cs.pending()
		if !ok || c.Start >= pos {
			return
		}
		t.cs.next++
		if c.OwnLine || c.Line() {
			t.hardline()
			t.WriteString(c.Text)
			t.hardline()
		} else {
			t.WriteString(c.Text)
			t.WriteString(" ")
		}
	}
}

// flushTrailing prints comments that trail the code just printed: not on
// their own line, on the same source line as the printer's cursor, and
// positioned before limit (the next element, so a comment trailing a later
// element on the same line is not stolen). A trailing line comment swallows
// the rest of its output line, so it forces the enclosing groups to break.
func (t *TrackedBuffer) flushTrailing(limit int) {
	if t.cs == nil {
		return
	}
	for {
		c, ok := t.cs.pending()
		if !ok || c.Start >= limit || c.OwnLine {
			return
		}
		if t.cs.cursor <= 0 || t.cs.lineOf(c.Start) != t.cs.lineOf(t.cs.cursor) {
			return
		}
		t.cs.next++
		t.WriteString(" ")
		t.WriteString(c.Text)
		if c.Line() {
			t.breaker()
		}
	}
}

// beforeClause flushes the comments that sit between the previous clause
// and the one about to be printed, so they land before the clause keyword.
// Call it before the clause's line break.
func (t *TrackedBuffer) beforeClause(n Node) {
	if t.cs == nil || n == nil {
		return
	}
	pos := firstPos(n)
	if pos <= 0 {
		return
	}
	t.flushTrailing(pos)
	t.flushLeading(pos)
}

// flushRemaining prints every comment not yet printed; the statement is
// over, so everything left trails it.
func (t *TrackedBuffer) flushRemaining() {
	if t.cs == nil {
		return
	}
	for {
		c, ok := t.cs.pending()
		if !ok {
			return
		}
		t.cs.next++
		if !c.OwnLine && t.cs.cursor > 0 && t.cs.lineOf(c.Start) == t.cs.lineOf(t.cs.cursor) {
			t.WriteString(" ")
			t.WriteString(c.Text)
			if c.Line() {
				t.breaker()
			}
		} else {
			t.hardline()
			t.WriteString(c.Text)
		}
	}
}

// firstPos returns a node's position, descending into lists to the first
// positioned item (lists themselves carry no position).
func firstPos(n Node) int {
	if n == nil {
		return 0
	}
	if l, ok := n.(*List); ok {
		for _, item := range l.Items {
			if _, ok := item.(*TODO); ok {
				continue
			}
			if p := firstPos(item); p > 0 {
				return p
			}
		}
		return 0
	}
	return n.Pos()
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
			t.flushTrailing(firstPos(item))
			t.line()
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
	t.group()
	t.indent()
	first := true
	t.conditionArgs(be, d, op, &first)
	t.endIndent()
	t.endGroup()
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
			t.flushTrailing(firstPos(item))
			t.line()
			t.WriteString(op)
		}
		*first = false
		t.astFormat(item, d)
	}
}

// docNode is the tree form of the token stream: text and break tokens are
// leaves, groups and indented regions are interior nodes.
type docNode struct {
	kind  tokenKind
	text  string
	kids  []*docNode
	width int // flat width in runes, filled in by measure
}

// tree assembles the flat token stream into a document tree. Stray close
// markers are ignored and unclosed regions end at the end of the stream.
func (t *TrackedBuffer) tree() *docNode {
	root := &docNode{kind: tokenOpenGroup}
	stack := []*docNode{root}
	for _, tok := range t.tokens {
		top := stack[len(stack)-1]
		switch tok.kind {
		case tokenText, tokenLine, tokenSoftline, tokenHardline, tokenBreaker:
			top.kids = append(top.kids, &docNode{kind: tok.kind, text: tok.text})
		case tokenOpenGroup, tokenOpenIndent:
			n := &docNode{kind: tok.kind}
			top.kids = append(top.kids, n)
			stack = append(stack, n)
		case tokenCloseGroup, tokenCloseIndent:
			if len(stack) > 1 {
				stack = stack[:len(stack)-1]
			}
		}
	}
	return root
}

// measure computes the width of every node when rendered flat. Hardlines
// and breakers measure wider than any real line, so no group containing one
// ever fits — Prettier's breakParent, by arithmetic.
func measure(d *docNode) int {
	switch d.kind {
	case tokenText:
		d.width = utf8.RuneCountInString(d.text)
	case tokenLine:
		d.width = 1
	case tokenSoftline:
		d.width = 0
	case tokenHardline, tokenBreaker:
		d.width = breakWidth
	default:
		for _, kid := range d.kids {
			d.width += measure(kid)
			if d.width > breakWidth {
				d.width = breakWidth
			}
		}
	}
	return d.width
}

// print renders the token stream. A negative width renders everything on
// one line; otherwise each group is rendered flat only when it fits in the
// space remaining on its line.
func (t *TrackedBuffer) print(width int) string {
	root := t.tree()
	measure(root)

	type frame struct {
		d      *docNode
		indent int
		flat   bool
	}
	var sb strings.Builder
	col := 0
	// atLineStart dedupes consecutive breaks: a break that fires when the
	// output is already at the start of a fresh line writes nothing, so a
	// comment's surrounding hardlines collapse against clause breaks
	// instead of leaving blank lines.
	atLineStart := true
	newline := func(indent int) {
		if atLineStart {
			return
		}
		sb.WriteByte('\n')
		for i := 0; i < indent; i++ {
			sb.WriteByte(' ')
		}
		col = indent
		atLineStart = true
	}
	stack := []frame{{d: root, flat: width < 0}}
	for len(stack) > 0 {
		f := stack[len(stack)-1]
		stack = stack[:len(stack)-1]
		switch f.d.kind {
		case tokenText:
			sb.WriteString(f.d.text)
			col += f.d.width
			atLineStart = false
		case tokenLine:
			if f.flat {
				sb.WriteByte(' ')
				col++
				atLineStart = false
			} else {
				newline(f.indent)
			}
		case tokenSoftline:
			if !f.flat {
				newline(f.indent)
			}
		case tokenHardline:
			newline(f.indent)
		case tokenBreaker:
			// forces enclosing groups broken; renders nothing
		case tokenOpenGroup:
			flat := f.flat || f.d.width <= width-col
			for i := len(f.d.kids) - 1; i >= 0; i-- {
				stack = append(stack, frame{d: f.d.kids[i], indent: f.indent, flat: flat})
			}
		case tokenOpenIndent:
			for i := len(f.d.kids) - 1; i >= 0; i-- {
				stack = append(stack, frame{d: f.d.kids[i], indent: f.indent + indentWidth, flat: f.flat})
			}
		}
	}
	return sb.String()
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

// PrettyWithComments renders the node as SQL like Pretty, weaving the
// statement's interior comments back in at the positions their offsets
// anchor them to. The printer flushes each comment before the first node
// positioned after it; a comment on the same source line as the code before
// it trails that code. Every comment in the set is printed: anything left
// when the statement ends trails it.
func PrettyWithComments(n Node, d format.Dialect, width int, cs *CommentSet) string {
	tb := NewTrackedBuffer()
	tb.cs = cs
	if ft, ok := n.(nodeFormatter); ok {
		ft.Format(tb, d)
	}
	return tb.print(width)
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

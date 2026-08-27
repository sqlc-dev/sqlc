// Package printer implements the document renderer the SQL formatter
// prints through: a small Wadler-style engine, the model behind Prettier
// and ruff (via Biome's printer). Callers emit a stream of tokens — literal
// text, break opportunities, and markers that open and close groups and
// indented regions — and Print lays each group out on a single line when
// its flat width fits, breaking it otherwise. The package knows nothing
// about SQL or the AST; the ast package's TrackedBuffer embeds a Buffer and
// layers node formatting and comment emission on top.
package printer

import (
	"strings"
	"unicode/utf8"
)

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
	// rest of the line, and at boundaries where the author broke the line.
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

// Buffer accumulates the document token stream.
type Buffer struct {
	tokens []token
}

func (t *Buffer) WriteString(s string) (int, error) {
	if s == "" {
		return 0, nil
	}
	t.tokens = append(t.tokens, token{kind: tokenText, text: s})
	return len(s), nil
}

// Write implements io.Writer so callers can use fmt.Fprintf.
func (t *Buffer) Write(p []byte) (int, error) {
	return t.WriteString(string(p))
}

func (t *Buffer) WriteRune(r rune) (int, error) {
	return t.WriteString(string(r))
}

// String renders the buffer on a single line.
func (t *Buffer) String() string {
	return t.Print(-1)
}

// Line marks a break opportunity between two tokens: a space when the
// surrounding group fits on one line, a line break when it does not.
func (t *Buffer) Line() {
	t.tokens = append(t.tokens, token{kind: tokenLine})
}

// Softline marks a break opportunity that disappears entirely when the
// surrounding group fits on one line.
func (t *Buffer) Softline() {
	t.tokens = append(t.tokens, token{kind: tokenSoftline})
}

// Hardline always renders a line break and forces every enclosing group to
// break.
func (t *Buffer) Hardline() {
	t.tokens = append(t.tokens, token{kind: tokenHardline})
}

// Breaker renders nothing but forces every enclosing group to break, so the
// next break opportunity becomes a real line break.
func (t *Buffer) Breaker() {
	t.tokens = append(t.tokens, token{kind: tokenBreaker})
}

// Group opens a region the renderer tries to lay out on a single line,
// breaking it only when its flat form does not fit. Must be paired with
// EndGroup.
func (t *Buffer) Group() {
	t.tokens = append(t.tokens, token{kind: tokenOpenGroup})
}

func (t *Buffer) EndGroup() {
	t.tokens = append(t.tokens, token{kind: tokenCloseGroup})
}

// Indent opens a region printed one indentation level deeper when line
// breaks occur inside it. Must be paired with EndIndent.
func (t *Buffer) Indent() {
	t.tokens = append(t.tokens, token{kind: tokenOpenIndent})
}

func (t *Buffer) EndIndent() {
	t.tokens = append(t.tokens, token{kind: tokenCloseIndent})
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
func (t *Buffer) tree() *docNode {
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

// Print renders the token stream. A negative width renders everything on
// one line; otherwise each group is rendered flat only when it fits in the
// space remaining on its line — and a group holding a hardline or breaker
// never fits, whatever the width, so forced breaks survive even an
// unlimited line length.
func (t *Buffer) Print(width int) string {
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
			flat := f.flat || (f.d.width < breakWidth && f.d.width <= width-col)
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

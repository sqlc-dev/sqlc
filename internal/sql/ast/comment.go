package ast

import (
	"sort"
	"strings"
)

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
	return !strings.HasPrefix(c.Text, "/*")
}

// CommentSet carries a statement's interior comments through the printer,
// gofmt-style: the printer holds a cursor over the comment list and flushes
// every comment positioned before the node it is about to print.
type CommentSet struct {
	comments []Comment
	// lineStarts holds the byte offset of the first byte of every line of
	// the source, sqlc's stand-in for go/token.FileSet: it turns a byte
	// offset into a line number, which is what placement decisions compare.
	lineStarts []int

	next   int // index of the next unprinted comment
	cursor int // greatest node position printed so far
}

// NewCommentSet prepares comments for printing against source text. Only
// comments inside the printed statement should be included; comments above
// the statement and after its terminator are the caller's to keep.
func NewCommentSet(comments []Comment, src string) *CommentSet {
	cs := &CommentSet{comments: comments, lineStarts: []int{0}}
	for i := 0; i < len(src); i++ {
		if src[i] == '\n' {
			cs.lineStarts = append(cs.lineStarts, i+1)
		}
	}
	sort.SliceStable(cs.comments, func(i, j int) bool {
		return cs.comments[i].Start < cs.comments[j].Start
	})
	return cs
}

func (cs *CommentSet) lineOf(pos int) int {
	return sort.SearchInts(cs.lineStarts, pos+1) - 1
}

// advance moves the cursor forward to pos; the cursor never moves back.
func (cs *CommentSet) advance(pos int) {
	if pos > cs.cursor {
		cs.cursor = pos
	}
}

// pending returns the next unprinted comment, if any.
func (cs *CommentSet) pending() (Comment, bool) {
	if cs.next >= len(cs.comments) {
		return Comment{}, false
	}
	return cs.comments[cs.next], true
}

// Exhausted reports whether every comment was printed.
func (cs *CommentSet) Exhausted() bool {
	return cs.next >= len(cs.comments)
}

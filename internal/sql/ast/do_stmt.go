package ast

import "github.com/sqlc-dev/sqlc/internal/sql/format"

type DoStmt struct {
	Args *List
}

func (n *DoStmt) Pos() int {
	return 0
}

func (n *DoStmt) Format(buf *TrackedBuffer, d format.Dialect) {
	if n == nil {
		return
	}
	buf.WriteString("DO ")
	// Find the "as" argument which contains the body
	if items(n.Args) {
		for _, arg := range n.Args.Items {
			if de, ok := arg.(*DefElem); ok && de.Defname != nil && *de.Defname == "as" {
				if s, ok := de.Arg.(*String); ok {
					buf.WriteString("$$")
					buf.WriteString(s.Str)
					buf.WriteString("$$")
				}
			}
		}
	}
}

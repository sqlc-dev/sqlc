package ast

import "github.com/sqlc-dev/sqlc/internal/sql/format"

type MultiAssignRef struct {
	Tag NodeTag[MultiAssignRef] `json:"tag"`

	Source   Node `json:"source,omitempty"`
	Colno    int  `json:"colno"`
	Ncolumns int  `json:"ncolumns"`
}

func (n *MultiAssignRef) Pos() int {
	return 0
}

func (n *MultiAssignRef) Format(buf *TrackedBuffer, d format.Dialect) {
	if n == nil {
		return
	}
	buf.astFormat(n.Source, d)
}

package ast

import (
	"strings"

	"github.com/sqlc-dev/sqlc/internal/sql/format"
)

type ColumnRef struct {
	Name string

	// From pg.ColumnRef
	Fields   *List
	Location int
}

func (n *ColumnRef) Pos() int {
	return n.Location
}

func (n *ColumnRef) Format(buf *TrackedBuffer, d format.Dialect) {
	if n == nil {
		return
	}

	if n.Fields != nil {
		var items []string
		for _, item := range n.Fields.Items {
			switch nn := item.(type) {
			case *String:
				items = append(items, d.QuoteIdent(nn.Str))
			case *A_Star:
				items = append(items, "*")
			}
		}
		buf.WriteString(strings.Join(items, "."))
	}
}

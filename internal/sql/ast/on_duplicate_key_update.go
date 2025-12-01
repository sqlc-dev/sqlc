package ast

import "github.com/sqlc-dev/sqlc/internal/sql/format"

// OnDuplicateKeyUpdate represents MySQL's ON DUPLICATE KEY UPDATE clause
type OnDuplicateKeyUpdate struct {
	// TargetList contains the assignments (column = value pairs)
	TargetList *List
	Location   int
}

func (n *OnDuplicateKeyUpdate) Pos() int {
	return n.Location
}

func (n *OnDuplicateKeyUpdate) Format(buf *TrackedBuffer, d format.Dialect) {
	if n == nil {
		return
	}
	buf.WriteString("ON DUPLICATE KEY UPDATE ")
	if n.TargetList != nil {
		for i, item := range n.TargetList.Items {
			if i > 0 {
				buf.WriteString(", ")
			}
			if rt, ok := item.(*ResTarget); ok {
				if rt.Name != nil {
					buf.WriteString(*rt.Name)
				}
				buf.WriteString(" = ")
				buf.astFormat(rt.Val, d)
			} else {
				buf.astFormat(item, d)
			}
		}
	}
}

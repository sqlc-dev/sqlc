package ast

import "strings"

type ColumnRef struct {
	Name string

	// From pg.ColumnRef
	Fields   *List
	Location int
}

func (n *ColumnRef) Pos() int {
	return n.Location
}

func (n *ColumnRef) Format(buf *TrackedBuffer) {
	if n == nil {
		return
	}

	if n.Fields != nil {
		var items []string
		for _, item := range n.Fields.Items {
			switch nn := item.(type) {
			case *String:
				if nn.Str == "user" {
					items = append(items, `"user"`)
				} else {
					items = append(items, nn.Str)
				}
			case *A_Star:
				items = append(items, "*")
			}
		}
		buf.WriteString(strings.Join(items, "."))
	}
}

package ast

import "github.com/sqlc-dev/sqlc/internal/sql/format"

// RangePivot is a FROM item that rotates the rows of another, as SQL
// Server's PIVOT and UNPIVOT do. PIVOT turns the values of ForColumn named
// in InList into columns of Aggregate; UNPIVOT turns the columns in InList
// into rows of ValueColumn, named in ForColumn of type ForType.
type RangePivot struct {
	Tag NodeTag[RangePivot] `json:"tag"`

	Source      Node      `json:"source,omitempty"`
	Unpivot     bool      `json:"unpivot"`
	Aggregate   *FuncCall `json:"aggregate,omitempty"`
	ValueColumn string    `json:"value_column"`
	ForColumn   string    `json:"for_column"`
	ForType     *TypeName `json:"for_type,omitempty"`
	InList      *List     `json:"in_list,omitempty"`
	Alias       *Alias    `json:"alias,omitempty"`
}

func (n *RangePivot) Pos() int {
	return 0
}

func (n *RangePivot) Format(buf *TrackedBuffer, d format.Dialect) {
	if n == nil {
		return
	}
	buf.astFormat(n.Source, d)
	if n.Unpivot {
		buf.WriteString(" UNPIVOT (")
		buf.WriteString(d.QuoteIdent(n.ValueColumn))
	} else {
		buf.WriteString(" PIVOT (")
		buf.astFormat(n.Aggregate, d)
	}
	buf.WriteString(" FOR ")
	buf.WriteString(d.QuoteIdent(n.ForColumn))
	buf.WriteString(" IN (")
	if n.InList != nil {
		for i, item := range n.InList.Items {
			if i > 0 {
				buf.WriteString(", ")
			}
			if s, ok := item.(*String); ok {
				buf.WriteString(d.QuoteIdent(s.Str))
			}
		}
	}
	buf.WriteString("))")
	if n.Alias != nil {
		buf.WriteString(" AS ")
		buf.astFormat(n.Alias, d)
	}
}

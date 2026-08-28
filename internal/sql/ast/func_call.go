package ast

import "github.com/sqlc-dev/sqlc/internal/sql/format"

type FuncCall struct {
	Tag NodeTag[FuncCall] `json:"tag"`

	Func           *FuncName  `json:"func,omitempty"`
	Funcname       *List      `json:"funcname,omitempty"`
	Args           *List      `json:"args,omitempty"`
	AggOrder       *List      `json:"agg_order,omitempty"`
	AggFilter      Node       `json:"agg_filter,omitempty"`
	AggWithinGroup bool       `json:"agg_within_group"`
	AggStar        bool       `json:"agg_star"`
	AggDistinct    bool       `json:"agg_distinct"`
	FuncVariadic   bool       `json:"func_variadic"`
	Over           *WindowDef `json:"over,omitempty"`
	Separator      *string    `json:"separator,omitempty"` // MySQL GROUP_CONCAT SEPARATOR
	Location       int        `json:"location"`
}

func (n *FuncCall) Pos() int {
	return n.Location
}

func (n *FuncCall) Format(buf *TrackedBuffer, d format.Dialect) {
	if n == nil {
		return
	}
	buf.astFormat(n.Func, d)
	buf.WriteString("(")
	if n.AggDistinct {
		buf.WriteString("DISTINCT ")
	}
	if n.AggStar {
		buf.WriteString("*")
	} else {
		buf.astFormat(n.Args, d)
	}
	// ORDER BY inside function call (not WITHIN GROUP)
	if items(n.AggOrder) && !n.AggWithinGroup {
		buf.WriteString(" ORDER BY ")
		buf.join(n.AggOrder, d, ", ")
	}
	// SEPARATOR for GROUP_CONCAT (MySQL)
	if n.Separator != nil {
		buf.WriteString(" SEPARATOR ")
		buf.WriteString("'")
		buf.WriteString(*n.Separator)
		buf.WriteString("'")
	}
	buf.WriteString(")")
	// WITHIN GROUP clause for ordered-set aggregates
	if items(n.AggOrder) && n.AggWithinGroup {
		buf.WriteString(" WITHIN GROUP (ORDER BY ")
		buf.join(n.AggOrder, d, ", ")
		buf.WriteString(")")
	}
	if set(n.AggFilter) {
		buf.WriteString(" FILTER (WHERE ")
		buf.astFormat(n.AggFilter, d)
		buf.WriteString(")")
	}
	if n.Over != nil {
		buf.WriteString(" OVER ")
		buf.astFormat(n.Over, d)
	}
}

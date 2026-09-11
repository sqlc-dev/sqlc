package ast

import "github.com/sqlc-dev/sqlc/internal/sql/format"

type ParamRef struct {
	Tag NodeTag[ParamRef] `json:"tag"`

	Number   int  `json:"number"`
	Location int  `json:"location"`
	Dollar   bool `json:"dollar"`
	// Name is the name the query gave the placeholder, when its syntax
	// has one: ClickHouse's {name:Type}.
	Name string `json:"name,omitempty"`
}

func (n *ParamRef) Pos() int {
	return n.Location
}

func (n *ParamRef) Format(buf *TrackedBuffer, d format.Dialect) {
	if n == nil {
		return
	}
	buf.WriteString(d.Param(n.Number, n.Dollar))
}

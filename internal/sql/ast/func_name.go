package ast

import "github.com/sqlc-dev/sqlc/internal/sql/format"

type FuncName struct {
	Tag NodeTag[FuncName] `json:"tag"`

	Catalog string `json:"catalog"`
	Schema  string `json:"schema"`
	Name    string `json:"name"`
}

func (n *FuncName) Pos() int {
	return 0
}

func (n *FuncName) Format(buf *TrackedBuffer, d format.Dialect) {
	if n == nil {
		return
	}
	if n.Schema != "" {
		buf.WriteString(n.Schema)
		buf.WriteString(".")
	}
	if n.Name != "" {
		buf.WriteString(n.Name)
	}
}

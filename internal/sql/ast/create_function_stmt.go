package ast

import "github.com/sqlc-dev/sqlc/internal/sql/format"

type CreateFunctionStmt struct {
	Replace    bool
	Params     *List
	ReturnType *TypeName
	Func       *FuncName
	// TODO: Understand these two fields
	Options    *List
	WithClause *List
}

func (n *CreateFunctionStmt) Pos() int {
	return 0
}

func (n *CreateFunctionStmt) Format(buf *TrackedBuffer, d format.Dialect) {
	if n == nil {
		return
	}
	buf.WriteString("CREATE ")
	if n.Replace {
		buf.WriteString("OR REPLACE ")
	}
	buf.WriteString("FUNCTION ")
	buf.astFormat(n.Func, d)
	buf.WriteString("(")
	if items(n.Params) {
		buf.join(n.Params, d, ", ")
	}
	buf.WriteString(")")
	if n.ReturnType != nil {
		buf.WriteString(" RETURNS ")
		buf.astFormat(n.ReturnType, d)
	}
	// Format options (AS, LANGUAGE, etc.)
	if items(n.Options) {
		for _, opt := range n.Options.Items {
			buf.WriteString(" ")
			buf.astFormat(opt, d)
		}
	}
}

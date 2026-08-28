package ast

import "github.com/sqlc-dev/sqlc/internal/sql/format"

type CreateTableStmt struct {
	Tag NodeTag[CreateTableStmt] `json:"tag"`

	IfNotExists bool         `json:"if_not_exists"`
	Name        *TableName   `json:"name,omitempty"`
	Cols        []*ColumnDef `json:"cols,omitempty"`
	ReferTable  *TableName   `json:"refer_table,omitempty"`
	Comment     string       `json:"comment"`
	Inherits    []*TableName `json:"inherits,omitempty"`
	// Using names the module of a virtual table (SQLite's CREATE VIRTUAL
	// TABLE ... USING module(...)), and ModuleArgs carries the module's
	// argument list as written — its grammar belongs to the module, so the
	// arguments pass through verbatim; nil means the declaration had no
	// argument list at all. Cols still holds the columns sqlc derives from
	// the arguments for the catalog; the statement prints as its
	// declaration, not as those columns.
	Using      string   `json:"using"`
	ModuleArgs []string `json:"module_args,omitempty"`
	// ModuleArgsMultiline records that the declaration was written across
	// lines. The arguments carry no positions, so the printer cannot keep
	// the author's exact breaks and prints the canonical broken form — one
	// argument per line — instead.
	ModuleArgsMultiline bool `json:"module_args_multiline"`
	// Incomplete marks a statement whose source carried syntax this node
	// does not model — column or table constraints beyond plain NOT NULL
	// and PRIMARY KEY, table options, TEMP, or an AS SELECT body. No
	// faithful rendering exists for it.
	Incomplete bool `json:"incomplete"`
}

func (n *CreateTableStmt) Pos() int {
	return 0
}

func (n *CreateTableStmt) Format(buf *TrackedBuffer, d format.Dialect) {
	if n == nil {
		return
	}
	// A virtual table prints as its declaration: the module name and the
	// argument list as written.
	if n.Using != "" {
		buf.WriteString("CREATE VIRTUAL TABLE ")
		if n.IfNotExists {
			buf.WriteString("IF NOT EXISTS ")
		}
		buf.astFormat(n.Name, d)
		buf.WriteString(" USING ")
		buf.WriteString(n.Using)
		if n.ModuleArgs != nil {
			buf.WriteString("(")
			buf.Group()
			buf.Indent()
			if n.ModuleArgsMultiline {
				buf.Breaker()
			}
			buf.Softline()
			for i, arg := range n.ModuleArgs {
				if i > 0 {
					buf.WriteString(",")
					buf.Line()
				}
				buf.WriteString(arg)
			}
			buf.EndIndent()
			buf.Softline()
			buf.EndGroup()
			buf.WriteString(")")
		}
		return
	}
	// An incomplete statement cannot be printed back: part of its source
	// was parsed away. Render nothing, which no verification accepts, so
	// the formatter keeps the statement as written.
	if n.Incomplete {
		return
	}
	buf.WriteString("CREATE TABLE ")
	if n.IfNotExists {
		buf.WriteString("IF NOT EXISTS ")
	}
	buf.astFormat(n.Name, d)

	buf.WriteString(" (")
	buf.Group()
	buf.Indent()
	if len(n.Cols) > 0 {
		buf.boundary(n.Cols[0])
	}
	buf.Softline()
	for i, col := range n.Cols {
		if i > 0 {
			buf.WriteString(",")
			buf.boundary(col)
			buf.Line()
		}
		buf.astFormat(col, d)
	}
	buf.EndIndent()
	buf.Softline()
	buf.EndGroup()
	buf.WriteString(")")
}

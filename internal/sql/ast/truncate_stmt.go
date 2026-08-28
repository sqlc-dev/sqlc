package ast

import "github.com/sqlc-dev/sqlc/internal/sql/format"

type TruncateStmt struct {
	Tag NodeTag[TruncateStmt] `json:"tag"`

	Relations   *List        `json:"relations,omitempty"`
	RestartSeqs bool         `json:"restart_seqs"`
	Behavior    DropBehavior `json:"behavior"`
}

func (n *TruncateStmt) Pos() int {
	return 0
}

func (n *TruncateStmt) Format(buf *TrackedBuffer, d format.Dialect) {
	if n == nil {
		return
	}
	buf.WriteString("TRUNCATE ")
	buf.astFormat(n.Relations, d)
}

package ast

import (
	"github.com/kyleconroy/sqlc/internal/sql/ast"
)

type RawStmt struct {
	Stmt         Node
	StmtLocation int
	StmtLen      int
}

func (n *RawStmt) Pos() int {
	return 0
}

package pg

import (
	"github.com/kyleconroy/sqlc/internal/sql/ast"
)

type CreateStmt struct {
	Relation       *RangeVar
	TableElts      *ast.List
	InhRelations   *ast.List
	Partbound      *PartitionBoundSpec
	Partspec       *PartitionSpec
	OfTypename     *TypeName
	Constraints    *ast.List
	Options        *ast.List
	Oncommit       OnCommitAction
	Tablespacename *string
	IfNotExists    bool
}

func (n *CreateStmt) Pos() int {
	return 0
}

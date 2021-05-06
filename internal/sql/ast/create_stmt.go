package ast

type CreateStmt struct {
	Relation       *RangeVar
	TableElts      *List
	InhRelations   *List
	Partbound      *PartitionBoundSpec
	Partspec       *PartitionSpec
	OfTypename     *TypeName
	Constraints    *List
	Options        *List
	Oncommit       OnCommitAction
	Tablespacename *string
	IfNotExists    bool
}

func (n *CreateStmt) Pos() int {
	return 0
}

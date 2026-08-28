package ast

type CreateStmt struct {
	Tag NodeTag[CreateStmt] `json:"tag"`

	Relation       *RangeVar           `json:",omitempty"`
	TableElts      *List               `json:",omitempty"`
	InhRelations   *List               `json:",omitempty"`
	Partbound      *PartitionBoundSpec `json:",omitempty"`
	Partspec       *PartitionSpec      `json:",omitempty"`
	OfTypename     *TypeName           `json:",omitempty"`
	Constraints    *List               `json:",omitempty"`
	Options        *List               `json:",omitempty"`
	Oncommit       OnCommitAction
	Tablespacename *string `json:",omitempty"`
	IfNotExists    bool
}

func (n *CreateStmt) Pos() int {
	return 0
}

package ast

type CreateStmt struct {
	Tag NodeTag[CreateStmt] `json:"tag"`

	Relation       *RangeVar           `json:"relation,omitempty"`
	TableElts      *List               `json:"table_elts,omitempty"`
	InhRelations   *List               `json:"inh_relations,omitempty"`
	Partbound      *PartitionBoundSpec `json:"partbound,omitempty"`
	Partspec       *PartitionSpec      `json:"partspec,omitempty"`
	OfTypename     *TypeName           `json:"of_typename,omitempty"`
	Constraints    *List               `json:"constraints,omitempty"`
	Options        *List               `json:"options,omitempty"`
	Oncommit       OnCommitAction      `json:"oncommit"`
	Tablespacename *string             `json:"tablespacename,omitempty"`
	IfNotExists    bool                `json:"if_not_exists"`
}

func (n *CreateStmt) Pos() int {
	return 0
}

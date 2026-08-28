package ast

type ReplicaIdentityStmt struct {
	Tag NodeTag[ReplicaIdentityStmt] `json:"tag"`

	IdentityType byte
	Name         *string
}

func (n *ReplicaIdentityStmt) Pos() int {
	return 0
}

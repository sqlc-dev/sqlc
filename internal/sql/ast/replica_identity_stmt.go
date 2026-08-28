package ast

type ReplicaIdentityStmt struct {
	Tag NodeTag[ReplicaIdentityStmt] `json:"tag"`

	IdentityType byte    `json:"identity_type"`
	Name         *string `json:"name,omitempty"`
}

func (n *ReplicaIdentityStmt) Pos() int {
	return 0
}

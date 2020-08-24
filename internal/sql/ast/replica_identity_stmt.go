package pg

type ReplicaIdentityStmt struct {
	IdentityType byte
	Name         *string
}

func (n *ReplicaIdentityStmt) Pos() int {
	return 0
}

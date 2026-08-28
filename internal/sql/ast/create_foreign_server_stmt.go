package ast

type CreateForeignServerStmt struct {
	Tag NodeTag[CreateForeignServerStmt] `json:"tag"`

	Servername  *string
	Servertype  *string
	Version     *string
	Fdwname     *string
	IfNotExists bool
	Options     *List
}

func (n *CreateForeignServerStmt) Pos() int {
	return 0
}

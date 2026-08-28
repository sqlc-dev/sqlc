package ast

type CreateForeignServerStmt struct {
	Tag NodeTag[CreateForeignServerStmt] `json:"tag"`

	Servername  *string `json:",omitempty"`
	Servertype  *string `json:",omitempty"`
	Version     *string `json:",omitempty"`
	Fdwname     *string `json:",omitempty"`
	IfNotExists bool
	Options     *List `json:",omitempty"`
}

func (n *CreateForeignServerStmt) Pos() int {
	return 0
}

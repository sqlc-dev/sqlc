package ast

type CreateForeignServerStmt struct {
	Tag NodeTag[CreateForeignServerStmt] `json:"tag"`

	Servername  *string `json:"servername,omitempty"`
	Servertype  *string `json:"servertype,omitempty"`
	Version     *string `json:"version,omitempty"`
	Fdwname     *string `json:"fdwname,omitempty"`
	IfNotExists bool    `json:"if_not_exists"`
	Options     *List   `json:"options,omitempty"`
}

func (n *CreateForeignServerStmt) Pos() int {
	return 0
}

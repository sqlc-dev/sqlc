package ast

type CreatePLangStmt struct {
	Tag NodeTag[CreatePLangStmt] `json:"tag"`

	Replace     bool
	Plname      *string
	Plhandler   *List
	Plinline    *List
	Plvalidator *List
	Pltrusted   bool
}

func (n *CreatePLangStmt) Pos() int {
	return 0
}

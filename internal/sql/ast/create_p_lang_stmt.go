package ast

type CreatePLangStmt struct {
	Tag NodeTag[CreatePLangStmt] `json:"tag"`

	Replace     bool
	Plname      *string `json:",omitempty"`
	Plhandler   *List   `json:",omitempty"`
	Plinline    *List   `json:",omitempty"`
	Plvalidator *List   `json:",omitempty"`
	Pltrusted   bool
}

func (n *CreatePLangStmt) Pos() int {
	return 0
}

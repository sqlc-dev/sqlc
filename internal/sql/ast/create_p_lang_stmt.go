package ast

type CreatePLangStmt struct {
	Tag NodeTag[CreatePLangStmt] `json:"tag"`

	Replace     bool    `json:"replace"`
	Plname      *string `json:"plname,omitempty"`
	Plhandler   *List   `json:"plhandler,omitempty"`
	Plinline    *List   `json:"plinline,omitempty"`
	Plvalidator *List   `json:"plvalidator,omitempty"`
	Pltrusted   bool    `json:"pltrusted"`
}

func (n *CreatePLangStmt) Pos() int {
	return 0
}

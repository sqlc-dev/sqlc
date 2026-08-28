package ast

type CreateConversionStmt struct {
	Tag NodeTag[CreateConversionStmt] `json:"tag"`

	ConversionName  *List   `json:",omitempty"`
	ForEncodingName *string `json:",omitempty"`
	ToEncodingName  *string `json:",omitempty"`
	FuncName        *List   `json:",omitempty"`
	Def             bool
}

func (n *CreateConversionStmt) Pos() int {
	return 0
}

package ast

type CreateConversionStmt struct {
	Tag NodeTag[CreateConversionStmt] `json:"tag"`

	ConversionName  *List
	ForEncodingName *string
	ToEncodingName  *string
	FuncName        *List
	Def             bool
}

func (n *CreateConversionStmt) Pos() int {
	return 0
}

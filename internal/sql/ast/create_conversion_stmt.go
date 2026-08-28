package ast

type CreateConversionStmt struct {
	Tag NodeTag[CreateConversionStmt] `json:"tag"`

	ConversionName  *List   `json:"conversion_name,omitempty"`
	ForEncodingName *string `json:"for_encoding_name,omitempty"`
	ToEncodingName  *string `json:"to_encoding_name,omitempty"`
	FuncName        *List   `json:"func_name,omitempty"`
	Def             bool    `json:"def"`
}

func (n *CreateConversionStmt) Pos() int {
	return 0
}

package ast

type InlineCodeBlock struct {
	Tag NodeTag[InlineCodeBlock] `json:"tag"`

	SourceText    *string
	LangOid       Oid
	LangIsTrusted bool
}

func (n *InlineCodeBlock) Pos() int {
	return 0
}

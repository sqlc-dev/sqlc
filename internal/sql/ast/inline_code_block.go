package ast

type InlineCodeBlock struct {
	SourceText    *string
	LangOid       Oid
	LangIsTrusted bool
}

func (n *InlineCodeBlock) Pos() int {
	return 0
}

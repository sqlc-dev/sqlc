package ast

type InlineCodeBlock struct {
	Tag NodeTag[InlineCodeBlock] `json:"tag"`

	SourceText    *string `json:"source_text,omitempty"`
	LangOid       Oid     `json:"lang_oid"`
	LangIsTrusted bool    `json:"lang_is_trusted"`
}

func (n *InlineCodeBlock) Pos() int {
	return 0
}

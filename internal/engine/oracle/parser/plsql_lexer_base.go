package parser

import (
	"github.com/antlr4-go/antlr/v4"
)

// PlSqlLexerBase state
type PlSqlLexerBase struct {
	*antlr.BaseLexer
}

func (l *PlSqlLexerBase) IsNewlineAtPos(pos int) bool {
	la := l.GetInputStream().LA(pos)
	return la == -1 || la == '\n'
}

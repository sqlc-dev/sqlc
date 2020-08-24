package ast

import (
	"github.com/kyleconroy/sqlc/internal/sql/ast"
)

type Constraint struct {
	Contype        ConstrType
	Conname        *string
	Deferrable     bool
	Initdeferred   bool
	Location       int
	IsNoInherit    bool
	RawExpr        ast.Node
	CookedExpr     *string
	GeneratedWhen  byte
	Keys           *List
	Exclusions     *List
	Options        *List
	Indexname      *string
	Indexspace     *string
	AccessMethod   *string
	WhereClause    ast.Node
	Pktable        *RangeVar
	FkAttrs        *List
	PkAttrs        *List
	FkMatchtype    byte
	FkUpdAction    byte
	FkDelAction    byte
	OldConpfeqop   *List
	OldPktableOid  Oid
	SkipValidation bool
	InitiallyValid bool
}

func (n *Constraint) Pos() int {
	return n.Location
}

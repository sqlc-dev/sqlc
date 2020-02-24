package ast

type Node interface {
	Pos() int
}

type Statement struct {
	Raw *RawStmt
}

type RawStmt struct {
	Stmt Node
}

func (n *RawStmt) Pos() int {
	return 0
}

type TableName struct {
	Catalog string
	Schema  string
	Name    string
}

func (n *TableName) Pos() int {
	return 0
}

type CreateTableStmt struct {
	IfNotExists bool
	Name        *TableName
	Cols        []*ColumnDef
}

func (n *CreateTableStmt) Pos() int {
	return 0
}

type DropTableStmt struct {
	IfExists bool
	Tables   []*TableName
}

func (n *DropTableStmt) Pos() int {
	return 0
}

type ColumnDef struct {
	Colname  string
	TypeName *TypeName
}

func (n *ColumnDef) Pos() int {
	return 0
}

type TypeName struct {
	Name string
}

func (n *TypeName) Pos() int {
	return 0
}

type SelectStmt struct {
	Fields *List
	From   *List
}

func (n *SelectStmt) Pos() int {
	return 0
}

type List struct {
	Items []Node
}

func (n *List) Pos() int {
	return 0
}

type ResTarget struct {
	Val Node
}

func (n *ResTarget) Pos() int {
	return 0
}

type ColumnRef struct {
	Name string
}

func (n *ColumnRef) Pos() int {
	return 0
}

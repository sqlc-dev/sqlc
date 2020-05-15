package ast

type Node interface {
	Pos() int
}

type Statement struct {
	Raw *RawStmt
}

func (n *Statement) Pos() int {
	return 0
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

type AlterTableStmt struct {
	Table *TableName
	Cmds  *List
	// MissingOk bool
}

func (n *AlterTableStmt) Pos() int {
	return 0
}

type AlterTableType int

const (
	AT_AddColumn AlterTableType = iota
	AT_AlterColumnType
	AT_DropColumn
	AT_DropNotNull
	AT_SetNotNull
)

type AlterTableCmd struct {
	Subtype   AlterTableType
	Name      *string
	Def       *ColumnDef
	MissingOk bool
}

func (n *AlterTableCmd) Pos() int {
	return 0
}

type AlterTableSetSchemaStmt struct {
	Table     *TableName
	NewSchema *string
}

func (n *AlterTableSetSchemaStmt) Pos() int {
	return 0
}

type CreateEnumStmt struct {
	TypeName *TypeName
	Vals     *List
}

func (n *CreateEnumStmt) Pos() int {
	return 0
}

type CreateSchemaStmt struct {
	Name        *string
	IfNotExists bool
}

func (n *CreateSchemaStmt) Pos() int {
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

type DropSchemaStmt struct {
	Schemas   []*String
	MissingOk bool
}

func (n *DropSchemaStmt) Pos() int {
	return 0
}

type DropTableStmt struct {
	IfExists bool
	Tables   []*TableName
}

func (n *DropTableStmt) Pos() int {
	return 0
}

type DropTypeStmt struct {
	IfExists bool
	Types    []*TypeName
}

func (n *DropTypeStmt) Pos() int {
	return 0
}

// TODO: Support array types
type ColumnDef struct {
	Colname   string
	TypeName  *TypeName
	IsNotNull bool
}

func (n *ColumnDef) Pos() int {
	return 0
}

type TypeName struct {
	Catalog string
	Schema  string
	Name    string
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

type String struct {
	Str string
}

func (n *String) Pos() int {
	return 0
}

type CommentOnSchemaStmt struct {
	Schema  *String
	Comment *string
}

func (n *CommentOnSchemaStmt) Pos() int {
	return 0
}

type CommentOnTableStmt struct {
	Table   *TableName
	Comment *string
}

func (n *CommentOnTableStmt) Pos() int {
	return 0
}

type CommentOnTypeStmt struct {
	Type    *TypeName
	Comment *string
}

func (n *CommentOnTypeStmt) Pos() int {
	return 0
}

type CommentOnColumnStmt struct {
	Table   *TableName
	Col     *ColumnRef
	Comment *string
}

func (n *CommentOnColumnStmt) Pos() int {
	return 0
}

type RenameColumnStmt struct {
	Table   *TableName
	Col     *ColumnRef
	NewName *string
}

func (n *RenameColumnStmt) Pos() int {
	return 0
}

type AlterTypeRenameValueStmt struct {
	Type     *TypeName
	OldValue *string
	NewValue *string
}

func (n *AlterTypeRenameValueStmt) Pos() int {
	return 0
}

type AlterTypeAddValueStmt struct {
	Type               *TypeName
	NewValue           *string
	SkipIfNewValExists bool
}

func (n *AlterTypeAddValueStmt) Pos() int {
	return 0
}

type RenameTableStmt struct {
	Table   *TableName
	NewName *string
}

func (n *RenameTableStmt) Pos() int {
	return 0
}

type FuncName struct {
	Catalog string
	Schema  string
	Name    string
}

func (n *FuncName) Pos() int {
	return 0
}

type CreateFunctionStmt struct {
	Replace    bool
	Params     []*FuncParam
	ReturnType *TypeName
	Func       *FuncName
}

func (n *CreateFunctionStmt) Pos() int {
	return 0
}

type DropFunctionStmt struct {
	Funcs     []*FuncSpec
	MissingOk bool
}

func (n *DropFunctionStmt) Pos() int {
	return 0
}

type FuncSpec struct {
	Name    *FuncName
	Args    []*TypeName
	HasArgs bool
}

func (n *FuncSpec) Pos() int {
	return 0
}

type FuncParam struct {
	Name    *string
	Type    *TypeName
	DefExpr Node // Will always be &ast.TODO
}

func (n *FuncParam) Pos() int {
	return 0
}

type TODO struct {
}

func (n *TODO) Pos() int {
	return 0
}

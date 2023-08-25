package compiler

import (
	"github.com/sqlc-dev/sqlc/internal/sql/ast"
)

type Function struct {
	Rel        *ast.FuncName
	ReturnType *ast.TypeName
}

type Table struct {
	Rel     *ast.TableName
	Columns []*Column
}

type Column struct {
	Name         string
	OriginalName string
	DataType     string
	NotNull      bool
	Unsigned     bool
	IsArray      bool
	ArrayDims    int
	Comment      string
	Length       *int
	IsNamedParam bool
	IsFuncCall   bool

	// XXX: Figure out what PostgreSQL calls `foo.id`
	Scope      string
	Table      *ast.TableName
	TableAlias string
	Type       *ast.TypeName
	EmbedTable *ast.TableName

	IsSqlcSlice bool // is this sqlc.slice()

	skipTableRequiredCheck bool
}

type Query struct {
	SQL      string
	Name     string
	Cmd      string // TODO: Pick a better name. One of: one, many, exec, execrows, copyFrom
	Flags    map[string]bool
	Columns  []*Column
	Params   []Parameter
	Comments []string

	// XXX: Hack
	Filename string

	// Needed for CopyFrom
	InsertIntoTable *ast.TableName

	// Needed for vet
	RawStmt *ast.RawStmt
}

type Parameter struct {
	Number int
	Column *Column
}

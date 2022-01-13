package compiler

import (
	"github.com/kyleconroy/sqlc/internal/sql/ast"
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
	DataType     string
	NotNull      bool
	IsArray      bool
	Comment      string
	Length       *int
	IsNamedParam bool
	IsFuncCall   bool

	// XXX: Figure out what PostgreSQL calls `foo.id`
	Scope      string
	Table      *ast.TableName
	TableAlias string
	Type       *ast.TypeName

	skipTableRequiredCheck bool
}

type Query struct {
	SQL      string
	Name     string
	Cmd      string // TODO: Pick a better name. One of: one, many, exec, execrows, copyFrom
	Columns  []*Column
	Params   []Parameter
	Comments []string

	// XXX: Hack
	Filename string

	// Needed for CopyFrom
	InsertIntoTable *ast.TableName
}

type Parameter struct {
	Number int
	Column *Column
}

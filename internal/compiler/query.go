package compiler

import "github.com/kyleconroy/sqlc/internal/sql/ast"

type Table struct {
	Rel     *ast.TableName
	Columns []*Column
}

type Column struct {
	Name     string
	DataType string
	NotNull  bool
	IsArray  bool
	Comment  string

	// XXX: Figure out what PostgreSQL calls `foo.id`
	Scope string
	Table *ast.TableName
	Type  *ast.TypeName
}

type Query struct {
	SQL      string
	Name     string
	Cmd      string // TODO: Pick a better name. One of: one, many, exec, execrows
	Columns  []*Column
	Params   []Parameter
	Comments []string

	// XXX: Hack
	Filename string
}

type Parameter struct {
	Number int
	Column *Column
}

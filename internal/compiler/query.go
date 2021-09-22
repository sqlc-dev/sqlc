package compiler

import (
	"fmt"

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

	// XXX: Figure out what PostgreSQL calls `foo.id`
	Scope      string
	Table      *ast.TableName
	TableAlias string
	Type       *ast.TypeName

	IsSlice bool // is this sqlc.slice

	skipTableRequiredCheck bool
}

// Named with "...Magic" because of the fixed string to be replaced
func (c *Column) InterpolatedMagic() string {
	return fmt.Sprintf(`"/*REPLACE:%s*/?"`, c.Name)
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

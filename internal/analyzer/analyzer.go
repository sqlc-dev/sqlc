package analyzer

import (
	"context"

	"github.com/sqlc-dev/sqlc/internal/sql/ast"
	"github.com/sqlc-dev/sqlc/internal/sql/named"
)

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
}

type Parameter struct {
	Number int
	Column *Column
}

type Analysis struct {
	Columns []Column
	Params  []Parameter
}

type Analyzer interface {
	Analyze(context.Context, ast.Node, string, []string, *named.ParamSet) (*Analysis, error)
	Close(context.Context) error
}

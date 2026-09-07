package compiler

import (
	"github.com/sqlc-dev/sqlc/internal/core"
	"github.com/sqlc-dev/sqlc/internal/metadata"
	"github.com/sqlc-dev/sqlc/internal/sql/ast"
	"github.com/sqlc-dev/sqlc/internal/sql/catalog"
)

type Function struct {
	Rel        *ast.FuncName
	ReturnType *ast.TypeName
	Outs       []*catalog.Argument
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

	// TypeExpr is the type as the analysis core wrote it, with the
	// arguments and nesting DataType and IsArray flatten away. It is unset
	// on the legacy path.
	TypeExpr *core.TypeExpr

	IsSqlcSlice bool // is this sqlc.slice()

	skipTableRequiredCheck bool
}

type Query struct {
	SQL      string
	Metadata metadata.Metadata
	Columns  []*Column
	Params   []Parameter

	// Needed for CopyFrom
	InsertIntoTable *ast.TableName

	// Needed for vet
	RawStmt *ast.RawStmt
}

type Parameter struct {
	Number int
	Column *Column
}

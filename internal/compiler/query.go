package compiler

import (
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

	// Target table for UPDATE queries
	UpdateTable *ast.TableName

	// Target table for DELETE queries
	DeleteFromTable *ast.TableName

	// Needed for vet
	RawStmt *ast.RawStmt
}

type ParameterContext int

const (
	ParameterContextUnspecified ParameterContext = iota
	ParameterContextSet
	ParameterContextValues
	ParameterContextWhere
	ParameterContextHaving
	ParameterContextFunctionArg
	ParameterContextLimit
	ParameterContextOffset
)

type Parameter struct {
	Number  int
	Column  *Column
	Context ParameterContext
}

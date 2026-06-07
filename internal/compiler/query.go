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

	// SourceTables lists the base tables the query reads from, including tables
	// that appear only in joins, subqueries, or common table expression bodies.
	// Names are schema-qualified when a schema is present, deduplicated, and
	// sorted. Common table expression names and the target relations of write
	// statements are excluded.
	SourceTables []string

	// Needed for CopyFrom
	InsertIntoTable *ast.TableName

	// Needed for vet
	RawStmt *ast.RawStmt
}

type Parameter struct {
	Number int
	Column *Column
}

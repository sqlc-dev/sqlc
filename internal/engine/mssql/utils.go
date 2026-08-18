package mssql

import (
	"log"
	"strings"

	tsql "github.com/sqlc-dev/teesql/ast"

	"github.com/sqlc-dev/sqlc/internal/debug"
	"github.com/sqlc-dev/sqlc/internal/sql/ast"
)

// fragmented is satisfied by every teesql node via its embedded Fragment.
type fragmented interface {
	Frag() *tsql.Fragment
}

// location returns a node's start offset, for the Location fields of the sqlc
// AST. Spans are recorded in UTF-16 code units, which agree with byte offsets
// for ASCII input.
func location(n any) int {
	if f, ok := n.(fragmented); ok {
		return f.Frag().StartOffset
	}
	return 0
}

func todo(n tsql.Node) *ast.TODO {
	if debug.Active {
		log.Printf("mssql.convert: Unknown node type %T\n", n)
	}
	return &ast.TODO{}
}

// identifier normalizes an identifier. T-SQL identifiers are
// case-insensitive under the default collations.
func identifier(id string) string {
	return strings.ToLower(id)
}

func NewIdentifier(t string) *ast.String {
	return &ast.String{Str: identifier(t)}
}

func identifierValue(id *tsql.Identifier) string {
	if id == nil {
		return ""
	}
	return identifier(id.Value)
}

// schemaName returns the schema an object name qualifies it with. The "dbo"
// default schema maps to the catalog's default namespace, so it is treated
// as unqualified.
func schemaName(n *tsql.SchemaObjectName) string {
	if n == nil || n.SchemaIdentifier == nil {
		return ""
	}
	s := identifier(n.SchemaIdentifier.Value)
	if s == "dbo" {
		return ""
	}
	return s
}

func parseTableName(n *tsql.SchemaObjectName) *ast.TableName {
	if n == nil {
		return &ast.TableName{}
	}
	return &ast.TableName{
		Schema: schemaName(n),
		Name:   identifierValue(n.BaseIdentifier),
	}
}

func parseRangeVar(n *tsql.SchemaObjectName) *ast.RangeVar {
	name := parseTableName(n)
	rv := &ast.RangeVar{
		Relname:  &name.Name,
		Location: location(n),
	}
	if name.Schema != "" {
		rv.Schemaname = &name.Schema
	}
	return rv
}

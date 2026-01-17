package clickhouse

import (
	"log"
	"strings"

	chast "github.com/sqlc-dev/doubleclick/ast"

	"github.com/sqlc-dev/sqlc/internal/debug"
	"github.com/sqlc-dev/sqlc/internal/sql/ast"
)

func todo(n chast.Node) *ast.TODO {
	if debug.Active {
		log.Printf("clickhouse.convert: Unknown node type %T\n", n)
	}
	return &ast.TODO{}
}

func identifier(id string) string {
	return strings.ToLower(id)
}

func NewIdentifier(t string) *ast.String {
	return &ast.String{Str: identifier(t)}
}

func parseTableName(n *chast.TableIdentifier) *ast.TableName {
	return &ast.TableName{
		Schema: identifier(n.Database),
		Name:   identifier(n.Table),
	}
}

func parseTableIdentifierToRangeVar(n *chast.TableIdentifier) *ast.RangeVar {
	schemaname := identifier(n.Database)
	relname := identifier(n.Table)
	return &ast.RangeVar{
		Schemaname: &schemaname,
		Relname:    &relname,
	}
}

func isNotNull(n *chast.ColumnDeclaration) bool {
	if n.Type == nil {
		return false
	}
	// Check if type is wrapped in Nullable()
	// If it's Nullable, it can be null, so return false
	// If it's not Nullable, it's NOT NULL by default in ClickHouse
	if n.Type.Name != "" && strings.ToLower(n.Type.Name) == "nullable" {
		return false
	}
	// Also check if Nullable field is explicitly set
	if n.Nullable != nil && *n.Nullable {
		return false
	}
	return true
}

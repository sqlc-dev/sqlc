package sqlite

import (
	"strings"

	meyer "github.com/sqlc-dev/meyer/ast"

	"github.com/sqlc-dev/sqlc/internal/sql/ast"
)

// identifier folds a name the way SQLite resolves one: an unquoted name is
// case-insensitive and normalizes to lower case, while a quoted name keeps
// the spelling that appeared between the quotes.
func identifier(n *meyer.Ident) string {
	if n == nil {
		return ""
	}
	if n.Quote == meyer.QuoteNone {
		return strings.ToLower(n.Name)
	}
	return n.Name
}

// NewIdentifier returns the folded name of id as an ast.String, the shape
// sqlc uses for the parts of a ColumnRef or a function name.
func NewIdentifier(id *meyer.Ident) *ast.String {
	return &ast.String{Str: identifier(id)}
}

func parseTableName(n *meyer.QualifiedName) *ast.TableName {
	if n == nil {
		return &ast.TableName{}
	}
	return &ast.TableName{
		Schema: identifier(n.Schema),
		Name:   identifier(n.Name),
	}
}

// parseRangeVar builds the relation reference sqlc expects for a table named
// in a FROM, INSERT, UPDATE or DELETE clause.
func parseRangeVar(n *meyer.QualifiedName, alias *meyer.Ident) *ast.RangeVar {
	rv := &ast.RangeVar{}
	if n == nil {
		return rv
	}
	rv.Location = n.Pos()
	name := identifier(n.Name)
	rv.Relname = &name
	if n.Schema != nil {
		schema := identifier(n.Schema)
		rv.Schemaname = &schema
	}
	if alias != nil {
		aliasName := identifier(alias)
		rv.Alias = &ast.Alias{Aliasname: &aliasName}
	}
	return rv
}

// typeName renders a type the way sqlc's type resolution expects to see it:
// the tokens joined without separators, matching what the ANTLR parser's
// GetText produced. "UNSIGNED BIG INT" becomes "UNSIGNEDBIGINT" and
// "DECIMAL(10, 5)" becomes "DECIMAL(10,5)".
func typeName(n *meyer.TypeName) string {
	if n == nil {
		return ""
	}
	var sb strings.Builder
	sb.WriteString(strings.ReplaceAll(n.Name, " ", ""))
	if len(n.Args) > 0 {
		sb.WriteString("(")
		sb.WriteString(strings.Join(n.Args, ","))
		sb.WriteString(")")
	}
	return sb.String()
}

// hasNotNullConstraint reports whether a column definition guarantees a
// value. A PRIMARY KEY implies NOT NULL for the purposes of code
// generation, matching the behavior of the previous parser.
func hasNotNullConstraint(constraints []*meyer.ColumnConstraint) bool {
	for _, constraint := range constraints {
		switch constraint.Kind {
		case meyer.ColumnNotNull, meyer.ColumnPrimaryKey:
			return true
		}
	}
	return false
}

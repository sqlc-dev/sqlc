package googlesql

import (
	"log"
	"strings"

	zjast "github.com/sqlc-dev/zetajones/ast"

	"github.com/sqlc-dev/sqlc/internal/debug"
	"github.com/sqlc-dev/sqlc/internal/sql/ast"
)

func todo(n zjast.Node) *ast.TODO {
	if debug.Active {
		log.Printf("googlesql.convert: Unknown node type %T\n", n)
	}
	return &ast.TODO{}
}

// identifier passes table and column identifiers through unchanged. GoogleSQL
// keywords are case-insensitive, but table names (on BigQuery) are
// case-sensitive, so unlike the MySQL and ClickHouse engines we do not
// lowercase identifiers here. Function and type names, which are
// case-insensitive, are lowercased at their call sites instead.
func identifier(id string) string {
	return id
}

func NewIdentifier(t string) *ast.String {
	return &ast.String{Str: identifier(t)}
}

// pathParts returns the identifier names of a dotted path expression.
func pathParts(n *zjast.PathExpression) []string {
	if n == nil {
		return nil
	}
	parts := make([]string, 0, len(n.Names))
	for _, id := range n.Names {
		if id == nil {
			continue
		}
		parts = append(parts, id.Name)
	}
	return parts
}

func parseTableName(n *zjast.PathExpression) *ast.TableName {
	name := &ast.TableName{}
	parts := pathParts(n)
	switch len(parts) {
	case 0:
	case 1:
		name.Name = identifier(parts[0])
	case 2:
		name.Schema = identifier(parts[0])
		name.Name = identifier(parts[1])
	default:
		name.Catalog = identifier(parts[0])
		name.Schema = identifier(parts[len(parts)-2])
		name.Name = identifier(parts[len(parts)-1])
	}
	return name
}

func parseRangeVar(n *zjast.PathExpression) *ast.RangeVar {
	tn := parseTableName(n)
	rv := &ast.RangeVar{}
	if n != nil {
		rv.Location = n.Pos()
	}
	if tn.Catalog != "" {
		catalog := tn.Catalog
		rv.Catalogname = &catalog
	}
	if tn.Schema != "" {
		schema := tn.Schema
		rv.Schemaname = &schema
	}
	if tn.Name != "" {
		relname := tn.Name
		rv.Relname = &relname
	}
	return rv
}

func convertAlias(a *zjast.Alias) *ast.Alias {
	if a == nil || a.Identifier == nil {
		return nil
	}
	name := identifier(a.Identifier.Name)
	return &ast.Alias{Aliasname: &name}
}

// whereExpr unwraps a WHERE clause node to its expression. The DELETE and
// UPDATE statement nodes hold either a *WhereClause or a bare expression.
func whereExpr(node zjast.Node) zjast.Node {
	if wc, ok := node.(*zjast.WhereClause); ok {
		return wc.Expr
	}
	return node
}

// stringLiteralValue concatenates and unquotes the components of a string
// literal.
func stringLiteralValue(n *zjast.StringLiteral) string {
	var sb strings.Builder
	for _, comp := range n.Components {
		if comp == nil {
			continue
		}
		sb.WriteString(unquote(comp.Image))
	}
	return sb.String()
}

// unquote strips the optional r/b prefix and the surrounding quotes from a
// string literal component and decodes escape sequences for non-raw strings.
// It is best-effort: malformed input is returned unchanged.
func unquote(image string) string {
	raw := false
	i := 0
prefix:
	for i < len(image) {
		switch image[i] {
		case 'r', 'R':
			raw = true
			i++
		case 'b', 'B':
			i++
		default:
			break prefix
		}
	}
	s := image[i:]
	for _, q := range []string{`'''`, `"""`, `'`, `"`} {
		if len(s) >= 2*len(q) && strings.HasPrefix(s, q) && strings.HasSuffix(s, q) {
			s = s[len(q) : len(s)-len(q)]
			if raw || !strings.ContainsRune(s, '\\') {
				return s
			}
			return decodeEscapes(s)
		}
	}
	return image
}

func decodeEscapes(s string) string {
	var sb strings.Builder
	for i := 0; i < len(s); i++ {
		ch := s[i]
		if ch != '\\' || i+1 >= len(s) {
			sb.WriteByte(ch)
			continue
		}
		i++
		switch s[i] {
		case 'n':
			sb.WriteByte('\n')
		case 't':
			sb.WriteByte('\t')
		case 'r':
			sb.WriteByte('\r')
		case 'a':
			sb.WriteByte('\a')
		case 'b':
			sb.WriteByte('\b')
		case 'f':
			sb.WriteByte('\f')
		case 'v':
			sb.WriteByte('\v')
		case '\\', '\'', '"', '?', '`':
			sb.WriteByte(s[i])
		default:
			// Unknown escape (\x, \u, ...): keep it verbatim.
			sb.WriteByte('\\')
			sb.WriteByte(s[i])
		}
	}
	return sb.String()
}

// typeName renders a zetajones type node (used by CAST) as a lowercased type
// name. Nested types degrade to a readable representation.
func typeName(node zjast.Node) string {
	switch t := node.(type) {
	case *zjast.SimpleType:
		return strings.ToLower(strings.Join(pathParts(t.Name), "."))
	case *zjast.ArrayType:
		return "array<" + typeName(t.ElementType) + ">"
	case *zjast.StructType:
		return "struct"
	case *zjast.RangeType:
		return "range<" + typeName(t.ElementType) + ">"
	case *zjast.MapType:
		return "map"
	default:
		return ""
	}
}

// columnSchemaTypeName renders a CREATE TABLE column schema as a lowercased
// type name. Nested types degrade to a readable representation.
func columnSchemaTypeName(node zjast.Node) string {
	switch t := node.(type) {
	case *zjast.SimpleColumnSchema:
		return strings.ToLower(strings.Join(pathParts(t.Type), "."))
	case *zjast.ArrayColumnSchema:
		return "array<" + columnSchemaTypeName(t.ElementSchema) + ">"
	case *zjast.StructColumnSchema:
		return "struct"
	default:
		return ""
	}
}

// columnAttributes returns the attribute list (NOT NULL, PRIMARY KEY, ...) of
// a column schema, if any.
func columnAttributes(node zjast.Node) *zjast.ColumnAttributeList {
	switch t := node.(type) {
	case *zjast.SimpleColumnSchema:
		return t.Attributes
	case *zjast.ArrayColumnSchema:
		return t.Attributes
	case *zjast.StructColumnSchema:
		return t.Attributes
	default:
		return nil
	}
}

// columnRefFields flattens an identifier, path expression, or dotted
// identifier chain into ColumnRef fields.
func columnRefFields(node zjast.Node) ([]ast.Node, bool) {
	switch n := node.(type) {
	case *zjast.Identifier:
		return []ast.Node{NewIdentifier(n.Name)}, true
	case *zjast.PathExpression:
		var fields []ast.Node
		for _, part := range pathParts(n) {
			fields = append(fields, NewIdentifier(part))
		}
		return fields, len(fields) > 0
	case *zjast.DotIdentifier:
		base, ok := columnRefFields(n.Expr)
		if !ok || n.Name == nil {
			return nil, false
		}
		return append(base, NewIdentifier(n.Name.Name)), true
	default:
		return nil, false
	}
}

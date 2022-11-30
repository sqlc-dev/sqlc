package rewrite

import (
	"fmt"
	"strings"

	"github.com/kyleconroy/sqlc/internal/sql/ast"
	"github.com/kyleconroy/sqlc/internal/sql/astutils"
)

// Embed is an instance of `sqlc.embed(param)`
type Embed struct {
	Table *ast.TableName
	param string
	Node  *ast.ColumnRef
}

// Orig string to replace
func (e Embed) Orig() string {
	return fmt.Sprintf("sqlc.embed(%s)", e.param)
}

// EmbedSet is a set of Embed instances
type EmbedSet []*Embed

// Find a matching embed by column ref
func (es EmbedSet) Find(node *ast.ColumnRef) (*Embed, bool) {
	for _, e := range es {
		if e.Node == node {
			return e, true
		}
	}
	return nil, false
}

// Embeds rewrites `sqlc.embed(param)` to a `ast.ColumnRef` of form `param.*`.
// The compiler can make use of the returned `EmbedSet` while expanding the
// `param.*` column refs to produce the correct source edits.
func Embeds(raw *ast.RawStmt) (*ast.RawStmt, EmbedSet) {
	var embeds []*Embed

	node := astutils.Apply(raw, func(cr *astutils.Cursor) bool {
		node := cr.Node()

		switch {
		case isEmbed(node):
			fun := node.(*ast.FuncCall)

			if len(fun.Args.Items) == 0 {
				return false
			}

			sw := &stringWalker{}
			astutils.Walk(sw, fun.Args)
			str := strings.Join(sw.Parts, ".")

			tableName, err := parseTable(sw)
			if err != nil {
				return false
			}

			node := &ast.ColumnRef{
				Fields: &ast.List{
					Items: []ast.Node{},
				},
			}

			for _, s := range sw.Parts {
				node.Fields.Items = append(node.Fields.Items, &ast.String{Str: s})
			}
			node.Fields.Items = append(node.Fields.Items, &ast.A_Star{})

			embeds = append(embeds, &Embed{
				Table: tableName,
				param: str,
				Node:  node,
			})

			cr.Replace(node)
			return false
		default:
			return true
		}
	}, nil)

	return node.(*ast.RawStmt), embeds
}

func isEmbed(node ast.Node) bool {
	call, ok := node.(*ast.FuncCall)
	if !ok {
		return false
	}

	if call.Func == nil {
		return false
	}

	isValid := call.Func.Schema == "sqlc" && call.Func.Name == "embed"
	return isValid
}

func parseTable(sw *stringWalker) (*ast.TableName, error) {
	parts := sw.Parts

	switch len(parts) {
	case 1:
		return &ast.TableName{
			Name: parts[0],
		}, nil
	case 2:
		return &ast.TableName{
			Schema: parts[0],
			Name:   parts[1],
		}, nil
	case 3:
		return &ast.TableName{
			Catalog: parts[0],
			Schema:  parts[1],
			Name:    parts[2],
		}, nil
	default:
		return nil, fmt.Errorf("invalid table name")
	}
}

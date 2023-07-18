package rewrite

import (
	"fmt"

	"github.com/kyleconroy/sqlc/internal/sql/ast"
	"github.com/kyleconroy/sqlc/internal/sql/astutils"
)

// Embed is an instance of `sqlc.embed(param)` or `sqlc.nembed(param)`.
// The only difference in an embed generated with `nembed` is that `Nullable`
// will always be `true`.
type Embed struct {
	Table    *ast.TableName
	param    string
	Node     *ast.ColumnRef
	Nullable bool
}

// Orig string to replace
func (e Embed) Orig() string {
	fName := "embed"
	if e.Nullable {
		fName = "nembed"
	}
	return fmt.Sprintf("sqlc.%s(%s)", fName, e.param)
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

// Embeds rewrites `sqlc.embed(param)` or `sqlc.nembed(param)` to an `ast.ColumnRef`
// of form `param.*`. The compiler can make use of the returned `EmbedSet` while
// expanding the `param.*` column refs to produce the correct source edits.
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

			param, _ := flatten(fun.Args)

			node := &ast.ColumnRef{
				Fields: &ast.List{
					Items: []ast.Node{
						&ast.String{Str: param},
						&ast.A_Star{},
					},
				},
			}

			nullable := false
			if fun.Func.Name == "nembed" {
				nullable = true
			}
			embeds = append(embeds, &Embed{
				Table:    &ast.TableName{Name: param},
				param:    param,
				Node:     node,
				Nullable: nullable,
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

	isValid := call.Func.Schema == "sqlc" && (call.Func.Name == "embed" || call.Func.Name == "nembed")
	return isValid
}

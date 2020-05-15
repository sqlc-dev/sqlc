package astutils

import (
	"strings"

	"github.com/kyleconroy/sqlc/internal/sql/ast"
	"github.com/kyleconroy/sqlc/internal/sql/ast/pg"
)

func Join(list *ast.List, sep string) string {
	items := []string{}
	for _, item := range list.Items {
		if n, ok := item.(*ast.String); ok {
			items = append(items, n.Str)
		}
		if n, ok := item.(*pg.String); ok {
			items = append(items, n.Str)
		}
	}
	return strings.Join(items, sep)
}

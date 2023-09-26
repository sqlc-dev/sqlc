package astutils

import (
	"strings"

	"github.com/Dionid/sqlc/internal/sql/ast"
)

func Join(list *ast.List, sep string) string {
	if list == nil {
		return ""
	}

	var items []string
	for _, item := range list.Items {
		if n, ok := item.(*ast.String); ok {
			items = append(items, n.Str)
		}
	}
	return strings.Join(items, sep)
}

package ast

import (
	"strings"

	nodes "github.com/lfittl/pg_query_go/nodes"
)

func Join(list nodes.List, sep string) string {
	items := []string{}
	for _, item := range list.Items {
		if n, ok := item.(nodes.String); ok {
			items = append(items, n.Str)
		}
	}
	return strings.Join(items, sep)
}

package dinosql

import (
	"fmt"

	nodes "github.com/lfittl/pg_query_go/nodes"
)

func validateParamRef(n nodes.Node) error {
	var allrefs []nodes.ParamRef

	// Find all parameter references
	Walk(VisitorFunc(func(node nodes.Node) {
		switch n := node.(type) {
		case nodes.ParamRef:
			allrefs = append(allrefs, n)
		}
	}), n)

	seen := map[int]struct{}{}
	for _, r := range allrefs {
		seen[r.Number] = struct{}{}
	}

	for i := 1; i <= len(seen); i += 1 {
		if _, ok := seen[i]; !ok {
			return fmt.Errorf("missing parameter reference: $%d", i)
		}
	}

	return nil
}

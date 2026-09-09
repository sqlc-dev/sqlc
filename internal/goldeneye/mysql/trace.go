package mysql

import (
	"encoding/json"
	"fmt"
	"sort"
	"strings"
)

// trace is what the optimizer trace says about a statement: the text of
// each query block after name resolution and before optimisation, by the
// number MySQL gives the block, and the block each derived table or CTE is
// materialised from.
//
// The trace is the one place MySQL prints a query block as it was resolved
// rather than as it was optimised. The note EXPLAIN leaves is printed
// after optimisation, and a lookup on a unique key against an empty table
// folds to `NULL = (@x)` there, losing the column the parameter was
// compared with; the expanded_query of the block still has it.
type trace struct {
	blocks  map[int]text
	derived map[string]int
}

// parseTrace reads the JSON of the OPTIMIZER_TRACE table's TRACE column.
func parseTrace(blob string) (trace, error) {
	tr := trace{blocks: map[int]text{}, derived: map[string]int{}}
	if strings.TrimSpace(blob) == "" {
		return tr, nil
	}
	var v any
	if err := json.Unmarshal([]byte(blob), &v); err != nil {
		return tr, fmt.Errorf("reading the optimizer trace: %w", err)
	}
	tr.walk(v, 0)
	return tr, nil
}

// walk collects the first expanded_query printed for each block — the one
// join_preparation prints, before join_optimization prints its rewrite —
// and the derived tables. Lists are walked in order and objects in key
// order, so the result does not depend on map iteration.
func (tr *trace) walk(v any, sel int) {
	switch v := v.(type) {
	case map[string]any:
		if n, ok := v["select#"].(float64); ok {
			sel = int(n)
		}
		if d, ok := v["derived"].(map[string]any); ok {
			if table, ok := d["table"].(string); ok {
				if n, ok := d["select#"].(float64); ok {
					tr.derived[strings.Trim(table, "` ")] = int(n)
				}
			}
		}
		if q, ok := v["expanded_query"].(string); ok {
			if _, seen := tr.blocks[sel]; !seen {
				tr.blocks[sel] = tokenize(q)
			}
		}
		keys := make([]string, 0, len(v))
		for k := range v {
			keys = append(keys, k)
		}
		sort.Strings(keys)
		for _, k := range keys {
			tr.walk(v[k], sel)
		}
	case []any:
		for _, c := range v {
			tr.walk(c, sel)
		}
	}
}

// order lists the block numbers innermost first: a block's text includes
// its subqueries' text, and a subquery is numbered after the block it is
// in, so the block a parameter belongs to is the highest-numbered block
// whose text has it.
func (tr *trace) order() []int {
	sels := make([]int, 0, len(tr.blocks))
	for sel := range tr.blocks {
		sels = append(sels, sel)
	}
	sort.Sort(sort.Reverse(sort.IntSlice(sels)))
	return sels
}

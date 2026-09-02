package main

import (
	"context"
	"encoding/json"
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

// The output follows the JSON `sqlc analyze` prints, with each type written
// as a call expression (see types.go) instead of a name and flags.

type analyzedQuery struct {
	Name    string           `json:"name"`
	Cmd     string           `json:"cmd"`
	Columns []analyzedColumn `json:"columns"`
	Params  []analyzedParam  `json:"params"`
}

type analyzedColumn struct {
	Name    string `json:"name"`
	Type    string `json:"type"`
	NotNull bool   `json:"not_null"`
	Table   string `json:"table,omitempty"`
}

type analyzedParam struct {
	Number int            `json:"number"`
	Column analyzedColumn `json:"column"`
}

// analyze runs every query against the schema and fixture and records what
// ClickHouse reports about each.
func analyze(ctx context.Context, l local, schema, fixture string, queries []query) ([]analyzedQuery, error) {
	out := make([]analyzedQuery, 0, len(queries))
	for _, q := range queries {
		aq, err := analyzeQuery(ctx, l, schema, fixture, q)
		if err != nil {
			return nil, fmt.Errorf("%s: %w", q.Name, err)
		}
		out = append(out, aq)
	}
	return out, nil
}

func analyzeQuery(ctx context.Context, l local, schema, fixture string, q query) (analyzedQuery, error) {
	sql, phs := bindPlaceholders(q.SQL)
	explain := returnsRows(sql)

	var script strings.Builder
	for _, stmt := range []string{schema, fixture} {
		if s := strings.TrimRight(strings.TrimSpace(stmt), ";"); s != "" {
			script.WriteString(s)
			script.WriteString(";\n")
		}
	}
	if explain {
		script.WriteString("EXPLAIN QUERY TREE " + sql + ";\n")
	}
	script.WriteString(sql + ";\n")

	results, err := l.run(ctx, script.String())
	if err != nil {
		return analyzedQuery{}, err
	}

	aq := analyzedQuery{
		Name:    q.Name,
		Cmd:     q.Cmd,
		Columns: []analyzedColumn{},
		Params:  []analyzedParam{},
	}
	if !explain {
		return analyzeExec(ctx, l, script.String(), sql, phs, aq)
	}
	if len(results) != 2 {
		return analyzedQuery{}, fmt.Errorf("expected the query tree and one result set, got %d results", len(results))
	}

	var lines []string
	for _, row := range results[0].Data {
		var line string
		if err := json.Unmarshal(row["explain"], &line); err != nil {
			return analyzedQuery{}, fmt.Errorf("reading query tree: %w", err)
		}
		lines = append(lines, line)
	}
	tree, err := parseQueryTree(lines)
	if err != nil {
		return analyzedQuery{}, err
	}

	// Names and types come from the block header of the executed query, the
	// same header a driver sees. The tree only adds where each came from.
	_, nodes := projection(firstNode(tree.root.children))
	for i, col := range results[1].Meta {
		ac := column(col.Name, col.Type)
		if i < len(nodes) && len(nodes) == len(results[1].Meta) {
			ac.Table = tree.sourceTable(nodes[i])
		}
		aq.Columns = append(aq.Columns, ac)
	}

	sentinels := tree.sentinels()
	for i, ph := range phs {
		ac := analyzedColumn{}
		if sentinel := sentinels[i+1]; sentinel != nil {
			ac = tree.paramColumn(sentinel)
		}
		if ph.Name != "" {
			ac.Name = ph.Name
		}
		aq.Params = append(aq.Params, analyzedParam{Number: ph.Number, Column: ac})
	}
	return aq, nil
}

func column(name, typ string) analyzedColumn {
	expr, notNull := typeExpr(typ)
	return analyzedColumn{Name: name, Type: expr, NotNull: notNull}
}

// returnsRows reports whether a statement produces a result set and so can
// be explained as a query tree.
func returnsRows(sql string) bool {
	head := strings.ToLower(strings.TrimSpace(sql))
	if strings.HasPrefix(head, "(") {
		return true
	}
	for _, kw := range []string{"select", "with", "show", "describe", "desc", "exists"} {
		if strings.HasPrefix(head, kw) && (len(head) == len(kw) || !isWordByte(head[len(kw)])) {
			return true
		}
	}
	return false
}

// sentinels finds the constants the placeholders were substituted with,
// keyed by placeholder ordinal.
func (t *queryTree) sentinels() map[int]*treeNode {
	found := map[int]*treeNode{}
	var walk func(n *treeNode)
	walk = func(n *treeNode) {
		if n.kind == "CONSTANT" {
			if k, ok := sentinelOrdinal(n); ok {
				found[k] = n
				return
			}
		}
		for _, c := range n.children {
			walk(c)
		}
	}
	walk(t.root)
	return found
}

// sentinelOrdinal recognises a constant folded directly from `NULL + k` or
// `4294967295 + k`, the shapes sentinelFor produces, and returns k. A
// constant folded from a larger expression that merely contains a sentinel
// does not match; the sentinel is found nested inside it instead.
func sentinelOrdinal(c *treeNode) (int, bool) {
	fn := firstNode(c.section("EXPRESSION").childrenOrNil())
	if fn == nil || fn.kind != "FUNCTION" || fn.attrs["function_name"] != "plus" {
		return 0, false
	}
	args := fn.section("ARGUMENTS").list()
	if len(args) != 2 || args[0].kind != "CONSTANT" || args[1].kind != "CONSTANT" {
		return 0, false
	}
	base := args[0].attrs["constant_value"]
	if base != "NULL" && base != "UInt64_"+limitBase {
		return 0, false
	}
	k, err := strconv.Atoi(strings.TrimPrefix(args[1].attrs["constant_value"], "UInt64_"))
	if err != nil {
		return 0, false
	}
	return k, true
}

// paramColumn describes what a placeholder is compared with or assigned to:
// the other operand of the function it is an argument of, preferring a
// column over an expression, or the projected column it stands for.
func (t *queryTree) paramColumn(sentinel *treeNode) analyzedColumn {
	list := sentinel.parent
	if list != nil && list.kind == "LIST" && list.parent != nil {
		switch owner := list.parent; {
		case owner.kind == "" && owner.text == "ARGUMENTS":
			// Prefer a column operand, then an expression, then a constant.
			rank := map[string]int{"COLUMN": 0, "FUNCTION": 1, "CONSTANT": 2}
			var best *treeNode
			for _, sib := range list.children {
				if sib == sentinel {
					continue
				}
				if r, ok := rank[sib.kind]; ok && (best == nil || r < rank[best.kind]) {
					best = sib
				}
			}
			if best != nil {
				return t.describe(best)
			}
		case owner.kind == "" && owner.text == "PROJECTION":
			names, nodes := projection(owner.parent)
			for i, n := range nodes {
				if n == sentinel && i < len(names) {
					ac := column(names[i], sentinel.attrs["constant_value_type"])
					return ac
				}
			}
		}
	}
	return column("", sentinel.attrs["constant_value_type"])
}

// describe turns a tree expression into a column description.
func (t *queryTree) describe(n *treeNode) analyzedColumn {
	switch n.kind {
	case "COLUMN":
		ac := column(n.attrs["column_name"], n.attrs["result_type"])
		ac.Table = t.sourceTable(n)
		return ac
	case "FUNCTION":
		return column(n.attrs["function_name"], n.attrs["result_type"])
	case "CONSTANT":
		return column("", n.attrs["constant_value_type"])
	}
	return analyzedColumn{}
}

var insertValuesRe = regexp.MustCompile(`(?is)^insert\s+into\s+(?:table\s+)?([\w.` + "`" + `"]+)\s*(?:\(([^)]*)\))?\s*(?:format\s+)?values\b`)

// analyzeExec runs a statement that returns no rows. The only parameters it
// can describe are those of an INSERT ... VALUES, which map positionally
// onto the target columns reported by DESCRIBE TABLE.
func analyzeExec(ctx context.Context, l local, script, sql string, phs []placeholder, aq analyzedQuery) (analyzedQuery, error) {
	m := insertValuesRe.FindStringSubmatch(sql)
	if m != nil {
		script += "DESCRIBE TABLE " + m[1] + ";\n"
	}
	results, err := l.run(ctx, script)
	if err != nil {
		return analyzedQuery{}, err
	}

	var targets []analyzedColumn
	if m != nil && len(results) == 1 {
		byName := map[string]analyzedColumn{}
		var all []analyzedColumn
		table := strings.Trim(m[1][strings.LastIndexByte(m[1], '.')+1:], "`\"")
		for _, row := range results[0].Data {
			var name, typ string
			json.Unmarshal(row["name"], &name)
			json.Unmarshal(row["type"], &typ)
			ac := column(name, typ)
			ac.Table = table
			byName[name] = ac
			all = append(all, ac)
		}
		if strings.TrimSpace(m[2]) == "" {
			targets = all
		} else {
			for _, name := range strings.Split(m[2], ",") {
				targets = append(targets, byName[strings.Trim(strings.TrimSpace(name), "`\"")])
			}
		}
	}
	for i, ph := range phs {
		ac := analyzedColumn{}
		if len(targets) > 0 {
			ac = targets[i%len(targets)]
		}
		if ph.Name != "" {
			ac.Name = ph.Name
		}
		aq.Params = append(aq.Params, analyzedParam{Number: ph.Number, Column: ac})
	}
	return aq, nil
}

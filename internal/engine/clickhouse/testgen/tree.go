package main

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

// treeNode is one line of EXPLAIN QUERY TREE output. Lines of the form
// `KIND id: N, key: value, ...` are nodes with a kind and attributes; every
// other line (section headers such as PROJECTION or ARGUMENTS, and the
// `name Type` lines under PROJECTION COLUMNS) keeps only its text.
type treeNode struct {
	kind     string
	id       int
	attrs    map[string]string
	text     string
	parent   *treeNode
	children []*treeNode
}

// queryTree is a parsed EXPLAIN QUERY TREE dump.
type queryTree struct {
	root *treeNode
	byID map[int]*treeNode
}

var nodeLineRe = regexp.MustCompile(`^([A-Z_]+) id: (\d+)(?:, (.*))?$`)

// parseQueryTree builds the tree from the dump's lines, using the two-space
// indentation to recover nesting.
func parseQueryTree(lines []string) (*queryTree, error) {
	root := &treeNode{id: -1, text: "<root>"}
	t := &queryTree{root: root, byID: map[int]*treeNode{}}
	stack := []*treeNode{root}
	for _, line := range lines {
		if strings.TrimSpace(line) == "" {
			continue
		}
		indent := len(line) - len(strings.TrimLeft(line, " "))
		if indent%2 != 0 {
			return nil, fmt.Errorf("query tree: unexpected indentation in %q", line)
		}
		depth := indent/2 + 1
		if depth > len(stack) {
			return nil, fmt.Errorf("query tree: line %q is nested too deeply", line)
		}
		stack = stack[:depth]
		parent := stack[len(stack)-1]

		n := &treeNode{id: -1, text: line[indent:], parent: parent}
		if m := nodeLineRe.FindStringSubmatch(n.text); m != nil {
			n.kind = m[1]
			n.id, _ = strconv.Atoi(m[2])
			n.attrs = parseAttrs(m[3])
			t.byID[n.id] = n
		}
		parent.children = append(parent.children, n)
		stack = append(stack, n)
	}
	return t, nil
}

var attrKeyRe = regexp.MustCompile(`^, ([a-z_]+): `)

// parseAttrs splits `key: value, key: value` where values may themselves
// contain commas inside parentheses or quotes, as types and constants do.
func parseAttrs(s string) map[string]string {
	attrs := map[string]string{}
	if s == "" {
		return attrs
	}
	// Positions at which a new `, key: ` begins at top level.
	var cuts []int
	depth := 0
	inQuote := false
	for i := 0; i < len(s); i++ {
		switch c := s[i]; {
		case c == '\'':
			inQuote = !inQuote
		case inQuote:
		case c == '(':
			depth++
		case c == ')':
			depth--
		case c == ',' && depth == 0 && attrKeyRe.MatchString(s[i:]):
			cuts = append(cuts, i)
		}
	}
	cuts = append(cuts, len(s))
	start := 0
	for _, cut := range cuts {
		pair := s[start:cut]
		if key, value, ok := strings.Cut(pair, ": "); ok {
			attrs[key] = value
		}
		start = cut + 2
	}
	return attrs
}

// section returns the child section header of a node, such as PROJECTION
// or JOIN TREE, or nil.
func (n *treeNode) section(name string) *treeNode {
	for _, c := range n.children {
		if c.kind == "" && c.text == name {
			return c
		}
	}
	return nil
}

// list returns the nodes of the LIST under a section header.
func (n *treeNode) list() []*treeNode {
	if n == nil {
		return nil
	}
	for _, c := range n.children {
		if c.kind == "LIST" {
			return c.children
		}
	}
	return nil
}

// firstQuery returns the QUERY node describing a result set: the node itself
// or, for a UNION, its first branch, whose projection names the union's
// columns.
func firstQuery(n *treeNode) *treeNode {
	for depth := 0; n != nil && depth < 32; depth++ {
		switch n.kind {
		case "QUERY":
			return n
		case "UNION":
			n = firstNode(n.section("QUERIES").list())
		default:
			return nil
		}
	}
	return nil
}

func firstNode(nodes []*treeNode) *treeNode {
	if len(nodes) == 0 {
		return nil
	}
	return nodes[0]
}

// projection returns the output column names of a query and the expression
// node producing each, in order.
func projection(q *treeNode) (names []string, nodes []*treeNode) {
	q = firstQuery(q)
	if q == nil {
		return nil, nil
	}
	if cols := q.section("PROJECTION COLUMNS"); cols != nil {
		for _, c := range cols.children {
			names = append(names, projectedName(c.text))
		}
	}
	return names, q.section("PROJECTION").list()
}

// projectedName strips the trailing type from a `name Type` line. The type
// is a single token that may carry a parenthesised argument list; the name
// may contain spaces of its own, as `plus(id, 1)` does.
func projectedName(line string) string {
	end := len(line)
	if strings.HasSuffix(line, ")") {
		depth := 0
		for end > 0 {
			end--
			if line[end] == ')' {
				depth++
			} else if line[end] == '(' {
				depth--
				if depth == 0 {
					break
				}
			}
		}
	}
	for end > 0 && isWordByte(line[end-1]) {
		end--
	}
	return strings.TrimSpace(line[:end])
}

// sourceTable resolves a COLUMN node to the table it reads from, following
// column references through subqueries and CTEs. Columns computed by an
// expression have no source and yield "".
func (t *queryTree) sourceTable(col *treeNode) string {
	for depth := 0; col != nil && col.kind == "COLUMN" && depth < 32; depth++ {
		id, err := strconv.Atoi(col.attrs["source_id"])
		if err != nil {
			return ""
		}
		src := t.byID[id]
		if src == nil {
			return ""
		}
		switch src.kind {
		case "TABLE":
			name := src.attrs["table_name"]
			if i := strings.LastIndexByte(name, '.'); i >= 0 {
				name = name[i+1:]
			}
			return name
		case "QUERY", "UNION":
			// The column is one of the subquery's output columns; keep
			// following whatever expression produces it there.
			want := col.attrs["column_name"]
			names, nodes := projection(src)
			col = nil
			for i, name := range names {
				if name == want && i < len(nodes) {
					col = nodes[i]
					break
				}
			}
		default:
			return ""
		}
	}
	return ""
}

// childrenOrNil returns a node's children, tolerating a nil node.
func (n *treeNode) childrenOrNil() []*treeNode {
	if n == nil {
		return nil
	}
	return n.children
}

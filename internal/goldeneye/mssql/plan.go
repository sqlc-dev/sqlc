package mssql

import (
	"context"
	"encoding/xml"
	"fmt"
	"strings"
)

// The estimated showplan is the one thing that says what a parameter is
// compared with or assigned to. It is compiled, not run, with SET
// SHOWPLAN_XML ON, and prints the plan as XML in which every column is a
// ColumnReference naming its schema, table and column, and every variable
// one naming only @variable. A parameter's partner is the column on the
// other side of the Compare it is an operand of, the column an Assign
// sets to it, or the column of a seek whose range expression it is. An
// expression the plan computes once and refers to by name, Expr1002, is
// followed to its DefinedValue. A parameter the query casts, CAST(@p AS
// T), has no partner: the cast says what it is, as it does to sqlc, and
// the plan shows it as a Convert the query asked for rather than one the
// server added.

// columnRef names a table column the plan refers to.
type columnRef struct {
	schema, table, column string
}

// node is one element of the showplan XML.
type node struct {
	XMLName  xml.Name
	Attrs    []xml.Attr `xml:",any,attr"`
	Children []node     `xml:",any"`
}

func (n *node) attr(name string) string {
	for _, a := range n.Attrs {
		if a.Name.Local == name {
			return a.Value
		}
	}
	return ""
}

// unbracket strips the [brackets] the plan quotes names with.
func unbracket(s string) string {
	if strings.HasPrefix(s, "[") && strings.HasSuffix(s, "]") {
		return strings.ReplaceAll(s[1:len(s)-1], "]]", "]")
	}
	return s
}

// plan is a parsed showplan with the expressions it defines by name.
type plan struct {
	root    node
	defined map[string]*node // Expr1002 -> the ScalarOperator defining it
}

// partners compiles a batch and returns the column each variable is
// compared with or assigned to, keyed by the variable's name without the
// @. A variable with no such column is absent.
func (a *analyzer) partners(ctx context.Context, batch string) (map[string]columnRef, error) {
	if _, err := a.conn.ExecContext(ctx, "SET SHOWPLAN_XML ON"); err != nil {
		return nil, err
	}
	defer a.conn.ExecContext(context.WithoutCancel(ctx), "SET SHOWPLAN_XML OFF")
	rows, err := a.conn.QueryContext(ctx, batch)
	if err != nil {
		return nil, fmt.Errorf("showplan: %w", err)
	}
	defer rows.Close()
	var blob string
	for rows.Next() {
		var s string
		if err := rows.Scan(&s); err != nil {
			return nil, err
		}
		blob += s
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	p := &plan{defined: map[string]*node{}}
	if err := xml.Unmarshal([]byte(blob), &p.root); err != nil {
		return nil, fmt.Errorf("showplan: %w", err)
	}
	p.collectDefined(&p.root)
	partners := map[string]columnRef{}
	p.walk(&p.root, partners)
	return partners, nil
}

// collectDefined records every DefinedValue that names an expression.
func (p *plan) collectDefined(n *node) {
	if n.XMLName.Local == "DefinedValue" && len(n.Children) >= 2 {
		ref := &n.Children[0]
		if ref.XMLName.Local == "ColumnReference" && ref.attr("Table") == "" && !strings.HasPrefix(ref.attr("Column"), "@") {
			p.defined[ref.attr("Column")] = &n.Children[1]
		}
	}
	for i := range n.Children {
		p.collectDefined(&n.Children[i])
	}
}

// walk finds the partners in a subtree, keeping the first found for each
// variable.
func (p *plan) walk(n *node, partners map[string]columnRef) {
	set := func(vars []string, ref columnRef) {
		for _, v := range vars {
			if _, ok := partners[v]; !ok {
				partners[v] = ref
			}
		}
	}
	switch n.XMLName.Local {
	case "Compare":
		var sides []*node
		for i := range n.Children {
			if n.Children[i].XMLName.Local == "ScalarOperator" {
				sides = append(sides, &n.Children[i])
			}
		}
		if len(sides) == 2 {
			for i, side := range sides {
				vars := p.variables(side, nil)
				if len(vars) == 0 {
					continue
				}
				if cols := p.columns(sides[1-i], nil); len(cols) > 0 {
					set(vars, cols[0])
				}
			}
		}
	case "Assign":
		var target *columnRef
		for i := range n.Children {
			c := &n.Children[i]
			switch c.XMLName.Local {
			case "ColumnReference":
				if ref, ok := refOf(c); ok && target == nil {
					target = &ref
				}
			case "ScalarOperator":
				if target != nil {
					set(p.variables(c, nil), *target)
				}
			}
		}
	case "Prefix", "StartRange", "EndRange":
		var cols []columnRef
		var exprs []*node
		for i := range n.Children {
			c := &n.Children[i]
			switch c.XMLName.Local {
			case "RangeColumns":
				for j := range c.Children {
					if ref, ok := refOf(&c.Children[j]); ok {
						cols = append(cols, ref)
					}
				}
			case "RangeExpressions":
				for j := range c.Children {
					exprs = append(exprs, &c.Children[j])
				}
			}
		}
		for i, e := range exprs {
			if i < len(cols) {
				set(p.variables(e, nil), cols[i])
			}
		}
	}
	for i := range n.Children {
		p.walk(&n.Children[i], partners)
	}
}

// refOf reads a ColumnReference that names a table column.
func refOf(n *node) (columnRef, bool) {
	if n.XMLName.Local != "ColumnReference" || n.attr("Table") == "" {
		return columnRef{}, false
	}
	return columnRef{
		schema: unbracket(n.attr("Schema")),
		table:  unbracket(n.attr("Table")),
		column: n.attr("Column"),
	}, true
}

// variables lists the variables a subtree refers to, following named
// expressions and leaving out the ones under a cast the query wrote.
func (p *plan) variables(n *node, seen map[string]bool) []string {
	var out []string
	if n.XMLName.Local == "Convert" && n.attr("Implicit") == "0" {
		return nil
	}
	if n.XMLName.Local == "ColumnReference" {
		col := n.attr("Column")
		switch {
		case strings.HasPrefix(col, "@"):
			return []string{col[1:]}
		case n.attr("Table") == "":
			if def, ok := p.defined[col]; ok && !seen[col] {
				if seen == nil {
					seen = map[string]bool{}
				}
				seen[col] = true
				return p.variables(def, seen)
			}
		}
	}
	for i := range n.Children {
		out = append(out, p.variables(&n.Children[i], seen)...)
	}
	return out
}

// columns lists the table columns a subtree refers to, following named
// expressions.
func (p *plan) columns(n *node, seen map[string]bool) []columnRef {
	var out []columnRef
	if n.XMLName.Local == "ColumnReference" {
		if ref, ok := refOf(n); ok {
			return []columnRef{ref}
		}
		col := n.attr("Column")
		if def, ok := p.defined[col]; ok && !seen[col] {
			if seen == nil {
				seen = map[string]bool{}
			}
			seen[col] = true
			return p.columns(def, seen)
		}
	}
	for i := range n.Children {
		out = append(out, p.columns(&n.Children[i], seen)...)
	}
	return out
}

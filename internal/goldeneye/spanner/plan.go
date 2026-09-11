package spanner

import (
	"regexp"
	"strings"

	"cloud.google.com/go/spanner/apiv1/spannerpb"
)

// The query plan is the one thing that says where a result column comes
// from and what a parameter stands for. It is a list of nodes: relational
// ones, such as a Scan of a table or the Serialize Result that produces the
// rows, and scalar ones, such as a Reference to a column or a variable, a
// Parameter, a Function. A relational node defines variables through its
// child links: a Scan of a table defines one per column it reads, named
// after the column, and any node may define one as another scalar, which
// a Reference names with a $. The children of Serialize Result after the
// relation it serializes are the result columns, in order; a DML plan
// lists the values it writes first — the key columns of the table, then
// for an UPDATE the columns it sets, or for an INSERT the columns it
// inserts — and the THEN RETURN columns after them. A comparison is a
// Function whose description reads ($col = @param).

// origin is what a scalar of the plan resolves to: a table column, a
// parameter, or nothing.
type origin struct {
	table, column string
	param         string
}

// plan is a parsed query plan.
type plan struct {
	nodes []*spannerpb.PlanNode
	// vars maps a variable to the node that defines it and the scalar it
	// is defined as.
	vars map[string]definition
}

type definition struct {
	owner *spannerpb.PlanNode
	child *spannerpb.PlanNode
}

func newPlan(qp *spannerpb.QueryPlan) *plan {
	p := &plan{vars: map[string]definition{}}
	if qp == nil {
		return p
	}
	p.nodes = qp.PlanNodes
	for _, n := range p.nodes {
		for _, l := range n.ChildLinks {
			if l.Variable != "" {
				p.vars[l.Variable] = definition{owner: n, child: p.node(l.ChildIndex)}
			}
		}
	}
	return p
}

func (p *plan) node(index int32) *spannerpb.PlanNode {
	for _, n := range p.nodes {
		if n.Index == index {
			return n
		}
	}
	return nil
}

func (p *plan) root() *spannerpb.PlanNode {
	return p.node(0)
}

func meta(n *spannerpb.PlanNode, key string) string {
	if n == nil || n.Metadata == nil {
		return ""
	}
	if v, ok := n.Metadata.Fields[key]; ok {
		return v.GetStringValue()
	}
	return ""
}

func description(n *spannerpb.PlanNode) string {
	if n == nil || n.ShortRepresentation == nil {
		return ""
	}
	return n.ShortRepresentation.Description
}

// serializeResult finds the node that produces the rows, the first
// Serialize Result reached from the root.
func (p *plan) serializeResult() *spannerpb.PlanNode {
	var found *spannerpb.PlanNode
	seen := map[int32]bool{}
	var walk func(n *spannerpb.PlanNode)
	walk = func(n *spannerpb.PlanNode) {
		if n == nil || found != nil || seen[n.Index] {
			return
		}
		seen[n.Index] = true
		if n.DisplayName == "Serialize Result" {
			found = n
			return
		}
		for _, l := range n.ChildLinks {
			walk(p.node(l.ChildIndex))
		}
	}
	walk(p.root())
	return found
}

// outputs lists the scalars Serialize Result produces, in order: every
// child after the relation it serializes.
func (p *plan) outputs() []*spannerpb.PlanNode {
	sr := p.serializeResult()
	if sr == nil {
		return nil
	}
	var out []*spannerpb.PlanNode
	for _, l := range sr.ChildLinks {
		if c := p.node(l.ChildIndex); c != nil && c.Kind == spannerpb.PlanNode_SCALAR {
			out = append(out, c)
		}
	}
	return out
}

// mutation names the table a DML plan writes, or "" for a query.
func (p *plan) mutation() (table, operation string) {
	for _, n := range p.nodes {
		if n.DisplayName == "Apply Mutations" {
			return meta(n, "table"), meta(n, "operation_type")
		}
	}
	return "", ""
}

// resolve follows a scalar to what it stands for.
func (p *plan) resolve(n *spannerpb.PlanNode, seen map[int32]bool) origin {
	if n == nil || seen[n.Index] {
		return origin{}
	}
	seen[n.Index] = true
	switch n.DisplayName {
	case "Parameter":
		return origin{param: meta(n, "name")}
	case "Reference":
		desc := description(n)
		if name, ok := strings.CutPrefix(desc, "$"); ok {
			return p.resolveVar(name, seen)
		}
		// A bare name is a column of the scan it is read by, or an input
		// of a union, which names it in the link's type.
		for _, m := range p.nodes {
			for _, l := range m.ChildLinks {
				if l.ChildIndex != n.Index {
					continue
				}
				switch {
				case m.DisplayName == "Scan" && meta(m, "scan_type") == "TableScan":
					return origin{table: meta(m, "scan_target"), column: desc}
				case m.DisplayName == "Scan":
					// A batch scan reads the batch variable of that name.
					target := strings.TrimPrefix(meta(m, "scan_target"), "$")
					return p.resolveVar(target+".Batch."+desc, seen)
				}
			}
		}
		for _, m := range p.nodes {
			for _, l := range m.ChildLinks {
				if l.Type == desc {
					return p.resolve(p.node(l.ChildIndex), seen)
				}
			}
		}
	}
	return origin{}
}

// resolveVar follows a variable to what it is defined as.
func (p *plan) resolveVar(name string, seen map[int32]bool) origin {
	d, ok := p.vars[name]
	if !ok {
		return origin{}
	}
	return p.resolve(d.child, seen)
}

var compareRe = regexp.MustCompile(`^\((.+) (=|!=|<>|<|<=|>|>=) (.+)\)$`)

// partners finds the table column each parameter is compared with, keyed
// by parameter name: the other operand of a comparison Function the
// parameter is an operand of. The first found wins.
func (p *plan) partners() map[string]origin {
	out := map[string]origin{}
	for _, n := range p.nodes {
		if n.DisplayName != "Function" || len(n.ChildLinks) != 2 || !compareRe.MatchString(description(n)) {
			continue
		}
		sides := []*spannerpb.PlanNode{p.node(n.ChildLinks[0].ChildIndex), p.node(n.ChildLinks[1].ChildIndex)}
		for i, side := range sides {
			o := p.resolve(side, map[int32]bool{})
			if o.param == "" {
				continue
			}
			other := p.resolve(sides[1-i], map[int32]bool{})
			if other.table != "" {
				if _, ok := out[o.param]; !ok {
					out[o.param] = other
				}
			}
		}
	}
	return out
}

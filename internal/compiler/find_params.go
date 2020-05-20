package compiler

import (
	"github.com/kyleconroy/sqlc/internal/sql/ast"
	"github.com/kyleconroy/sqlc/internal/sql/ast/pg"
	"github.com/kyleconroy/sqlc/internal/sql/astutils"
)

func findParameters(root ast.Node) []paramRef {
	refs := make([]paramRef, 0)
	v := paramSearch{seen: make(map[int]struct{}), refs: &refs}
	astutils.Walk(v, root)
	return refs
}

type paramRef struct {
	parent ast.Node
	rv     *pg.RangeVar
	ref    *pg.ParamRef
	name   string // Named parameter support
}

type paramSearch struct {
	parent   ast.Node
	rangeVar *pg.RangeVar
	refs     *[]paramRef
	seen     map[int]struct{}

	// XXX: Gross state hack for limit
	limitCount  ast.Node
	limitOffset ast.Node
}

type limitCount struct {
}

func (l *limitCount) Pos() int {
	return 0
}

type limitOffset struct {
}

func (l *limitOffset) Pos() int {
	return 0
}

func (p paramSearch) Visit(node ast.Node) astutils.Visitor {
	switch n := node.(type) {

	case *pg.A_Expr:
		p.parent = node

	case *ast.FuncCall:
		p.parent = node

	case *pg.InsertStmt:
		if s, ok := n.SelectStmt.(*pg.SelectStmt); ok {
			for i, item := range s.TargetList.Items {
				target, ok := item.(*pg.ResTarget)
				if !ok {
					continue
				}
				ref, ok := target.Val.(*pg.ParamRef)
				if !ok {
					continue
				}
				// TODO: Out-of-bounds panic
				*p.refs = append(*p.refs, paramRef{parent: n.Cols.Items[i], ref: ref, rv: n.Relation})
				p.seen[ref.Location] = struct{}{}
			}
			for _, item := range s.ValuesLists.Items {
				vl, ok := item.(*ast.List)
				if !ok {
					continue
				}
				for i, v := range vl.Items {
					ref, ok := v.(*pg.ParamRef)
					if !ok {
						continue
					}
					// TODO: Out-of-bounds panic
					*p.refs = append(*p.refs, paramRef{parent: n.Cols.Items[i], ref: ref, rv: n.Relation})
					p.seen[ref.Location] = struct{}{}
				}
			}
		}

	case *pg.RangeVar:
		p.rangeVar = n

	case *pg.ResTarget:
		p.parent = node

	case *pg.SelectStmt:
		if n.LimitCount != nil {
			p.limitCount = n.LimitCount
		}
		if n.LimitOffset != nil {
			p.limitOffset = n.LimitOffset
		}

	case *pg.TypeCast:
		p.parent = node

	case *pg.ParamRef:
		parent := p.parent

		if count, ok := p.limitCount.(*pg.ParamRef); ok {
			if n.Number == count.Number {
				parent = &limitCount{}
			}
		}

		if offset, ok := p.limitOffset.(*pg.ParamRef); ok {
			if n.Number == offset.Number {
				parent = &limitOffset{}
			}
		}
		if _, found := p.seen[n.Location]; found {
			break
		}

		// Special, terrible case for *pg.MultiAssignRef
		set := true
		if res, ok := parent.(*pg.ResTarget); ok {
			if multi, ok := res.Val.(*pg.MultiAssignRef); ok {
				set = false
				if row, ok := multi.Source.(*pg.RowExpr); ok {
					for i, arg := range row.Args.Items {
						if ref, ok := arg.(*pg.ParamRef); ok {
							if multi.Colno == i+1 && ref.Number == n.Number {
								set = true
							}
						}
					}
				}
			}
		}

		if set {
			*p.refs = append(*p.refs, paramRef{parent: parent, ref: n, rv: p.rangeVar})
			p.seen[n.Location] = struct{}{}
		}
		return nil
	}
	return p
}

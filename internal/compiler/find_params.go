package compiler

import (
	"fmt"

	"github.com/sqlc-dev/sqlc/internal/sql/ast"
	"github.com/sqlc-dev/sqlc/internal/sql/astutils"
)

func findParameters(root ast.Node) ([]paramRef, []error) {
	refs := make([]paramRef, 0)
	errors := make([]error, 0)
	v := paramSearch{seen: make(map[int]struct{}), refs: &refs, errs: &errors}
	astutils.Walk(v, root)
	if len(*v.errs) > 0 {
		return refs, *v.errs
	} else {
		return refs, nil
	}
}

type paramRef struct {
	parent ast.Node
	rv     *ast.RangeVar
	ref    *ast.ParamRef
	name   string // Named parameter support
}

type paramSearch struct {
	parent   ast.Node
	rangeVar *ast.RangeVar
	refs     *[]paramRef
	seen     map[int]struct{}
	errs     *[]error

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

	case *ast.A_Expr:
		p.parent = node

	case *ast.BetweenExpr:
		p.parent = node

	case *ast.CallStmt:
		p.parent = n.FuncCall

	case *ast.DeleteStmt:
		if n.LimitCount != nil {
			p.limitCount = n.LimitCount
		}

	case *ast.FuncCall:
		p.parent = node

	case *ast.InsertStmt:
		if s, ok := n.SelectStmt.(*ast.SelectStmt); ok {
			for i, item := range s.TargetList.Items {
				target, ok := item.(*ast.ResTarget)
				if !ok {
					continue
				}
				ref, ok := target.Val.(*ast.ParamRef)
				if !ok {
					continue
				}
				if len(n.Cols.Items) <= i {
					*p.errs = append(*p.errs, fmt.Errorf("INSERT has more expressions than target columns"))
					return p
				}
				*p.refs = append(*p.refs, paramRef{parent: n.Cols.Items[i], ref: ref, rv: n.Relation})
				p.seen[ref.Location] = struct{}{}
			}
			for _, item := range s.ValuesLists.Items {
				vl, ok := item.(*ast.List)
				if !ok {
					continue
				}
				for i, v := range vl.Items {
					ref, ok := v.(*ast.ParamRef)
					if !ok {
						continue
					}
					if len(n.Cols.Items) <= i {
						*p.errs = append(*p.errs, fmt.Errorf("INSERT has more expressions than target columns"))
						return p
					}
					*p.refs = append(*p.refs, paramRef{parent: n.Cols.Items[i], ref: ref, rv: n.Relation})
					p.seen[ref.Location] = struct{}{}
				}
			}
		}

	case *ast.UpdateStmt:
		for _, item := range n.TargetList.Items {
			target, ok := item.(*ast.ResTarget)
			if !ok {
				continue
			}
			ref, ok := target.Val.(*ast.ParamRef)
			if !ok {
				continue
			}
			for _, relation := range n.Relations.Items {
				rv, ok := relation.(*ast.RangeVar)
				if !ok {
					continue
				}
				*p.refs = append(*p.refs, paramRef{parent: target, ref: ref, rv: rv})
			}
			p.seen[ref.Location] = struct{}{}
		}
		if n.LimitCount != nil {
			p.limitCount = n.LimitCount
		}

	case *ast.RangeVar:
		p.rangeVar = n

	case *ast.ResTarget:
		p.parent = node

	case *ast.SelectStmt:
		if n.LimitCount != nil {
			p.limitCount = n.LimitCount
		}
		if n.LimitOffset != nil {
			p.limitOffset = n.LimitOffset
		}

	case *ast.TypeCast:
		p.parent = node

	case *ast.ParamRef:
		parent := p.parent

		if count, ok := p.limitCount.(*ast.ParamRef); ok {
			if n.Number == count.Number {
				parent = &limitCount{}
			}
		}

		if offset, ok := p.limitOffset.(*ast.ParamRef); ok {
			if n.Number == offset.Number {
				parent = &limitOffset{}
			}
		}
		if _, found := p.seen[n.Location]; found {
			break
		}

		// Special, terrible case for *ast.MultiAssignRef
		set := true
		if res, ok := parent.(*ast.ResTarget); ok {
			if multi, ok := res.Val.(*ast.MultiAssignRef); ok {
				set = false
				if row, ok := multi.Source.(*ast.RowExpr); ok {
					for i, arg := range row.Args.Items {
						if ref, ok := arg.(*ast.ParamRef); ok {
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

	case *ast.In:
		if n.Sel == nil {
			p.parent = node
		} else {
			if sel, ok := n.Sel.(*ast.SelectStmt); ok && sel.FromClause != nil && len(sel.FromClause.Items) > 0 {
				from := sel.FromClause
				if schema, ok := from.Items[0].(*ast.RangeVar); ok && schema != nil {
					p.rangeVar = &ast.RangeVar{
						Catalogname: schema.Catalogname,
						Schemaname:  schema.Schemaname,
						Relname:     schema.Relname,
					}
				}
			}
		}
		if _, ok := n.Expr.(*ast.ParamRef); ok {
			p.Visit(n.Expr)
		}
	}
	return p
}

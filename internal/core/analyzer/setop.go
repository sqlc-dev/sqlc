package analyzer

import (
	"github.com/sqlc-dev/sqlc/internal/core"
	"github.com/sqlc-dev/sqlc/internal/sql/ast"
)

// combineColumns types each column of a set operation from both of its
// branches by the dialect's rules. A column keeps the first branch's name
// and is read from no one table. It may be NULL in a UNION when either
// branch's may, in an INTERSECT when both may, and in an EXCEPT when the
// first branch's may, since every row it returns is one of those.
func (a *analyzer) combineColumns(op ast.SetOperation, left, right []core.Column) ([]core.Column, error) {
	context := core.CombineUnion
	switch op {
	case ast.Intersect:
		context = core.CombineIntersect
	case ast.Except:
		context = core.CombineExcept
	}
	out := make([]core.Column, len(left))
	for i, l := range left {
		if i >= len(right) {
			out[i] = l
			continue
		}
		r := right[i]
		col := core.Column{Name: l.Name, NotNull: l.NotNull && r.NotNull}
		switch op {
		case ast.Intersect:
			col.NotNull = l.NotNull || r.NotNull
		case ast.Except:
			col.NotNull = l.NotNull
		}
		t, handled, err := a.combineTypes(context, []exprType{columnExprType(l), columnExprType(r)})
		if err != nil {
			return nil, err
		}
		if !handled {
			out[i] = l
			continue
		}
		t.nullable = !col.NotNull
		col.TypeOID = t.typeOID
		col.DataType, col.IsArray = a.typeNameOf(t)
		col.Type = a.typeExprOf(t)
		out[i] = col
	}
	return out, nil
}

// combineTypes delegates selection and modifiers to engine-supplied behavior.
// The whole list is passed without discarding untyped candidate positions.
func (a *analyzer) combineTypes(context core.TypeCombinationContext, types []exprType) (exprType, bool, error) {
	if a.combiner == nil {
		return exprType{}, false, nil
	}
	candidates := make([]*core.TypeExpr, len(types))
	typed := 0
	for i, t := range types {
		candidates[i] = a.exprOf(t)
		if candidates[i] != nil {
			typed++
		}
	}
	result, err := a.combiner.Combine(core.TypeCombinationRequest{Context: context, Candidates: candidates})
	if err != nil {
		return exprType{}, false, err
	}
	if !result.Handled {
		return exprType{}, false, nil
	}
	out := result.Type
	if out == nil {
		return exprType{}, true, nil
	}
	// Preserve the existing column provenance when only the first candidate
	// is typed and the combiner leaves its type unchanged.
	if typed == 1 && candidates[0] != nil && out.Key() == candidates[0].Key() {
		return types[0], true, nil
	}
	// A type with no arguments is its family. Looked up as a declaration
	// it would take the dialect's defaults, as varchar is varchar(1).
	if len(out.Args) == 0 {
		if oid, err := a.cat.TypeOID(out.Name); err == nil {
			return exprType{typeOID: oid}, true, nil
		}
	}
	return a.lookupType(out), true, nil
}

// commonType retains first-typed-value inference when no combiner handles
// the context. This compatibility fallback is not database coercion semantics.
func (a *analyzer) commonType(context core.TypeCombinationContext, types []exprType) (exprType, error) {
	out, handled, err := a.combineTypes(context, types)
	if err != nil || handled {
		return out, err
	}
	return firstTyped(types), nil
}

func firstTyped(types []exprType) exprType {
	for i, t := range types {
		if t.typeOID == 0 && t.expr == nil {
			continue
		}
		if i > 0 {
			t.sourceClassOID, t.sourceAttributeOID, t.sourceTableAlias = 0, 0, ""
		}
		return t
	}
	return exprType{}
}

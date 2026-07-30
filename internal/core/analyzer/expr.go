package analyzer

import (
	"fmt"
	"strings"

	"github.com/sqlc-dev/sqlc/internal/core"
	"github.com/sqlc-dev/sqlc/internal/sql/ast"
)

type exprType struct {
	typeOID            int64
	nullable           bool
	sourceClassOID     int64
	sourceAttributeOID int64
	sourceTableAlias   string
}

func (a *analyzer) typeExpr(n ast.Node) (exprType, error) {
	switch e := n.(type) {
	case nil:
		return exprType{}, nil

	case *ast.TODO:
		return exprType{}, nil

	case *ast.A_Const:
		return a.typeConst(e)

	case *ast.ColumnRef:
		return a.typeColumnRef(e)

	case *ast.ParamRef:
		return a.typeParamRef(e)

	case *ast.A_Expr:
		return a.typeAExpr(e)

	case *ast.BoolExpr:
		return a.typeBoolExpr(e)

	case *ast.FuncCall:
		return a.typeFuncCall(e)

	case *ast.TypeCast:
		return a.typeTypeCast(e)

	case *ast.NullTest:
		if _, err := a.typeExpr(e.Arg); err != nil {
			return exprType{}, err
		}
		return a.boolType(false)
	}
	return exprType{}, fmt.Errorf("typeExpr: unsupported %T", n)
}

func (a *analyzer) typeConst(c *ast.A_Const) (exprType, error) {
	switch v := c.Val.(type) {
	case *ast.Integer:
		_ = v
		oid, err := a.cat.TypeOID("int4")
		if err != nil {
			return exprType{}, err
		}
		return exprType{typeOID: oid}, nil
	case *ast.Float:
		oid, err := a.cat.TypeOID("numeric")
		if err != nil {
			return exprType{}, err
		}
		return exprType{typeOID: oid}, nil
	case *ast.String:
		oid, err := a.cat.TypeOID("text")
		if err != nil {
			return exprType{}, err
		}
		return exprType{typeOID: oid}, nil
	case *ast.Boolean:
		return a.boolType(false)
	case nil:
		return exprType{nullable: true}, nil
	}
	return exprType{}, fmt.Errorf("typeConst: unsupported %T", c.Val)
}

func (a *analyzer) boolType(nullable bool) (exprType, error) {
	oid, err := a.cat.TypeOID("bool")
	if err != nil {
		return exprType{}, err
	}
	return exprType{typeOID: oid, nullable: nullable}, nil
}

func (a *analyzer) typeColumnRef(c *ast.ColumnRef) (exprType, error) {
	parts := flattenFields(c.Fields)
	if len(parts) == 0 {
		return exprType{}, fmt.Errorf("column ref: empty")
	}
	relation := ""
	column := parts[0]
	if len(parts) >= 2 {
		relation = parts[0]
		column = parts[1]
	}
	rel, col, ok, err := a.scope.resolveColumn(relation, column)
	if err != nil {
		return exprType{}, err
	}
	if !ok {
		if relation != "" {
			return exprType{}, fmt.Errorf("unknown column %q.%q", relation, column)
		}
		return exprType{}, fmt.Errorf("unknown column %q", column)
	}
	return exprType{
		typeOID:            col.typeOID,
		nullable:           !col.notNull,
		sourceClassOID:     rel.classOID,
		sourceAttributeOID: col.attOID,
		sourceTableAlias:   rel.alias,
	}, nil
}

func flattenFields(fields *ast.List) []string {
	if fields == nil {
		return nil
	}
	out := make([]string, 0, len(fields.Items))
	for _, item := range fields.Items {
		switch v := item.(type) {
		case *ast.String:
			out = append(out, v.Str)
		case *ast.A_Star:
			out = append(out, "*")
			return out
		}
	}
	return out
}

func (a *analyzer) typeParamRef(p *ast.ParamRef) (exprType, error) {
	cur, ok := a.params[p.Number]
	if !ok {
		cur = &core.Parameter{Number: p.Number}
		a.params[p.Number] = cur
	}
	return exprType{typeOID: cur.TypeOID, nullable: !cur.NotNull}, nil
}

func (a *analyzer) inferParam(number int, t exprType) {
	cur, ok := a.params[number]
	if !ok {
		cur = &core.Parameter{Number: number}
		a.params[number] = cur
	}
	if cur.TypeOID == 0 && t.typeOID != 0 {
		cur.TypeOID = t.typeOID
		if name, err := a.cat.TypeName(t.typeOID); err == nil {
			cur.DataType = name
		}
		cur.NotNull = !t.nullable
	}
	if cur.Source == nil && t.sourceAttributeOID != 0 {
		ad, err := a.cat.LookupAttribute(t.sourceAttributeOID)
		if err == nil {
			cur.Source = &core.ColumnSource{
				Schema:     ad.Schema,
				Table:      ad.Table,
				TableAlias: t.sourceTableAlias,
				Column:     ad.Column,
			}
		}
	}
}

func (a *analyzer) typeAExpr(e *ast.A_Expr) (exprType, error) {
	if e.Kind != ast.A_Expr_Kind_OP {
		return exprType{}, fmt.Errorf("a_expr: unsupported kind %v", e.Kind)
	}
	opName := opNameFromList(e.Name)
	if opName == "" {
		return exprType{}, fmt.Errorf("a_expr: unnamed operator")
	}

	leftT, err := a.typeExpr(e.Lexpr)
	if err != nil {
		return exprType{}, err
	}
	rightT, err := a.typeExpr(e.Rexpr)
	if err != nil {
		return exprType{}, err
	}

	if pr, ok := e.Lexpr.(*ast.ParamRef); ok && rightT.typeOID != 0 {
		a.inferParam(pr.Number, rightT)
		leftT = rightT
	}
	if pr, ok := e.Rexpr.(*ast.ParamRef); ok && leftT.typeOID != 0 {
		a.inferParam(pr.Number, leftT)
		rightT = leftT
	}

	overload, err := a.resolveOperator(opName, leftT.typeOID, rightT.typeOID)
	if err != nil {
		return exprType{}, err
	}
	return exprType{
		typeOID:  overload.ResultTypeOID,
		nullable: leftT.nullable || rightT.nullable,
	}, nil
}

func opNameFromList(l *ast.List) string {
	if l == nil {
		return ""
	}
	parts := make([]string, 0, len(l.Items))
	for _, item := range l.Items {
		if s, ok := item.(*ast.String); ok {
			parts = append(parts, s.Str)
		}
	}
	return strings.Join(parts, ".")
}

func (a *analyzer) resolveOperator(name string, leftOID, rightOID int64) (core.OperatorOverload, error) {
	candidates, err := a.cat.FindOperators(name, leftOID, rightOID)
	if err != nil {
		return core.OperatorOverload{}, err
	}
	if len(candidates) > 0 {
		return candidates[0], nil
	}

	all, err := a.cat.FindOperators(name, 0, 0)
	if err != nil {
		return core.OperatorOverload{}, err
	}
	for _, ov := range all {
		if leftOID != 0 && ov.LeftTypeOID != 0 && leftOID != ov.LeftTypeOID {
			ok, _ := a.cat.CastAllowed(leftOID, ov.LeftTypeOID, "i")
			if !ok {
				continue
			}
		}
		if rightOID != 0 && ov.RightTypeOID != 0 && rightOID != ov.RightTypeOID {
			ok, _ := a.cat.CastAllowed(rightOID, ov.RightTypeOID, "i")
			if !ok {
				continue
			}
		}
		if (leftOID == 0) != (ov.LeftTypeOID == 0) {
			continue
		}
		if (rightOID == 0) != (ov.RightTypeOID == 0) {
			continue
		}
		return ov, nil
	}
	return core.OperatorOverload{}, fmt.Errorf("no operator %q for (%d, %d)", name, leftOID, rightOID)
}

func (a *analyzer) typeBoolExpr(b *ast.BoolExpr) (exprType, error) {
	for _, item := range listItems(b.Args) {
		if _, err := a.typeExpr(item); err != nil {
			return exprType{}, err
		}
	}
	return a.boolType(false)
}

func (a *analyzer) typeFuncCall(f *ast.FuncCall) (exprType, error) {
	name := funcCallName(f)
	if name == "" {
		return exprType{}, fmt.Errorf("func call: missing name")
	}

	if f.AggStar && (name == "count" || name == "count.*") {
		if overloads, err := a.cat.FindProcs("count", nil); err == nil && len(overloads) > 0 {
			return exprType{typeOID: overloads[0].ReturnTypeOID, nullable: overloads[0].ReturnNullable}, nil
		}
		oid, err := a.cat.TypeOID("int8")
		if err != nil {
			return exprType{}, err
		}
		return exprType{typeOID: oid, nullable: false}, nil
	}

	for _, arg := range listItems(f.Args) {
		if _, err := a.typeExpr(arg); err != nil {
			return exprType{}, err
		}
	}

	overloads, err := a.cat.FindProcs(name, nil)
	if err != nil {
		return exprType{}, err
	}
	if len(overloads) == 0 {
		return exprType{}, fmt.Errorf("unknown function %q", name)
	}
	p := overloads[0]
	return exprType{typeOID: p.ReturnTypeOID, nullable: p.ReturnNullable}, nil
}

func funcCallName(f *ast.FuncCall) string {
	if f.Funcname != nil {
		return opNameFromList(f.Funcname)
	}
	if f.Func != nil {
		return strings.ToLower(f.Func.Name)
	}
	return ""
}

func (a *analyzer) typeTypeCast(c *ast.TypeCast) (exprType, error) {
	if c.TypeName == nil {
		return exprType{}, fmt.Errorf("cast: missing target type")
	}
	if _, err := a.typeExpr(c.Arg); err != nil {
		return exprType{}, err
	}
	oid, err := a.cat.TypeOID(strings.ToLower(c.TypeName.Name))
	if err != nil {
		return exprType{}, fmt.Errorf("cast target %q: %w", c.TypeName.Name, err)
	}
	return exprType{typeOID: oid}, nil
}

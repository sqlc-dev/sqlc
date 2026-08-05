package analyzer

import (
	"slices"

	"github.com/sqlc-dev/sqlc/internal/core"
	"github.com/sqlc-dev/sqlc/internal/sql/ast"
)

func (a *analyzer) projectTarget(rt *ast.ResTarget) error {
	// A column reference's field list is flattened once here and threaded
	// through the star check, the star expansion and the output name.
	var fields []string
	if cr, ok := rt.Val.(*ast.ColumnRef); ok {
		fields = flattenFields(cr.Fields)
		if isStar(fields) {
			a.emitStar(fields)
			return nil
		}
	}

	t, err := a.typeExpr(rt.Val)
	if err != nil {
		return err
	}
	col := core.Column{
		Name:               targetName(rt, fields),
		TypeOID:            t.typeOID,
		NotNull:            !t.nullable,
		SourceClassOID:     t.sourceClassOID,
		SourceAttributeOID: t.sourceAttributeOID,
	}
	col.DataType, col.IsArray = a.typeNameOf(t)
	a.decorateSource(&col, t.sourceAttributeOID, t.sourceTableAlias)
	a.columns = append(a.columns, col)
	return nil
}

func (a *analyzer) decorateSource(col *core.Column, attOID int64, tableAlias string) {
	if attOID == 0 {
		return
	}
	ad, err := a.cat.LookupAttribute(attOID)
	if err != nil {
		return
	}
	col.Source = &core.ColumnSource{
		Schema:     ad.Schema,
		Table:      ad.Table,
		TableAlias: tableAlias,
		Column:     ad.Column,
	}
	col.DeclType = ad.DeclType
	col.TypeLength = ad.TypeLength
	col.TypeScale = ad.TypeScale
	col.IsPrimaryKey = ad.IsPrimaryKey
	col.IsUnique = ad.IsUnique
	col.IsAutoIncrement = ad.AutoIncrement
}

// targetName picks the output name for a target. fields is the already
// flattened field list when rt.Val is a column reference, and nil otherwise.
func targetName(rt *ast.ResTarget, fields []string) string {
	if rt.Name != nil && *rt.Name != "" {
		return *rt.Name
	}
	if len(fields) > 0 {
		return fields[len(fields)-1]
	}
	if fc, ok := rt.Val.(*ast.FuncCall); ok {
		if name := funcCallName(fc); name != "" {
			return name
		}
	}
	return "?column?"
}

func isStar(fields []string) bool {
	return len(fields) > 0 && fields[len(fields)-1] == "*"
}

func (a *analyzer) emitStar(fields []string) {
	relName := ""
	if len(fields) > 1 {
		relName = fields[0]
	}
	for _, rel := range a.scope.rels {
		if relName != "" && rel.alias != relName {
			continue
		}
		a.columns = slices.Grow(a.columns, len(rel.cols))
		for _, c := range rel.cols {
			col := core.Column{
				Name:               c.Name,
				TypeOID:            c.TypeOID,
				NotNull:            c.NotNull,
				SourceClassOID:     rel.classOID,
				SourceAttributeOID: c.AttOID,
			}
			col.DataType, col.IsArray = a.typeNameOf(exprType{typeOID: c.TypeOID})
			a.decorateSource(&col, c.AttOID, rel.alias)
			a.columns = append(a.columns, col)
		}
	}
}

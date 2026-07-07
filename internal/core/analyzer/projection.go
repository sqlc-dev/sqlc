package analyzer

import (
	"github.com/sqlc-dev/sqlc/internal/core"
	"github.com/sqlc-dev/sqlc/internal/sql/ast"
)

// projectTarget evaluates one entry in the SELECT list and records an
// output Column. Star expansions emit one column per source column.
func (a *analyzer) projectTarget(rt *ast.ResTarget) error {
	if cr, ok := rt.Val.(*ast.ColumnRef); ok {
		// "* " or "t.*" expansion
		if isStarRef(cr) {
			a.emitStar(cr)
			return nil
		}
	}

	t, err := a.typeExpr(rt.Val)
	if err != nil {
		return err
	}
	col := core.Column{
		Name:               targetName(rt),
		TypeOID:            t.typeOID,
		NotNull:            !t.nullable,
		SourceClassOID:     t.sourceClassOID,
		SourceAttributeOID: t.sourceAttributeOID,
	}
	if t.typeOID != 0 {
		if name, err := a.cat.TypeName(t.typeOID); err == nil {
			col.DataType = name
		}
	}
	a.decorateSource(&col, t.sourceAttributeOID, t.sourceTableAlias)
	a.columns = append(a.columns, col)
	return nil
}

// decorateSource fills in the resolved Source struct and per-column
// flag fields for a Column whose value came from a real table column.
// Computed expressions (zero attOID) are left untouched.
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

// targetName picks the user-visible name for a SELECT-list entry: the
// alias if one was given, the source column name for direct refs, or
// the conventional "?column?" placeholder for computed expressions.
func targetName(rt *ast.ResTarget) string {
	if rt.Name != nil && *rt.Name != "" {
		return *rt.Name
	}
	if cr, ok := rt.Val.(*ast.ColumnRef); ok {
		parts := flattenFields(cr.Fields)
		if len(parts) > 0 {
			return parts[len(parts)-1]
		}
	}
	if fc, ok := rt.Val.(*ast.FuncCall); ok {
		if name := funcCallName(fc); name != "" {
			return name
		}
	}
	return "?column?"
}

func isStarRef(c *ast.ColumnRef) bool {
	parts := flattenFields(c.Fields)
	return len(parts) > 0 && parts[len(parts)-1] == "*"
}

// emitStar expands * (or t.*) into one Column per source attribute.
func (a *analyzer) emitStar(cr *ast.ColumnRef) {
	parts := flattenFields(cr.Fields)
	relName := ""
	if len(parts) > 1 {
		relName = parts[0]
	}
	for _, rel := range a.scope.rels {
		if relName != "" && rel.alias != relName {
			continue
		}
		for _, c := range rel.cols {
			col := core.Column{
				Name:               c.name,
				TypeOID:            c.typeOID,
				NotNull:            c.notNull,
				SourceClassOID:     rel.classOID,
				SourceAttributeOID: c.attOID,
			}
			if name, err := a.cat.TypeName(c.typeOID); err == nil {
				col.DataType = name
			}
			a.decorateSource(&col, c.attOID, rel.alias)
			a.columns = append(a.columns, col)
		}
	}
}

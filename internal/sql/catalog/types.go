package catalog

import (
	"errors"
	"fmt"

	"github.com/sqlc-dev/sqlc/internal/sql/ast"
	"github.com/sqlc-dev/sqlc/internal/sql/sqlerr"
)

type Type interface {
	isType()

	SetComment(string)
}

type Enum struct {
	Name    string
	Vals    []string
	Comment string
}

func (e *Enum) SetComment(c string) {
	e.Comment = c
}

func (e *Enum) isType() {
}

type CompositeType struct {
	Name    string
	Columns []*Column
	Comment string
}

func (ct *CompositeType) isType() {
}

func (ct *CompositeType) SetComment(c string) {
	ct.Comment = c
}

func sameType(a, b *ast.TypeName) bool {
	if a.Catalog != b.Catalog {
		return false
	}
	// The pg_catalog schema is searched by default, so take that into
	// account when comparing schemas
	aSchema := a.Schema
	bSchema := b.Schema
	if aSchema == "pg_catalog" {
		aSchema = ""
	}
	if bSchema == "pg_catalog" {
		bSchema = ""
	}
	if aSchema != bSchema {
		return false
	}
	if a.Name != b.Name {
		return false
	}
	return true
}

func (c *Catalog) createEnum(stmt *ast.CreateEnumStmt) error {
	ns := stmt.TypeName.Schema
	if ns == "" {
		ns = c.DefaultSchema
	}
	schema, err := c.getSchema(ns)
	if err != nil {
		return err
	}
	// Because tables have associated data types, the type name must also
	// be distinct from the name of any existing table in the same
	// schema.
	// https://www.postgresql.org/docs/current/sql-createtype.html
	tbl := &ast.TableName{
		Name: stmt.TypeName.Name,
	}
	if _, _, err := schema.getTable(tbl); err == nil {
		return sqlerr.RelationExists(tbl.Name)
	}
	if _, _, err := schema.getType(stmt.TypeName); err == nil {
		return sqlerr.TypeExists(tbl.Name)
	}
	schema.Types = append(schema.Types, &Enum{
		Name: stmt.TypeName.Name,
		Vals: stringSlice(stmt.Vals),
	})
	return nil
}

func stringSlice(list *ast.List) []string {
	items := []string{}
	for _, item := range list.Items {
		if n, ok := item.(*ast.String); ok {
			items = append(items, n.Str)
		}
	}
	return items
}

func columnDefsSlice(list *ast.List) []*ast.ColumnDef {
	items := []*ast.ColumnDef{}
	for _, item := range list.Items {
		if n, ok := item.(*ast.ColumnDef); ok {
			items = append(items, n)
		}
	}
	return items
}

func (c *Catalog) getType(rel *ast.TypeName) (Type, int, error) {
	ns := rel.Schema
	if ns == "" {
		ns = c.DefaultSchema
	}
	s, err := c.getSchema(ns)
	if err != nil {
		return nil, -1, err
	}
	return s.getType(rel)
}

func (c *Catalog) createCompositeType(stmt *ast.CompositeTypeStmt) error {
	ns := stmt.TypeName.Schema
	if ns == "" {
		ns = c.DefaultSchema
	}
	schema, err := c.getSchema(ns)
	if err != nil {
		return err
	}
	// Because tables have associated data types, the type name must also
	// be distinct from the name of any existing table in the same
	// schema.
	// https://www.postgresql.org/docs/current/sql-createtype.html
	tbl := &ast.TableName{
		Name: stmt.TypeName.Name,
	}
	if _, _, err := schema.getTable(tbl); err == nil {
		return sqlerr.RelationExists(tbl.Name)
	}
	if _, _, err := schema.getType(stmt.TypeName); err == nil {
		return sqlerr.TypeExists(tbl.Name)
	}
	ct := &CompositeType{
		Name: stmt.TypeName.Name,
	}
	for _, col := range columnDefsSlice(stmt.ColDefList) {
		ct.Columns = append(ct.Columns, &Column{
			Name:       col.Colname,
			Type:       *col.TypeName,
			IsNotNull:  col.IsNotNull,
			IsUnsigned: col.IsUnsigned,
			IsArray:    col.IsArray,
			ArrayDims:  col.ArrayDims,
			Comment:    col.Comment,
			Length:     col.Length,
		})
	}
	schema.Types = append(schema.Types, ct)
	return nil
}

func (c *Catalog) alterTypeRenameValue(stmt *ast.AlterTypeRenameValueStmt) error {
	ns := stmt.Type.Schema
	if ns == "" {
		ns = c.DefaultSchema
	}
	schema, err := c.getSchema(ns)
	if err != nil {
		return err
	}
	typ, _, err := schema.getType(stmt.Type)
	if err != nil {
		return err
	}
	enum, ok := typ.(*Enum)
	if !ok {
		return fmt.Errorf("type is not an enum: %T", stmt.Type)
	}

	oldIndex := -1
	newIndex := -1
	for i, val := range enum.Vals {
		if val == *stmt.OldValue {
			oldIndex = i
		}
		if val == *stmt.NewValue {
			newIndex = i
		}
	}
	if oldIndex < 0 {
		return fmt.Errorf("type %T does not have value %s", stmt.Type, *stmt.OldValue)
	}
	if newIndex >= 0 {
		return fmt.Errorf("type %T already has value %s", stmt.Type, *stmt.NewValue)
	}
	enum.Vals[oldIndex] = *stmt.NewValue
	return nil
}

func (c *Catalog) alterTypeAddValue(stmt *ast.AlterTypeAddValueStmt) error {
	ns := stmt.Type.Schema
	if ns == "" {
		ns = c.DefaultSchema
	}
	schema, err := c.getSchema(ns)
	if err != nil {
		return err
	}
	typ, _, err := schema.getType(stmt.Type)
	if err != nil {
		return err
	}
	enum, ok := typ.(*Enum)
	if !ok {
		return fmt.Errorf("type is not an enum: %T", stmt.Type)
	}

	existingIndex := -1
	for i, val := range enum.Vals {
		if val == *stmt.NewValue {
			existingIndex = i
		}
	}

	if existingIndex >= 0 {
		if !stmt.SkipIfNewValExists {
			return fmt.Errorf("enum %s already has value %s", enum.Name, *stmt.NewValue)
		} else {
			return nil
		}
	}

	insertIndex := len(enum.Vals)
	if stmt.NewValHasNeighbor {
		foundNeighbor := false
		for i, val := range enum.Vals {
			if val == *stmt.NewValNeighbor {
				if stmt.NewValIsAfter {
					insertIndex = i + 1
				} else {
					insertIndex = i
				}
				foundNeighbor = true
				break
			}
		}

		if !foundNeighbor {
			return fmt.Errorf("enum %s unable to find existing neighbor value %s for new value %s", enum.Name, *stmt.NewValNeighbor, *stmt.NewValue)
		}
	}

	if insertIndex == len(enum.Vals) {
		enum.Vals = append(enum.Vals, *stmt.NewValue)
	} else {
		enum.Vals = append(enum.Vals[:insertIndex+1], enum.Vals[insertIndex:]...)
		enum.Vals[insertIndex] = *stmt.NewValue
	}

	return nil
}

func (c *Catalog) alterTypeSetSchema(stmt *ast.AlterTypeSetSchemaStmt) error {
	ns := stmt.Type.Schema
	if ns == "" {
		ns = c.DefaultSchema
	}
	oldSchema, err := c.getSchema(ns)
	if err != nil {
		return err
	}
	typ, idx, err := oldSchema.getType(stmt.Type)
	if err != nil {
		return err
	}
	oldType := *stmt.Type
	stmt.Type.Schema = *stmt.NewSchema
	newSchema, err := c.getSchema(*stmt.NewSchema)
	if err != nil {
		return err
	}
	// Because tables have associated data types, the type name must also
	// be distinct from the name of any existing table in the same
	// schema.
	// https://www.postgresql.org/docs/current/sql-createtype.html
	tbl := &ast.TableName{
		Name: stmt.Type.Name,
	}
	if _, _, err := newSchema.getTable(tbl); err == nil {
		return sqlerr.RelationExists(tbl.Name)
	}
	if _, _, err := newSchema.getType(stmt.Type); err == nil {
		return sqlerr.TypeExists(stmt.Type.Name)
	}
	oldSchema.Types = append(oldSchema.Types[:idx], oldSchema.Types[idx+1:]...)
	newSchema.Types = append(newSchema.Types, typ)

	c.updateTypeNames(func(t *ast.TypeName) {
		if *t == oldType {
			t.Schema = *stmt.NewSchema
		}
	})
	return nil
}

func (c *Catalog) dropType(stmt *ast.DropTypeStmt) error {
	for _, name := range stmt.Types {
		ns := name.Schema
		if ns == "" {
			ns = c.DefaultSchema
		}
		schema, err := c.getSchema(ns)
		if errors.Is(err, sqlerr.NotFound) && stmt.IfExists {
			continue
		} else if err != nil {
			return err
		}

		_, idx, err := schema.getType(name)
		if errors.Is(err, sqlerr.NotFound) && stmt.IfExists {
			continue
		} else if err != nil {
			return err
		}

		schema.Types = append(schema.Types[:idx], schema.Types[idx+1:]...)
	}
	return nil
}

func (c *Catalog) renameType(stmt *ast.RenameTypeStmt) error {
	if stmt.NewName == nil {
		return fmt.Errorf("rename type: empty name")
	}
	newName := *stmt.NewName
	ns := stmt.Type.Schema
	if ns == "" {
		ns = c.DefaultSchema
	}
	schema, err := c.getSchema(ns)
	if err != nil {
		return err
	}
	ityp, idx, err := schema.getType(stmt.Type)
	if err != nil {
		return err
	}
	if _, _, err := schema.getTable(&ast.TableName{Name: newName}); err == nil {
		return sqlerr.RelationExists(newName)
	}
	if _, _, err := schema.getType(&ast.TypeName{Name: newName}); err == nil {
		return sqlerr.TypeExists(newName)
	}

	switch typ := ityp.(type) {

	case *CompositeType:
		schema.Types[idx] = &CompositeType{
			Name:    newName,
			Columns: typ.Columns,
			Comment: typ.Comment,
		}

	case *Enum:
		schema.Types[idx] = &Enum{
			Name:    newName,
			Vals:    typ.Vals,
			Comment: typ.Comment,
		}

	default:
		return fmt.Errorf("unsupported type: %T", typ)

	}

	c.updateTypeNames(func(t *ast.TypeName) {
		if *t == *stmt.Type {
			t.Name = newName
		}
	})

	return nil
}

func (c *Catalog) updateTypeNames(typeUpdater func(t *ast.TypeName)) error {
	for _, schema := range c.Schemas {
		// Update all the table columns with the new type
		for _, table := range schema.Tables {
			for _, column := range table.Columns {
				typeUpdater(&column.Type)
			}
		}
		// Update all the composite fields with the new type
		for _, typ := range schema.Types {
			composite, ok := typ.(*CompositeType)
			if !ok {
				continue
			}
			for _, fieldType := range composite.Columns {
				typeUpdater(&fieldType.Type)
			}
		}
	}
	return nil
}

package catalog

import (
	"github.com/sqlc-dev/sqlc/internal/sql/ast"
)

// Catalog describes a database instance consisting of metadata in which database objects are defined
type Catalog struct {
	Comment       string
	DefaultSchema string
	Name          string
	Schemas       []*Schema
	SearchPath    []string
	LoadExtension func(string) *Schema

	// TODO: un-export
	Extensions map[string]struct{}
}

// New creates a new catalog
func New(defaultSchema string) *Catalog {

	newCatalog := &Catalog{
		DefaultSchema: defaultSchema,
		Schemas:       make([]*Schema, 0),
		Extensions:    make(map[string]struct{}),
	}

	if newCatalog.DefaultSchema != "" {
		newCatalog.Schemas = append(newCatalog.Schemas, &Schema{Name: defaultSchema})
	}

	return newCatalog
}

// Clone creates a deep copy of the catalog, preserving all schemas, tables, types, and functions.
// This is used to create isolated copies for query parsing where functions might be registered
// with context-dependent types without affecting the original catalog.
func (c *Catalog) Clone() *Catalog {
	if c == nil {
		return nil
	}

	cloned := &Catalog{
		Comment:       c.Comment,
		DefaultSchema: c.DefaultSchema,
		Name:          c.Name,
		Schemas:       make([]*Schema, 0, len(c.Schemas)),
		SearchPath:    make([]string, len(c.SearchPath)),
		LoadExtension: c.LoadExtension,
		Extensions:    make(map[string]struct{}),
	}

	// Copy search path
	copy(cloned.SearchPath, c.SearchPath)

	// Copy extensions
	for k, v := range c.Extensions {
		cloned.Extensions[k] = v
	}

	// Clone schemas
	for _, schema := range c.Schemas {
		cloned.Schemas = append(cloned.Schemas, cloneSchema(schema))
	}

	return cloned
}

func cloneSchema(s *Schema) *Schema {
	if s == nil {
		return nil
	}

	cloned := &Schema{
		Name:    s.Name,
		Comment: s.Comment,
		Tables:  make([]*Table, len(s.Tables)),
		Types:   make([]Type, len(s.Types)),
		Funcs:   make([]*Function, len(s.Funcs)),
	}

	// Clone tables
	for i, table := range s.Tables {
		cloned.Tables[i] = cloneTable(table)
	}

	// Clone types
	for i, t := range s.Types {
		cloned.Types[i] = cloneType(t)
	}

	// Clone functions
	for i, fn := range s.Funcs {
		cloned.Funcs[i] = cloneFunction(fn)
	}

	return cloned
}

func cloneFunction(f *Function) *Function {
	if f == nil {
		return nil
	}

	cloned := &Function{
		Name:               f.Name,
		Comment:            f.Comment,
		Desc:               f.Desc,
		ReturnType:         f.ReturnType, // ast.TypeName is immutable for our purposes
		ReturnTypeNullable: f.ReturnTypeNullable,
		Args:               make([]*Argument, len(f.Args)),
	}

	for i, arg := range f.Args {
		if arg != nil {
			cloned.Args[i] = &Argument{
				Name:       arg.Name,
				Type:       arg.Type, // ast.TypeName is immutable
				HasDefault: arg.HasDefault,
				Mode:       arg.Mode,
			}
		}
	}

	return cloned
}

func cloneTable(t *Table) *Table {
	if t == nil {
		return nil
	}

	cloned := &Table{
		Rel:     t.Rel,
		Comment: t.Comment,
		Columns: make([]*Column, len(t.Columns)),
	}

	for i, col := range t.Columns {
		if col != nil {
			colClone := &Column{
				Name:       col.Name,
				Type:       col.Type,
				IsNotNull:  col.IsNotNull,
				IsUnsigned: col.IsUnsigned,
				IsArray:    col.IsArray,
				ArrayDims:  col.ArrayDims,
				Comment:    col.Comment,
				linkedType: col.linkedType,
			}
			if col.Length != nil {
				length := *col.Length
				colClone.Length = &length
			}
			cloned.Columns[i] = colClone
		}
	}

	return cloned
}

func cloneType(t Type) Type {
	if t == nil {
		return nil
	}

	switch typ := t.(type) {
	case *Enum:
		return &Enum{
			Name:    typ.Name,
			Vals:    append([]string{}, typ.Vals...),
			Comment: typ.Comment,
		}
	case *CompositeType:
		return &CompositeType{
			Name:    typ.Name,
			Comment: typ.Comment,
		}
	default:
		return t
	}
}

func (c *Catalog) Build(stmts []ast.Statement) error {
	for i := range stmts {
		if err := c.Update(stmts[i], nil); err != nil {
			return err
		}
	}
	return nil
}

func (c *Catalog) Update(stmt ast.Statement, colGen columnGenerator) error {
	if stmt.Raw == nil {
		return nil
	}
	var err error
	switch n := stmt.Raw.Stmt.(type) {

	case *ast.AlterTableStmt:
		err = c.alterTable(n)

	case *ast.AlterTableSetSchemaStmt:
		err = c.alterTableSetSchema(n)

	case *ast.AlterTypeAddValueStmt:
		err = c.alterTypeAddValue(n)

	case *ast.AlterTypeRenameValueStmt:
		err = c.alterTypeRenameValue(n)

	case *ast.AlterTypeSetSchemaStmt:
		err = c.alterTypeSetSchema(n)

	case *ast.CommentOnColumnStmt:
		err = c.commentOnColumn(n)

	case *ast.CommentOnSchemaStmt:
		err = c.commentOnSchema(n)

	case *ast.CommentOnTableStmt:
		err = c.commentOnTable(n)

	case *ast.CommentOnTypeStmt:
		err = c.commentOnType(n)

	case *ast.CommentOnViewStmt:
		err = c.commentOnView(n)

	case *ast.CompositeTypeStmt:
		err = c.createCompositeType(n)

	case *ast.CreateEnumStmt:
		err = c.createEnum(n)

	case *ast.CreateExtensionStmt:
		err = c.createExtension(n)

	case *ast.CreateFunctionStmt:
		err = c.createFunction(n)

	case *ast.CreateSchemaStmt:
		err = c.createSchema(n)

	case *ast.CreateTableStmt:
		err = c.createTable(n)

	case *ast.CreateTableAsStmt:
		err = c.createTableAs(n, colGen)

	case *ast.ViewStmt:
		err = c.createView(n, colGen)

	case *ast.DropFunctionStmt:
		err = c.dropFunction(n)

	case *ast.DropSchemaStmt:
		err = c.dropSchema(n)

	case *ast.DropTableStmt:
		err = c.dropTable(n)

	case *ast.DropTypeStmt:
		err = c.dropType(n)

	case *ast.RenameColumnStmt:
		err = c.renameColumn(n)

	case *ast.RenameTableStmt:
		err = c.renameTable(n)

	case *ast.RenameTypeStmt:
		err = c.renameType(n)

	case *ast.List:
		for _, nn := range n.Items {
			if err = c.Update(ast.Statement{
				Raw: &ast.RawStmt{
					Stmt:         nn,
					StmtLocation: stmt.Raw.StmtLocation,
					StmtLen:      stmt.Raw.StmtLen,
				},
			}, colGen); err != nil {
				return err
			}
		}

	}
	return err
}

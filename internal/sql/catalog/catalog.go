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

package catalog

import (
	"github.com/sqlc-dev/sqlc/internal/sql/ast"
	"github.com/sqlc-dev/sqlc/internal/sql/sqlerr"
)

func (c *Catalog) commentOnColumn(stmt *ast.CommentOnColumnStmt) error {
	_, t, err := c.getTable(stmt.Table)
	if err != nil {
		return err
	}
	for i := range t.Columns {
		if t.Columns[i].Name == stmt.Col.Name {
			if stmt.Comment != nil {
				t.Columns[i].Comment = *stmt.Comment
			} else {
				t.Columns[i].Comment = ""
			}
			return nil
		}
	}
	return sqlerr.ColumnNotFound(stmt.Table.Name, stmt.Col.Name)
}

func (c *Catalog) commentOnSchema(stmt *ast.CommentOnSchemaStmt) error {
	s, err := c.getSchema(stmt.Schema.Str)
	if err != nil {
		return err
	}
	if stmt.Comment != nil {
		s.Comment = *stmt.Comment
	} else {
		s.Comment = ""
	}
	return nil
}

func (c *Catalog) commentOnTable(stmt *ast.CommentOnTableStmt) error {
	_, t, err := c.getTable(stmt.Table)
	if err != nil {
		return err
	}
	if stmt.Comment != nil {
		t.Comment = *stmt.Comment
	} else {
		t.Comment = ""
	}
	return nil
}

func (c *Catalog) commentOnType(stmt *ast.CommentOnTypeStmt) error {
	t, _, err := c.getType(stmt.Type)
	if err != nil {
		return err
	}
	if stmt.Comment != nil {
		t.SetComment(*stmt.Comment)
	} else {
		t.SetComment("")
	}
	return nil
}

func (c *Catalog) commentOnView(stmt *ast.CommentOnViewStmt) error {
	_, t, err := c.getTable(stmt.View)
	if err != nil {
		return err
	}
	if stmt.Comment != nil {
		t.Comment = *stmt.Comment
	} else {
		t.Comment = ""
	}
	return nil
}

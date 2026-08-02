package schema

import (
	"fmt"
	"strings"

	"github.com/sqlc-dev/sqlc/internal/core"
	"github.com/sqlc-dev/sqlc/internal/sql/ast"
)

func Apply(cat *core.Catalog, n ast.Node) error {
	switch v := n.(type) {
	case nil:
		return nil
	case *ast.RawStmt:
		return Apply(cat, v.Stmt)
	case *ast.List:
		for _, it := range v.Items {
			if err := Apply(cat, it); err != nil {
				return err
			}
		}
		return nil
	case *ast.CreateTableStmt:
		return applyCreateTable(cat, v)
	case *ast.DropTableStmt:
		return applyDropTable(cat, v)
	}
	return nil
}

func applyCreateTable(cat *core.Catalog, stmt *ast.CreateTableStmt) error {
	if stmt.Name == nil {
		return fmt.Errorf("create table with nil name")
	}
	nsOID, err := resolveOrCreateNamespace(cat, stmt.Name.Schema)
	if err != nil {
		return err
	}
	if _, err := cat.ClassOID(nsOID, stmt.Name.Name); err == nil {
		if stmt.IfNotExists {
			return nil
		}
		return fmt.Errorf("relation %q already exists", stmt.Name.Name)
	}
	classOID, err := cat.CreateClass(nsOID, stmt.Name.Name, "r")
	if err != nil {
		return err
	}
	for i, col := range stmt.Cols {
		if col == nil || col.TypeName == nil {
			return fmt.Errorf("column %d on %q: missing type", i+1, stmt.Name.Name)
		}
		typeOID, err := resolveOrCreateType(cat, col.TypeName.Name)
		if err != nil {
			return fmt.Errorf("column %s.%s: %w", stmt.Name.Name, col.Colname, err)
		}
		if err := cat.CreateAttributeSpec(core.AttributeSpec{
			ClassOID:     classOID,
			Name:         col.Colname,
			TypeOID:      typeOID,
			Num:          i + 1,
			NotNull:      col.IsNotNull || col.PrimaryKey,
			IsPrimaryKey: col.PrimaryKey,
			DeclType:     col.TypeName.Name,
		}); err != nil {
			return fmt.Errorf("attr %s.%s: %w", stmt.Name.Name, col.Colname, err)
		}
	}
	return nil
}

func applyDropTable(cat *core.Catalog, stmt *ast.DropTableStmt) error {
	for _, tn := range stmt.Tables {
		if tn == nil {
			continue
		}
		nsOID, err := cat.NamespaceOID(nsName(tn.Schema))
		if err != nil {
			if stmt.IfExists {
				continue
			}
			return err
		}
		classOID, err := cat.ClassOID(nsOID, tn.Name)
		if err != nil {
			if stmt.IfExists {
				continue
			}
			return fmt.Errorf("drop table %q: %w", tn.Name, err)
		}
		if err := cat.DropClass(classOID); err != nil {
			return fmt.Errorf("drop table %q: %w", tn.Name, err)
		}
	}
	return nil
}

func nsName(schema string) string {
	if schema == "" {
		return "public"
	}
	return schema
}

func resolveOrCreateNamespace(cat *core.Catalog, schema string) (int64, error) {
	name := nsName(schema)
	if oid, err := cat.NamespaceOID(name); err == nil {
		return oid, nil
	}
	return cat.CreateNamespace(name)
}

func resolveOrCreateType(cat *core.Catalog, name string) (int64, error) {
	canonical := strings.ToLower(strings.TrimSpace(name))
	if oid, err := cat.TypeOID(canonical); err == nil {
		return oid, nil
	}
	return cat.CreateType(canonical, 0)
}

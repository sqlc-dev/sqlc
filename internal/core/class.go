package core

import (
	"context"
	"fmt"

	"github.com/sqlc-dev/sqlc/internal/core/catalogdb"
)

// CreateClass inserts a relation (table, view, etc.) and returns its OID.
// Kind should be 'r' (table), 'v' (view), or 'i' (index).
func (c *Catalog) CreateClass(namespaceOID int64, name string, kind string) (int64, error) {
	oid, err := c.q.CreateClass(context.Background(), catalogdb.CreateClassParams{
		NamespaceOid: namespaceOID,
		Name:         name,
		Kind:         kind,
	})
	if err != nil {
		return 0, fmt.Errorf("create class %q: %w", name, err)
	}
	return oid, nil
}

// ClassOID returns the OID for the given relation in the given namespace.
func (c *Catalog) ClassOID(namespaceOID int64, name string) (int64, error) {
	oid, err := c.q.ClassOID(context.Background(), catalogdb.ClassOIDParams{
		NamespaceOid: namespaceOID,
		Name:         name,
	})
	if err != nil {
		return 0, fmt.Errorf("class %q: %w", name, err)
	}
	return oid, nil
}

// ClassOIDByName looks up a relation by name across all namespaces.
// If multiple matches exist, it returns the first found.
func (c *Catalog) ClassOIDByName(name string) (int64, error) {
	oid, err := c.q.ClassOIDByName(context.Background(), name)
	if err != nil {
		return 0, fmt.Errorf("class %q: %w", name, err)
	}
	return oid, nil
}

// DropClass removes a relation and all of its columns from the catalog.
func (c *Catalog) DropClass(classOID int64) error {
	ctx := context.Background()
	if err := c.q.DeleteAttributesByClass(ctx, classOID); err != nil {
		return fmt.Errorf("drop class %d attributes: %w", classOID, err)
	}
	if err := c.q.DeleteClass(ctx, classOID); err != nil {
		return fmt.Errorf("drop class %d: %w", classOID, err)
	}
	return nil
}

// ClassInfo is a relation row exposed for codegen catalog projection.
type ClassInfo struct {
	OID  int64
	Name string
}

// TablesInNamespace returns the tables (kind 'r') in a namespace, by OID.
func (c *Catalog) TablesInNamespace(namespaceOID int64) ([]ClassInfo, error) {
	rows, err := c.q.ListTablesInNamespace(context.Background(), namespaceOID)
	if err != nil {
		return nil, fmt.Errorf("list tables in namespace %d: %w", namespaceOID, err)
	}
	out := make([]ClassInfo, 0, len(rows))
	for _, r := range rows {
		out = append(out, ClassInfo{OID: r.Oid, Name: r.Name})
	}
	return out, nil
}

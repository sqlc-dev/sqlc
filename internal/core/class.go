package core

import "fmt"

// CreateClass inserts a relation (table, view, etc.) and returns its OID.
// Kind should be 'r' (table), 'v' (view), or 'i' (index).
func (c *Catalog) CreateClass(namespaceOID int64, name string, kind string) (int64, error) {
	res, err := c.db.Exec(
		`INSERT INTO sql_class (namespace_oid, name, kind) VALUES (?, ?, ?)`,
		namespaceOID, name, kind,
	)
	if err != nil {
		return 0, fmt.Errorf("create class %q: %w", name, err)
	}
	return res.LastInsertId()
}

// ClassOID returns the OID for the given relation in the given namespace.
func (c *Catalog) ClassOID(namespaceOID int64, name string) (int64, error) {
	var oid int64
	err := c.db.QueryRow(
		`SELECT oid FROM sql_class WHERE namespace_oid = ? AND name = ?`,
		namespaceOID, name,
	).Scan(&oid)
	if err != nil {
		return 0, fmt.Errorf("class %q: %w", name, err)
	}
	return oid, nil
}

// ClassOIDByName looks up a relation by name across all namespaces.
// If multiple matches exist, it returns the first found.
func (c *Catalog) ClassOIDByName(name string) (int64, error) {
	var oid int64
	err := c.db.QueryRow(
		`SELECT oid FROM sql_class WHERE name = ? LIMIT 1`, name,
	).Scan(&oid)
	if err != nil {
		return 0, fmt.Errorf("class %q: %w", name, err)
	}
	return oid, nil
}

package core

import "fmt"

// CreateNamespace inserts a namespace and returns its OID.
func (c *Catalog) CreateNamespace(name string) (int64, error) {
	res, err := c.db.Exec(`INSERT INTO sql_namespace (name) VALUES (?)`, name)
	if err != nil {
		return 0, fmt.Errorf("create namespace %q: %w", name, err)
	}
	return res.LastInsertId()
}

// NamespaceOID returns the OID for the given namespace name.
func (c *Catalog) NamespaceOID(name string) (int64, error) {
	var oid int64
	err := c.db.QueryRow(`SELECT oid FROM sql_namespace WHERE name = ?`, name).Scan(&oid)
	if err != nil {
		return 0, fmt.Errorf("namespace %q: %w", name, err)
	}
	return oid, nil
}

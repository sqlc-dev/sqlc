package core

import "fmt"

// CreateConstraint inserts a constraint for the given relation.
// Kind should be 'p' (primary key), 'f' (foreign key), 'u' (unique), or 'c' (check).
// columns is a comma-separated list of attribute ordinal positions.
func (c *Catalog) CreateConstraint(classOID int64, name string, kind string, columns string) error {
	_, err := c.db.Exec(
		`INSERT INTO sql_constraint (class_oid, name, kind, columns) VALUES (?, ?, ?, ?)`,
		classOID, name, kind, columns,
	)
	if err != nil {
		return fmt.Errorf("create constraint %q on class %d: %w", name, classOID, err)
	}
	return nil
}

package core

import (
	"context"
	"fmt"

	"github.com/sqlc-dev/sqlc/internal/core/catalogdb"
)

// CreateConstraint inserts a constraint for the given relation.
// Kind should be 'p' (primary key), 'f' (foreign key), 'u' (unique), or 'c' (check).
// columns is a comma-separated list of attribute ordinal positions.
func (c *Catalog) CreateConstraint(classOID int64, name string, kind string, columns string) error {
	err := c.q.CreateConstraint(context.Background(), catalogdb.CreateConstraintParams{
		ClassOid: classOID,
		Name:     name,
		Kind:     kind,
		Columns:  columns,
	})
	if err != nil {
		return fmt.Errorf("create constraint %q on class %d: %w", name, classOID, err)
	}
	return nil
}

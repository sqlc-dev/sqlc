package core

import (
	"context"
	"fmt"

	"github.com/sqlc-dev/sqlc/internal/core/catalogdb"
)

// CreateDialect inserts a SQL dialect and returns its OID.
func (c *Catalog) CreateDialect(name string) (int64, error) {
	oid, err := c.q.CreateDialect(context.Background(), name)
	if err != nil {
		return 0, fmt.Errorf("create dialect %q: %w", name, err)
	}
	return oid, nil
}

// DialectOID returns the OID for a registered dialect.
func (c *Catalog) DialectOID(name string) (int64, error) {
	oid, err := c.q.DialectOID(context.Background(), name)
	if err != nil {
		return 0, fmt.Errorf("dialect %q: %w", name, err)
	}
	return oid, nil
}

// SetDialectFlag stores a per-dialect configuration value.
// If the key already exists, the value is replaced.
func (c *Catalog) SetDialectFlag(dialectOID int64, key, value string) error {
	err := c.q.SetDialectFlag(context.Background(), catalogdb.SetDialectFlagParams{
		DialectOid: dialectOID,
		Key:        key,
		Value:      value,
	})
	if err != nil {
		return fmt.Errorf("set dialect flag %s.%s: %w", key, value, err)
	}
	return nil
}

// DialectFlag returns the value of a dialect flag, or "" if not set.
func (c *Catalog) DialectFlag(dialectOID int64, key string) (string, error) {
	value, err := c.q.DialectFlag(context.Background(), catalogdb.DialectFlagParams{
		DialectOid: dialectOID,
		Key:        key,
	})
	if err != nil {
		return "", nil // missing flag = empty string, not an error
	}
	return value, nil
}

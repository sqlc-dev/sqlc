package core

import (
	"context"
	"fmt"

	"github.com/sqlc-dev/sqlc/internal/core/catalogdb"
)

func (c *Catalog) CreateDialect(name string) (int64, error) {
	oid, err := c.q.CreateDialect(context.Background(), name)
	if err != nil {
		return 0, fmt.Errorf("create dialect %q: %w", name, err)
	}
	return oid, nil
}

func (c *Catalog) DialectOID(name string) (int64, error) {
	oid, err := c.q.DialectOID(context.Background(), name)
	if err != nil {
		return 0, fmt.Errorf("dialect %q: %w", name, err)
	}
	return oid, nil
}

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

func (c *Catalog) DialectFlag(dialectOID int64, key string) (string, error) {
	value, err := c.q.DialectFlag(context.Background(), catalogdb.DialectFlagParams{
		DialectOid: dialectOID,
		Key:        key,
	})
	if err != nil {
		return "", nil
	}
	return value, nil
}

package core

import (
	"context"
	"fmt"
	"slices"
	"strings"

	"github.com/sqlc-dev/sqlc/internal/core/catalogdb"
)

func (c *Catalog) CreateDialect(name string) (int64, error) {
	oid, err := c.q.CreateDialect(context.Background(), name)
	if err != nil {
		return 0, fmt.Errorf("create dialect %q: %w", name, err)
	}
	c.dialectOID = oid
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

// FlagComparisonOperators holds the dialect's comparison operators as a
// comma-separated list, so that a type declared by a schema can be given the
// same ones the seeded types have.
const FlagComparisonOperators = "operators.comparison"

// FlagCastCategories holds the categories whose types are all implicitly
// castable to one another, as the dialect's seed declared them, so that a type
// arriving after the seed — an extension's, say — can join its category.
const FlagCastCategories = "casts.categories"

// Constant kinds. A literal in a query has no declared type, so each dialect
// names the type its literals take on.
const (
	ConstInteger = "integer"
	ConstFloat   = "float"
	ConstString  = "string"
	ConstBool    = "bool"
)

// defaultConstTypes are the PostgreSQL type names, used by any dialect that
// does not name its own.
var defaultConstTypes = map[string]string{
	ConstInteger: "int4",
	ConstFloat:   "numeric",
	ConstString:  "text",
	ConstBool:    "bool",
}

// IsComparisonOperator reports whether the dialect counts name among the
// operators that compare two values and yield a boolean.
func (c *Catalog) IsComparisonOperator(name string) bool {
	if c.dialectOID == 0 {
		return false
	}
	ops, _ := c.DialectFlag(c.dialectOID, FlagComparisonOperators)
	if ops == "" {
		return false
	}
	return slices.Contains(strings.Split(ops, ","), name)
}

// SetConstType records the type a literal of the given kind takes on in this
// dialect.
func (c *Catalog) SetConstType(dialectOID int64, kind, typeName string) error {
	return c.SetDialectFlag(dialectOID, "const."+kind, typeName)
}

// ConstTypeOID returns the type a literal of the given kind takes on.
func (c *Catalog) ConstTypeOID(kind string) (int64, error) {
	if c.dialectOID != 0 {
		if name, _ := c.DialectFlag(c.dialectOID, "const."+kind); name != "" {
			return c.TypeOID(name)
		}
	}
	name, ok := defaultConstTypes[kind]
	if !ok {
		return 0, fmt.Errorf("unknown constant kind %q", kind)
	}
	return c.TypeOID(name)
}

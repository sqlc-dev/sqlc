package core

import "sync"

// A Canonicalizer rewrites a type expression into the form its engine
// stores and reports: ClickHouse turns Decimal32(4) into Decimal(9, 4) and
// Enum('a', 'b') into Enum8('a' = 1, 'b' = 2), SQL Server turns float(24)
// into real. It sees each expression as a whole before its arguments are
// interned, and again on each argument, so it has to be idempotent. Aliases
// and argument defaults are data in the dialect's seed; a canonicalizer is
// for what only code can say.
type Canonicalizer func(*TypeExpr) *TypeExpr

// A UserTypeBase says what a type family the schema declared and the
// dialect did not seed stands on: SQLite gives every declared spelling one
// of five affinities by a rule over its words, so FOO BAR(3) compares as a
// numeric. It returns the base family's name and the category the new type
// takes, or an empty name for a type that stands on nothing.
type UserTypeBase func(name string) (base, category string)

var (
	hooksMu        sync.RWMutex
	canonicalizers = map[string]Canonicalizer{}
	userTypeBases  = map[string]UserTypeBase{}
)

// RegisterUserTypeBase installs the rule a dialect resolves an unseeded
// type family by, under the dialect's name.
func RegisterUserTypeBase(dialect string, fn UserTypeBase) {
	hooksMu.Lock()
	defer hooksMu.Unlock()
	userTypeBases[dialect] = fn
}

// userTypeBase applies the catalog's dialect's rule for an unseeded family.
func (c *Catalog) userTypeBase(name string) (base, category string) {
	dialect := c.dialectName()
	if dialect == "" {
		return "", ""
	}
	hooksMu.RLock()
	fn := userTypeBases[dialect]
	hooksMu.RUnlock()
	if fn == nil {
		return "", ""
	}
	return fn(name)
}

// RegisterCanonicalizer installs the canonicalizer for a dialect, by the
// name its dialect.json records. An engine registers its own at init, so
// that a catalog restored from the cache — which runs no seed — finds it by
// the dialect it was seeded with.
func RegisterCanonicalizer(dialect string, fn Canonicalizer) {
	hooksMu.Lock()
	defer hooksMu.Unlock()
	canonicalizers[dialect] = fn
}

// canonicalize applies the catalog's dialect's canonicalizer, if any.
func (c *Catalog) canonicalize(t *TypeExpr) *TypeExpr {
	name := c.dialectName()
	if name == "" {
		return t
	}
	hooksMu.RLock()
	fn := canonicalizers[name]
	hooksMu.RUnlock()
	if fn == nil {
		return t
	}
	return fn(t)
}

// dialectName is the name of the dialect the catalog was seeded with.
func (c *Catalog) dialectName() string {
	if c.dialectOID == 0 {
		return ""
	}
	c.dialectNameOnce.Do(func() {
		if row, err := c.q.SeededDialect(contextBackground()); err == nil {
			c.dialect = row.Name
		}
	})
	return c.dialect
}

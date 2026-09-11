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

// A ResultArg is one argument of a function call as a result-type rule
// sees it: its type, when known, and its value when it is an integer
// literal, which is what the scale of toDecimal64(x, 4) is.
type ResultArg struct {
	Type *TypeExpr
	Int  *int64
}

// A ResultType says what a function returns when that depends on its
// arguments in a way no seed can spell: ClickHouse's toDecimal64(x, s) is
// Decimal(18, s). It returns nil to leave the answer to the catalog.
type ResultType func(name string, args []ResultArg) *TypeExpr

var (
	hooksMu        sync.RWMutex
	canonicalizers = map[string]Canonicalizer{}
	userTypeBases  = map[string]UserTypeBase{}
	resultTypes    = map[string]ResultType{}
)

// RegisterResultType installs a dialect's result-type rule, under the
// dialect's name.
func RegisterResultType(dialect string, fn ResultType) {
	hooksMu.Lock()
	defer hooksMu.Unlock()
	resultTypes[dialect] = fn
}

// ResultTypeOf applies the catalog's dialect's result-type rule to a call,
// or returns nil when there is none or it has nothing to say.
func (c *Catalog) ResultTypeOf(name string, args []ResultArg) *TypeExpr {
	dialect := c.dialectName()
	if dialect == "" {
		return nil
	}
	hooksMu.RLock()
	fn := resultTypes[dialect]
	hooksMu.RUnlock()
	if fn == nil {
		return nil
	}
	return fn(name, args)
}

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

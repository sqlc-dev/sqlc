package core

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/sqlc-dev/sqlc/internal/core/catalogdb"
	"github.com/sqlc-dev/sqlc/internal/core/catalogdef"

	_ "modernc.org/sqlite"
)

//go:generate go run github.com/sqlc-dev/sqlc/cmd/sqlc generate

type Catalog struct {
	db    *sql.DB
	stmts *stmtCache
	q     *catalogdb.Queries

	// dialectOID is the dialect this catalog was seeded with. A catalog is
	// built for one dialect, so dialect-wide lookups need no other input.
	dialectOID int64

	// deferred holds the parts of a seed that are only worth loading if a
	// query asks for them. They run at most once, when a namespace lookup
	// misses.
	deferred     []func(*Catalog) error
	deferredDone bool
}

// SeedLater registers a part of the seed to run the first time the catalog is
// asked for a namespace it does not have. A dialect's system catalogs run to
// thousands of columns that most queries never reference, so loading them is
// left until one does.
func (c *Catalog) SeedLater(fn func(*Catalog) error) {
	c.deferred = append(c.deferred, fn)
}

// runDeferred runs the deferred seeds, reporting whether it had any to run.
func (c *Catalog) runDeferred() (bool, error) {
	if c.deferredDone || len(c.deferred) == 0 {
		return false, nil
	}
	// Marked done up front: a deferred seed looks namespaces up itself, and
	// must not set itself running again.
	c.deferredDone = true
	for _, fn := range c.deferred {
		if err := fn(c); err != nil {
			return false, err
		}
	}
	return true, nil
}

type Option func(*Catalog) error

func WithSeed(fn func(*Catalog) error) Option {
	return func(c *Catalog) error { return fn(c) }
}

func New(opts ...Option) (*Catalog, error) {
	db, err := sql.Open("sqlite", ":memory:")
	if err != nil {
		return nil, fmt.Errorf("core: open catalog: %w", err)
	}
	// Every ":memory:" connection is its own empty database, so the pool has to
	// be pinned to a single connection that is never retired — otherwise a
	// second connection would see a catalog with no tables in it.
	db.SetMaxOpenConns(1)
	db.SetMaxIdleConns(1)
	db.SetConnMaxIdleTime(0)
	db.SetConnMaxLifetime(0)

	// The catalog is scratch state rebuilt on every run and never read back
	// from disk, so durability buys nothing and costs a journal write per
	// statement.
	for _, pragma := range []string{
		"PRAGMA journal_mode = OFF",
		"PRAGMA synchronous = OFF",
	} {
		if _, err := db.Exec(pragma); err != nil {
			db.Close()
			return nil, fmt.Errorf("core: %s: %w", pragma, err)
		}
	}

	// catalogdef.Schema is a batch of DDL statements, so it goes to the
	// connection directly rather than through the prepared-statement cache.
	if _, err := db.Exec(catalogdef.Schema); err != nil {
		db.Close()
		return nil, fmt.Errorf("core: init schema: %w", err)
	}

	stmts := newStmtCache(db)
	c := &Catalog{db: db, stmts: stmts, q: catalogdb.New(stmts)}

	// Run the seeds inside one transaction rather than paying a commit for each
	// of the hundreds of rows a dialect registers.
	ctx := context.Background()
	if err := stmts.begin(ctx); err != nil {
		db.Close()
		return nil, fmt.Errorf("core: begin setup: %w", err)
	}
	if err := c.setup(opts); err != nil {
		stmts.rollback()
		db.Close()
		return nil, err
	}
	if err := stmts.commit(); err != nil {
		db.Close()
		return nil, fmt.Errorf("core: commit setup: %w", err)
	}
	return c, nil
}

func (c *Catalog) setup(opts []Option) error {
	if err := c.bootstrap(); err != nil {
		return fmt.Errorf("core: bootstrap: %w", err)
	}
	for i, opt := range opts {
		if err := opt(c); err != nil {
			return fmt.Errorf("core: option %d: %w", i, err)
		}
	}
	return nil
}

func (c *Catalog) Close() error {
	err := c.stmts.Close()
	if cerr := c.db.Close(); err == nil {
		err = cerr
	}
	return err
}

func (c *Catalog) DB() *sql.DB {
	return c.db
}

func (c *Catalog) bootstrap() error {
	_, err := c.CreateNamespace("public")
	return err
}

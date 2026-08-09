package core

import (
	"context"
	"database/sql"
	"fmt"
	"runtime"

	"github.com/sqlc-dev/sqlc/internal/core/catalogdb"
	"github.com/sqlc-dev/sqlc/internal/core/catalogdef"

	_ "github.com/ncruces/go-sqlite3/driver"
)

//go:generate go run github.com/sqlc-dev/sqlc/cmd/sqlc generate

type Catalog struct {
	db    *sql.DB
	stmts *stmtCache
	q     *catalogdb.Queries

	// dialectOID is the dialect this catalog was seeded with. A catalog is
	// built for one dialect, so dialect-wide lookups need no other input.
	dialectOID int64

	// loadExtension applies a named extension's seed, installed by the
	// dialect's seed for the engines that ship extension data. It runs at most
	// once per extension name: a schema is free to say CREATE EXTENSION twice.
	loadExtension func(name string) error
	extensions    map[string]bool
}

type Option func(*Catalog) error

func WithSeed(fn func(*Catalog) error) Option {
	return func(c *Catalog) error { return fn(c) }
}

func New(opts ...Option) (*Catalog, error) {
	db, err := openDB()
	if err != nil {
		return nil, err
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

// openFile opens a catalog held in a file, read only. The file is named after
// the hash of its contents, so it can never change under a reader: SQLite is
// told as much, which lets any number of processes share it with no locking
// and no copy.
func openFile(path string) (*sql.DB, error) {
	db, err := sql.Open("sqlite3", "file:"+path+"?mode=ro&immutable=1")
	if err != nil {
		return nil, fmt.Errorf("core: open catalog %s: %w", path, err)
	}
	// An immutable file has no locking to contend for, so queries analyzed
	// concurrently each get their own connection rather than queueing on one.
	db.SetMaxOpenConns(runtime.GOMAXPROCS(0))
	db.SetMaxIdleConns(runtime.GOMAXPROCS(0))
	db.SetConnMaxIdleTime(0)
	db.SetConnMaxLifetime(0)
	return db, nil
}

// pinPool holds a catalog being built to a single connection that is never
// retired: every ":memory:" connection is its own empty database, so a second
// one would see a catalog with no tables in it.
func pinPool(db *sql.DB) {
	db.SetMaxOpenConns(1)
	db.SetMaxIdleConns(1)
	db.SetConnMaxIdleTime(0)
	db.SetConnMaxLifetime(0)
}

// openDB opens the empty SQLite database a catalog is built in.
func openDB() (*sql.DB, error) {
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		return nil, fmt.Errorf("core: open catalog: %w", err)
	}
	pinPool(db)

	// The catalog being built is scratch state that is never read back from
	// disk, so durability buys nothing and costs a journal write per statement.
	// Foreign keys are declared in the schema as documentation but not
	// enforced: dialect seeds carry dangling references (a proc whose return
	// type was never registered), and enforcement would also charge an index
	// lookup per seeded row. The driver compiles SQLite with enforcement on by
	// default, so it is switched off explicitly.
	for _, pragma := range []string{
		"PRAGMA journal_mode = OFF",
		"PRAGMA synchronous = OFF",
		"PRAGMA foreign_keys = OFF",
	} {
		if _, err := db.Exec(pragma); err != nil {
			db.Close()
			return nil, fmt.Errorf("core: %s: %w", pragma, err)
		}
	}
	return db, nil
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

// SeededDialectOID returns the dialect this catalog was seeded with, which is
// what anything registered after seeding — an extension's types and functions
// — records as its dialect.
func (c *Catalog) SeededDialectOID() int64 {
	return c.dialectOID
}

// SetExtensionLoader installs the function CREATE EXTENSION calls to apply an
// extension's seed. A dialect without extension data leaves it unset.
func (c *Catalog) SetExtensionLoader(fn func(name string) error) {
	c.loadExtension = fn
}

// LoadExtension applies the named extension's seed to the catalog. An
// extension already applied, or a dialect with no loader, contributes
// nothing; an extension the loader does not know is the loader's to ignore.
func (c *Catalog) LoadExtension(name string) error {
	if c.loadExtension == nil || c.extensions[name] {
		return nil
	}
	if c.extensions == nil {
		c.extensions = map[string]bool{}
	}
	c.extensions[name] = true
	return c.loadExtension(name)
}

package core

import (
	"context"
	"database/sql"
	"fmt"
	"runtime"

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
	db, err := sql.Open("sqlite", "file:"+path+"?mode=ro&immutable=1")
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
	db, err := sql.Open("sqlite", ":memory:")
	if err != nil {
		return nil, fmt.Errorf("core: open catalog: %w", err)
	}
	pinPool(db)

	// The catalog being built is scratch state that is never read back from
	// disk, so durability buys nothing and costs a journal write per statement.
	for _, pragma := range []string{
		"PRAGMA journal_mode = OFF",
		"PRAGMA synchronous = OFF",
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

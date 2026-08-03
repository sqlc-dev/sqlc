package core

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log/slog"

	"github.com/sqlc-dev/sqlc/internal/cache"
	"github.com/sqlc-dev/sqlc/internal/core/catalogdb"
)

// The name of the sole output blob a CoreCatalog action produces.
const catalogOutput = "catalog.db"

// serializer is the driver's serialization API, which the SQLite driver
// implements on its connections.
type serializer interface {
	Serialize() ([]byte, error)
}

// NewCached returns the catalog a dialect and a schema produce, restored from
// the cache when a run has produced it before.
//
// Producing one means seeding the dialect — thousands of rows, since
// PostgreSQL's functions and system catalogs alone run to five figures — and
// then parsing the schema and applying its DDL. The result is the same every
// time, so the finished database is stored in the cache and opened directly on
// the next run, turning all of that work into opening a file.
// It is a single action: the dialect and the schema together are what the
// catalog is, so they are one key, and a hit means apply is never called.
//
// The dialect data is embedded in the sqlc binary, whose digest is an input to
// every action, so the dialect only has to be named. The schema is hashed in
// the order it is applied, since DDL order decides what the catalog holds.
//
// A catalog whose apply fails is returned but not cached, so that a caller can
// report the failure against the catalog it got. A nil catalog means the
// catalog could not be created at all. A cache that cannot be opened, read or
// written is not an error: the catalog is simply built the long way.
func NewCached(dialect string, schema []string, apply func(*Catalog) error, opts ...Option) (*Catalog, error) {
	store, action := openAction(dialect, schema)
	if store != nil {
		defer store.Close()
		if cat, err := restore(store, action); err == nil {
			return cat, nil
		} else if !errors.Is(err, cache.ErrNotFound) {
			slog.Debug("restoring the catalog from the cache failed", "err", err)
		}
	}

	cat, err := New(opts...)
	if err != nil {
		return nil, err
	}
	if err := apply(cat); err != nil {
		return cat, err
	}
	if store != nil {
		if err := save(store, action, cat); err != nil {
			slog.Debug("saving the catalog to the cache failed", "err", err)
		}
	}
	return cat, nil
}

// openAction opens the cache and digests the action that produces a catalog. A
// nil store means the cache is unavailable and the catalog has to be built.
func openAction(dialect string, schema []string) (*cache.Cache, cache.Digest) {
	store, err := cache.Open()
	if err != nil {
		slog.Debug("opening the cache failed", "err", err)
		return nil, cache.Digest{}
	}
	action := store.NewAction("CoreCatalog").AddInput("dialect", []byte(dialect))
	for _, s := range schema {
		action.AddInput("schema", []byte(s))
	}
	return store, action.Digest()
}

// restore opens the catalog a previous run cached. The stored blob is opened
// where it lies rather than copied: analysis only ever reads the catalog, so
// every run that wants this one reads the same file.
func restore(store *cache.Cache, action cache.Digest) (*Catalog, error) {
	result, err := store.Actions.Get(action)
	if err != nil {
		return nil, err
	}
	path, ok := store.CAS.Filename(result.Outputs[catalogOutput])
	if !ok {
		return nil, cache.ErrNotFound
	}

	db, err := openFile(path)
	if err != nil {
		return nil, err
	}

	stmts := newStmtCache(db)
	cat := &Catalog{db: db, stmts: stmts, q: catalogdb.New(stmts)}

	// The seeds recorded the dialect on the catalog as they ran. A restored
	// catalog reads it back from the database instead.
	row, err := cat.q.SeededDialect(context.Background())
	if err != nil {
		cat.Close()
		return nil, fmt.Errorf("core: restored catalog has no dialect: %w", err)
	}
	cat.dialectOID = row.Oid
	return cat, nil
}

// save stores a finished catalog for the next run.
func save(store *cache.Cache, action cache.Digest, cat *Catalog) error {
	blob, err := serialize(cat.db)
	if err != nil {
		return err
	}
	digest, err := store.CAS.Put(blob)
	if err != nil {
		return err
	}
	return store.Actions.Put(action, &cache.ActionResult{
		Outputs: map[string]cache.Digest{catalogOutput: digest},
	})
}

// serialize returns the bytes a database being built in memory would be
// written to disk as, which is what gets stored.
func serialize(db *sql.DB) ([]byte, error) {
	var blob []byte
	err := withRawConn(db, func(s serializer) error {
		var err error
		blob, err = s.Serialize()
		return err
	})
	if err != nil {
		return nil, fmt.Errorf("core: serialize catalog: %w", err)
	}
	return blob, nil
}

func withRawConn(db *sql.DB, fn func(serializer) error) error {
	conn, err := db.Conn(context.Background())
	if err != nil {
		return err
	}
	defer conn.Close()
	return conn.Raw(func(driverConn any) error {
		s, ok := driverConn.(serializer)
		if !ok {
			return fmt.Errorf("driver does not serialize databases")
		}
		return fn(s)
	})
}

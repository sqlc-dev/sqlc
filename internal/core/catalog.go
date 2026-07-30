package core

import (
	"database/sql"
	"fmt"

	"github.com/sqlc-dev/sqlc/internal/core/catalogdb"
	"github.com/sqlc-dev/sqlc/internal/core/catalogdef"

	_ "modernc.org/sqlite"
)

//go:generate go run github.com/sqlc-dev/sqlc/cmd/sqlc generate

type Catalog struct {
	db *sql.DB
	q  *catalogdb.Queries
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
	if _, err := db.Exec(catalogdef.Schema); err != nil {
		db.Close()
		return nil, fmt.Errorf("core: init schema: %w", err)
	}
	c := &Catalog{db: db, q: catalogdb.New(db)}
	if err := c.bootstrap(); err != nil {
		db.Close()
		return nil, fmt.Errorf("core: bootstrap: %w", err)
	}
	for i, opt := range opts {
		if err := opt(c); err != nil {
			db.Close()
			return nil, fmt.Errorf("core: option %d: %w", i, err)
		}
	}
	return c, nil
}

func (c *Catalog) Close() error {
	return c.db.Close()
}

func (c *Catalog) DB() *sql.DB {
	return c.db
}

func (c *Catalog) bootstrap() error {
	_, err := c.CreateNamespace("public")
	return err
}

package core

import (
	"context"
	"database/sql"
	"sync"
)

// stmtCache is a catalogdb.DBTX that keeps one prepared statement per distinct
// query string.
//
// The catalog runs the same dozen or so statements over and over — seeding a
// dialect is hundreds of inserts through a handful of queries, and analysis
// issues one lookup per table, column, type, operator and function a statement
// touches. Handing the raw *sql.DB to catalogdb means SQLite re-parses and
// re-plans that SQL text on every call, which dominated both paths. Preparing
// once and rebinding is the whole win.
//
// Statements live for the lifetime of the catalog and are released by Close.
type stmtCache struct {
	db *sql.DB

	mu    sync.RWMutex
	stmts map[string]*sql.Stmt

	// While a transaction is open every statement must be prepared on the
	// transaction's own connection: the catalog pins the pool to a single
	// connection, so reaching for the *sql.DB mid-transaction would block
	// forever waiting for a connection the transaction already holds. These
	// statements are owned by the transaction and die with it.
	tx      *sql.Tx
	txStmts map[string]*sql.Stmt
}

func newStmtCache(db *sql.DB) *stmtCache {
	return &stmtCache{db: db, stmts: make(map[string]*sql.Stmt)}
}

func (c *stmtCache) prepared(ctx context.Context, query string) (*sql.Stmt, error) {
	c.mu.RLock()
	tx, cache := c.tx, c.stmts
	if tx != nil {
		cache = c.txStmts
	}
	stmt, ok := cache[query]
	c.mu.RUnlock()
	if ok {
		return stmt, nil
	}

	if tx != nil {
		stmt, err := tx.PrepareContext(ctx, query)
		if err != nil {
			return nil, err
		}
		c.mu.Lock()
		defer c.mu.Unlock()
		// Bail out if the transaction ended while we were preparing; the
		// statement is dead either way.
		if c.tx != tx {
			return stmt, nil
		}
		if existing, ok := c.txStmts[query]; ok {
			stmt.Close()
			return existing, nil
		}
		c.txStmts[query] = stmt
		return stmt, nil
	}

	stmt, err := c.db.PrepareContext(ctx, query)
	if err != nil {
		return nil, err
	}

	c.mu.Lock()
	defer c.mu.Unlock()
	// Another goroutine may have prepared the same query while we were not
	// holding the lock; keep the winner so each query has a single statement.
	if existing, ok := c.stmts[query]; ok {
		stmt.Close()
		return existing, nil
	}
	c.stmts[query] = stmt
	return stmt, nil
}

func (c *stmtCache) ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error) {
	stmt, err := c.prepared(ctx, query)
	if err != nil {
		return nil, err
	}
	return stmt.ExecContext(ctx, args...)
}

func (c *stmtCache) QueryContext(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error) {
	stmt, err := c.prepared(ctx, query)
	if err != nil {
		return nil, err
	}
	return stmt.QueryContext(ctx, args...)
}

func (c *stmtCache) QueryRowContext(ctx context.Context, query string, args ...interface{}) *sql.Row {
	stmt, err := c.prepared(ctx, query)
	if err != nil {
		// *sql.Row carries an error but has no exported constructor, so fall
		// back to the unprepared path and let it report the failure on Scan.
		return c.db.QueryRowContext(ctx, query, args...)
	}
	return stmt.QueryRowContext(ctx, args...)
}

// PrepareContext returns a statement owned by the caller; it is deliberately
// not cached, since the caller is responsible for closing it.
func (c *stmtCache) PrepareContext(ctx context.Context, query string) (*sql.Stmt, error) {
	return c.db.PrepareContext(ctx, query)
}

// begin opens a transaction that every subsequent statement runs inside, until
// commit or rollback. Batching writes this way saves SQLite a commit per
// statement, which is most of the cost of installing a dialect seed.
func (c *stmtCache) begin(ctx context.Context) error {
	tx, err := c.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	c.mu.Lock()
	c.tx = tx
	c.txStmts = make(map[string]*sql.Stmt)
	c.mu.Unlock()
	return nil
}

// endTx detaches the current transaction and returns it, or nil if none is
// open. The transaction's statements are released along with it.
func (c *stmtCache) endTx() *sql.Tx {
	c.mu.Lock()
	defer c.mu.Unlock()
	tx := c.tx
	c.tx, c.txStmts = nil, nil
	return tx
}

func (c *stmtCache) commit() error {
	if tx := c.endTx(); tx != nil {
		return tx.Commit()
	}
	return nil
}

func (c *stmtCache) rollback() {
	if tx := c.endTx(); tx != nil {
		tx.Rollback()
	}
}

func (c *stmtCache) Close() error {
	c.rollback()

	c.mu.Lock()
	defer c.mu.Unlock()
	var firstErr error
	for query, stmt := range c.stmts {
		if err := stmt.Close(); err != nil && firstErr == nil {
			firstErr = err
		}
		delete(c.stmts, query)
	}
	return firstErr
}

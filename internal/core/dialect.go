package core

import "fmt"

// CreateDialect inserts a SQL dialect and returns its OID.
func (c *Catalog) CreateDialect(name string) (int64, error) {
	res, err := c.db.Exec(`INSERT INTO sql_dialect (name) VALUES (?)`, name)
	if err != nil {
		return 0, fmt.Errorf("create dialect %q: %w", name, err)
	}
	return res.LastInsertId()
}

// DialectOID returns the OID for a registered dialect.
func (c *Catalog) DialectOID(name string) (int64, error) {
	var oid int64
	err := c.db.QueryRow(`SELECT oid FROM sql_dialect WHERE name = ?`, name).Scan(&oid)
	if err != nil {
		return 0, fmt.Errorf("dialect %q: %w", name, err)
	}
	return oid, nil
}

// SetDialectFlag stores a per-dialect configuration value.
// If the key already exists, the value is replaced.
func (c *Catalog) SetDialectFlag(dialectOID int64, key, value string) error {
	_, err := c.db.Exec(
		`INSERT INTO sql_dialect_flag (dialect_oid, key, value) VALUES (?, ?, ?)
		 ON CONFLICT(dialect_oid, key) DO UPDATE SET value = excluded.value`,
		dialectOID, key, value,
	)
	if err != nil {
		return fmt.Errorf("set dialect flag %s.%s: %w", key, value, err)
	}
	return nil
}

// DialectFlag returns the value of a dialect flag, or "" if not set.
func (c *Catalog) DialectFlag(dialectOID int64, key string) (string, error) {
	var value string
	err := c.db.QueryRow(
		`SELECT value FROM sql_dialect_flag WHERE dialect_oid = ? AND key = ?`,
		dialectOID, key,
	).Scan(&value)
	if err != nil {
		return "", nil // missing flag = empty string, not an error
	}
	return value, nil
}

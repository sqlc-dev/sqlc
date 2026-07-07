package core

import (
	"database/sql"
	"fmt"
	"strings"
)

// TypeSpec describes a SQL type for insertion into the catalog.
// All fields are optional; Name is required.
type TypeSpec struct {
	Name         string
	Size         int
	Typtype      string // 'b'ase | 'c'omposite | 'd'omain | 'e'num | 'p'seudo | 'r'ange
	Category     string // 'N'umeric | 'S'tring | 'B'oolean | 'D'atetime | 'A'rray | ...
	Preferred    bool
	NamespaceOID int64 // 0 = NULL
	DialectOID   int64 // 0 = NULL (shared)
	ElementOID   int64 // for arrays; 0 = NULL
}

// CreateType inserts a base type (typtype='b') and returns its OID.
// For full control, use CreateTypeSpec.
func (c *Catalog) CreateType(name string, size int) (int64, error) {
	return c.CreateTypeSpec(TypeSpec{Name: name, Size: size, Typtype: "b"})
}

// CreateTypeSpec inserts a type with full metadata. NamespaceOID defaults
// to the "public" namespace when zero so the simpler CreateType API and
// engines that don't track namespaces (sqlite) continue to work.
func (c *Catalog) CreateTypeSpec(t TypeSpec) (int64, error) {
	if t.Typtype == "" {
		t.Typtype = "b"
	}
	if t.NamespaceOID == 0 {
		oid, err := c.NamespaceOID("public")
		if err != nil {
			return 0, fmt.Errorf("create type %q: default namespace: %w", t.Name, err)
		}
		t.NamespaceOID = oid
	}
	res, err := c.db.Exec(
		`INSERT INTO sql_type
		   (name, size, typtype, category, preferred, namespace_oid, dialect_oid, element_oid)
		 VALUES (?, ?, ?, ?, ?, ?, ?, ?)`,
		strings.ToLower(t.Name), t.Size, t.Typtype,
		nullableString(t.Category), boolToInt(t.Preferred),
		t.NamespaceOID, nullableOID(t.DialectOID), nullableOID(t.ElementOID),
	)
	if err != nil {
		return 0, fmt.Errorf("create type %q: %w", t.Name, err)
	}
	return res.LastInsertId()
}

// TypeOID returns the OID for the given (unqualified) type name. When
// the same name lives in multiple namespaces (e.g. a system-internal
// duplicate in information_schema and pg_catalog), the lookup prefers
// pg_catalog, then public, then any other namespace by name.
func (c *Catalog) TypeOID(name string) (int64, error) {
	var oid int64
	err := c.db.QueryRow(
		`SELECT t.oid FROM sql_type t
		 JOIN sql_namespace ns ON ns.oid = t.namespace_oid
		 WHERE t.name = ?
		 ORDER BY CASE ns.name
		   WHEN 'pg_catalog' THEN 0
		   WHEN 'public'     THEN 1
		   ELSE 2 END,
		   ns.name
		 LIMIT 1`,
		strings.ToLower(name),
	).Scan(&oid)
	if err != nil {
		return 0, fmt.Errorf("type %q: %w", name, err)
	}
	return oid, nil
}

// TypeName returns the name for the given type OID.
func (c *Catalog) TypeName(oid int64) (string, error) {
	var name string
	err := c.db.QueryRow(`SELECT name FROM sql_type WHERE oid = ?`, oid).Scan(&name)
	if err != nil {
		return "", fmt.Errorf("type oid %d: %w", oid, err)
	}
	return name, nil
}

// TypeInfo bundles the metadata returned by LookupType.
type TypeInfo struct {
	OID       int64
	Name      string
	Category  string
	Typtype   string
	Preferred bool
}

// LookupType returns metadata for the given type OID.
func (c *Catalog) LookupType(oid int64) (TypeInfo, error) {
	var (
		ti       TypeInfo
		category sql.NullString
		pref     int
	)
	err := c.db.QueryRow(
		`SELECT oid, name, COALESCE(category, ''), typtype, preferred FROM sql_type WHERE oid = ?`,
		oid,
	).Scan(&ti.OID, &ti.Name, &category, &ti.Typtype, &pref)
	if err != nil {
		return ti, fmt.Errorf("lookup type oid %d: %w", oid, err)
	}
	ti.Category = category.String
	ti.Preferred = pref != 0
	return ti, nil
}

func nullableOID(oid int64) any {
	if oid == 0 {
		return nil
	}
	return oid
}

func nullableString(s string) any {
	if s == "" {
		return nil
	}
	return s
}

func boolToInt(b bool) int {
	if b {
		return 1
	}
	return 0
}

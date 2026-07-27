package core

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"github.com/sqlc-dev/sqlc/internal/core/catalogdb"
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
	oid, err := c.q.CreateType(context.Background(), catalogdb.CreateTypeParams{
		Name:         strings.ToLower(t.Name),
		Size:         int64(t.Size),
		Typtype:      t.Typtype,
		Category:     nullString(t.Category),
		Preferred:    boolToInt64(t.Preferred),
		NamespaceOid: t.NamespaceOID,
		DialectOid:   nullInt64(t.DialectOID),
		ElementOid:   nullInt64(t.ElementOID),
	})
	if err != nil {
		return 0, fmt.Errorf("create type %q: %w", t.Name, err)
	}
	return oid, nil
}

// TypeOID returns the OID for the given (unqualified) type name. When
// the same name lives in multiple namespaces (e.g. a system-internal
// duplicate in information_schema and pg_catalog), the lookup prefers
// pg_catalog, then public, then any other namespace by name.
func (c *Catalog) TypeOID(name string) (int64, error) {
	oid, err := c.q.TypeOIDByName(context.Background(), strings.ToLower(name))
	if err != nil {
		return 0, fmt.Errorf("type %q: %w", name, err)
	}
	return oid, nil
}

// TypeName returns the name for the given type OID.
func (c *Catalog) TypeName(oid int64) (string, error) {
	name, err := c.q.TypeNameByOID(context.Background(), oid)
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
	row, err := c.q.LookupType(context.Background(), oid)
	if err != nil {
		return TypeInfo{}, fmt.Errorf("lookup type oid %d: %w", oid, err)
	}
	return TypeInfo{
		OID:       row.Oid,
		Name:      row.Name,
		Category:  row.Category.String,
		Typtype:   row.Typtype,
		Preferred: row.Preferred != 0,
	}, nil
}

// nullableOID / nullableString / boolToInt build the loosely-typed
// arguments the not-yet-migrated raw-SQL call sites pass to database/sql.
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

// nullInt64 / nullString / boolToInt64 build the sql.Null* and int64
// arguments the sqlc-generated catalogdb queries expect.
func nullInt64(oid int64) sql.NullInt64 {
	return sql.NullInt64{Int64: oid, Valid: oid != 0}
}

func nullString(s string) sql.NullString {
	return sql.NullString{String: s, Valid: s != ""}
}

func boolToInt64(b bool) int64 {
	if b {
		return 1
	}
	return 0
}

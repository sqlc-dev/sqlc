package core

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"github.com/sqlc-dev/sqlc/internal/core/catalogdb"
)

type TypeSpec struct {
	Name         string
	Size         int
	Typtype      string
	Category     string
	Preferred    bool
	NamespaceOID int64
	DialectOID   int64
	ElementOID   int64
}

func (c *Catalog) CreateType(name string, size int) (int64, error) {
	return c.CreateTypeSpec(TypeSpec{Name: name, Size: size, Typtype: "b"})
}

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

func (c *Catalog) TypeOID(name string) (int64, error) {
	oid, err := c.q.TypeOIDByName(context.Background(), strings.ToLower(name))
	if err != nil {
		return 0, fmt.Errorf("type %q: %w", name, err)
	}
	return oid, nil
}

func (c *Catalog) TypeName(oid int64) (string, error) {
	name, err := c.q.TypeNameByOID(context.Background(), oid)
	if err != nil {
		return "", fmt.Errorf("type oid %d: %w", oid, err)
	}
	return name, nil
}

type TypeInfo struct {
	OID       int64
	Name      string
	Category  string
	Typtype   string
	Preferred bool
}

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

func orZero(n sql.NullInt64) int64 {
	if n.Valid {
		return n.Int64
	}
	return 0
}

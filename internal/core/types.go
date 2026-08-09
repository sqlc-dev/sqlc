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

// ArraySuffix marks the type of an array of the type it is appended to. The
// catalog names an array type after its element, the way a schema spells it.
const ArraySuffix = "[]"

// CreateUserType registers a type a schema declared rather than the dialect,
// such as an enum or a name the dialect's seed does not list. The type gains
// the dialect's comparison operators, so a column of it can be compared.
func (c *Catalog) CreateUserType(name, category string) (int64, error) {
	typtype := "b"
	if category == "E" {
		typtype = "e"
	}
	oid, err := c.CreateTypeSpec(TypeSpec{
		Name:       name,
		Typtype:    typtype,
		Category:   category,
		DialectOID: c.dialectOID,
	})
	if err != nil {
		return 0, err
	}
	if err := c.createComparisons(oid); err != nil {
		return 0, err
	}
	return oid, nil
}

// CreateArrayType registers the array type over elementOID.
func (c *Catalog) CreateArrayType(name string, elementOID int64) (int64, error) {
	oid, err := c.CreateTypeSpec(TypeSpec{
		Name:       name,
		Typtype:    "b",
		Category:   "A",
		DialectOID: c.dialectOID,
		ElementOID: elementOID,
	})
	if err != nil {
		return 0, err
	}
	if err := c.createComparisons(oid); err != nil {
		return 0, err
	}
	return oid, nil
}

// createComparisons gives a type the dialect's comparison operators, which the
// seed registered for the types it knew about up front.
func (c *Catalog) createComparisons(typeOID int64) error {
	if c.dialectOID == 0 {
		return nil
	}
	ops, _ := c.DialectFlag(c.dialectOID, FlagComparisonOperators)
	boolName, _ := c.DialectFlag(c.dialectOID, "const."+ConstBool)
	if ops == "" || boolName == "" {
		return nil
	}
	boolOID, err := c.TypeOID(boolName)
	if err != nil {
		return nil
	}
	for _, op := range strings.Split(ops, ",") {
		if op == "" {
			continue
		}
		if _, err := c.CreateOperator(OperatorSpec{
			Name:          op,
			DialectOID:    c.dialectOID,
			LeftTypeOID:   typeOID,
			RightTypeOID:  typeOID,
			ResultTypeOID: boolOID,
		}); err != nil {
			return err
		}
	}
	return nil
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

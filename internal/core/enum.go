package core

import (
	"context"
	"fmt"
	"strings"

	"github.com/sqlc-dev/sqlc/internal/core/catalogdb"
)

// EnumInfo describes an enum type the schema declared.
type EnumInfo struct {
	OID int64
	// Name is the name the catalog stores, which carries the schema when the
	// type was declared with one ("foo.mood").
	Name string
}

// CreateEnumType registers an enum type with its labels in declaration order.
func (c *Catalog) CreateEnumType(name string, labels []string) (int64, error) {
	oid, err := c.CreateUserType(name, "E")
	if err != nil {
		return 0, err
	}
	if err := c.SetEnumLabels(oid, labels); err != nil {
		return 0, err
	}
	return oid, nil
}

// EnumLabels returns an enum type's labels in declaration order.
func (c *Catalog) EnumLabels(oid int64) ([]string, error) {
	labels, err := c.q.ListEnumLabels(context.Background(), oid)
	if err != nil {
		return nil, fmt.Errorf("enum labels of type oid %d: %w", oid, err)
	}
	return labels, nil
}

// SetEnumLabels replaces an enum type's labels.
func (c *Catalog) SetEnumLabels(oid int64, labels []string) error {
	ctx := context.Background()
	if err := c.q.DeleteEnumLabels(ctx, oid); err != nil {
		return fmt.Errorf("enum labels of type oid %d: %w", oid, err)
	}
	for i, label := range labels {
		if err := c.q.CreateEnumLabel(ctx, catalogdb.CreateEnumLabelParams{TypeOid: oid, Ord: int64(i + 1), Label: label}); err != nil {
			return fmt.Errorf("enum label %q of type oid %d: %w", label, oid, err)
		}
	}
	return nil
}

// Enums lists the enum types the schema declared, in declaration order.
func (c *Catalog) Enums() ([]EnumInfo, error) {
	rows, err := c.q.ListEnumTypes(context.Background())
	if err != nil {
		return nil, fmt.Errorf("list enum types: %w", err)
	}
	out := make([]EnumInfo, 0, len(rows))
	for _, r := range rows {
		out = append(out, EnumInfo{OID: r.Oid, Name: r.Name})
	}
	return out, nil
}

// RenameType gives a type a new name. Columns refer to the type by OID, so
// they follow the rename.
func (c *Catalog) RenameType(oid int64, name string) error {
	if err := c.q.RenameType(context.Background(), catalogdb.RenameTypeParams{Oid: oid, Name: strings.ToLower(name)}); err != nil {
		return fmt.Errorf("rename type oid %d to %q: %w", oid, name, err)
	}
	return nil
}

// DropType removes a type and, when it is an enum, its labels.
func (c *Catalog) DropType(oid int64) error {
	ctx := context.Background()
	if err := c.q.DeleteEnumLabels(ctx, oid); err != nil {
		return fmt.Errorf("drop type oid %d: %w", oid, err)
	}
	if err := c.q.DeleteType(ctx, oid); err != nil {
		return fmt.Errorf("drop type oid %d: %w", oid, err)
	}
	return nil
}

// SplitTypeName separates the schema a type name carries from the name
// itself: "foo.mood" is the type mood in schema foo, and a bare name is in
// the default schema.
func SplitTypeName(name string) (schema, typ string) {
	if i := strings.LastIndex(name, "."); i >= 0 {
		return name[:i], name[i+1:]
	}
	return "", name
}

// JoinTypeName is the inverse of SplitTypeName: the name the catalog stores
// for a type in a schema. The default schema is not spelled.
func JoinTypeName(schema, typ string) string {
	if schema == "" || schema == "public" {
		return typ
	}
	return schema + "." + typ
}

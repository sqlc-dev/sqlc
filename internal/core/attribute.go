package core

import (
	"context"
	"fmt"

	"github.com/sqlc-dev/sqlc/internal/core/catalogdb"
)

type AttributeSpec struct {
	ClassOID      int64
	Name          string
	TypeOID       int64
	Num           int
	NotNull       bool
	HasDefault    bool
	DeclType      string
	TypeLength    int
	TypeScale     int
	AutoIncrement bool
	IsPrimaryKey  bool
	IsUnique      bool
}

func (c *Catalog) CreateAttributeSpec(s AttributeSpec) error {
	err := c.q.CreateAttribute(context.Background(), catalogdb.CreateAttributeParams{
		ClassOid:      s.ClassOID,
		Name:          s.Name,
		TypeOid:       s.TypeOID,
		NotNull:       boolToInt64(s.NotNull),
		HasDefault:    boolToInt64(s.HasDefault),
		Num:           int64(s.Num),
		DeclType:      s.DeclType,
		TypeLength:    int64(s.TypeLength),
		TypeScale:     int64(s.TypeScale),
		AutoIncrement: boolToInt64(s.AutoIncrement),
		IsPrimaryKey:  boolToInt64(s.IsPrimaryKey),
		IsUnique:      boolToInt64(s.IsUnique),
	})
	if err != nil {
		return fmt.Errorf("create attribute %q on class %d: %w", s.Name, s.ClassOID, err)
	}
	return nil
}

func (c *Catalog) CreateAttribute(classOID int64, name string, typeOID int64, notNull bool, hasDefault bool, num int) error {
	return c.CreateAttributeSpec(AttributeSpec{
		ClassOID:   classOID,
		Name:       name,
		TypeOID:    typeOID,
		Num:        num,
		NotNull:    notNull,
		HasDefault: hasDefault,
	})
}

func (c *Catalog) SetAttributePrimaryKey(classOID int64, columns []string) error {
	ctx := context.Background()
	for _, name := range columns {
		err := c.q.SetAttributePrimaryKey(ctx, catalogdb.SetAttributePrimaryKeyParams{
			ClassOid: classOID,
			Name:     name,
		})
		if err != nil {
			return fmt.Errorf("mark pk %s on class %d: %w", name, classOID, err)
		}
	}
	return nil
}

func (c *Catalog) SetAttributeUnique(classOID int64, columns []string) error {
	ctx := context.Background()
	for _, name := range columns {
		err := c.q.SetAttributeUnique(ctx, catalogdb.SetAttributeUniqueParams{
			ClassOid: classOID,
			Name:     name,
		})
		if err != nil {
			return fmt.Errorf("mark unique %s on class %d: %w", name, classOID, err)
		}
	}
	return nil
}

type ColumnInfo struct {
	Name          string
	TypeName      string
	TypeOID       int64
	NotNull       bool
	DeclType      string
	TypeLength    int
	TypeScale     int
	AutoIncrement bool
	IsPrimaryKey  bool
	IsUnique      bool
	AttributeOID  int64
	ClassOID      int64
}

func (c *Catalog) ResolveColumn(table, column string) (*ColumnInfo, error) {
	r, err := c.q.ResolveColumn(context.Background(), catalogdb.ResolveColumnParams{
		TableName:  table,
		ColumnName: column,
	})
	if err != nil {
		return nil, fmt.Errorf("resolve %s.%s: %w", table, column, err)
	}
	info := ColumnInfo{
		AttributeOID:  r.Oid,
		ClassOID:      r.ClassOid,
		Name:          r.Name,
		TypeName:      r.TypeName,
		TypeOID:       r.TypeOid,
		NotNull:       r.NotNull != 0,
		DeclType:      r.DeclType,
		TypeLength:    int(r.TypeLength),
		TypeScale:     int(r.TypeScale),
		AutoIncrement: r.AutoIncrement != 0,
		IsPrimaryKey:  r.IsPrimaryKey != 0,
		IsUnique:      r.IsUnique != 0,
	}
	return &info, nil
}

func (c *Catalog) TableColumns(table string) ([]ColumnInfo, error) {
	rows, err := c.q.TableColumns(context.Background(), table)
	if err != nil {
		return nil, fmt.Errorf("table columns %q: %w", table, err)
	}
	cols := make([]ColumnInfo, 0, len(rows))
	for _, r := range rows {
		cols = append(cols, ColumnInfo{
			AttributeOID:  r.Oid,
			ClassOID:      r.ClassOid,
			Name:          r.Name,
			TypeName:      r.TypeName,
			TypeOID:       r.TypeOid,
			NotNull:       r.NotNull != 0,
			DeclType:      r.DeclType,
			TypeLength:    int(r.TypeLength),
			TypeScale:     int(r.TypeScale),
			AutoIncrement: r.AutoIncrement != 0,
			IsPrimaryKey:  r.IsPrimaryKey != 0,
			IsUnique:      r.IsUnique != 0,
		})
	}
	return cols, nil
}

type ClassColumn struct {
	AttOID  int64
	Name    string
	TypeOID int64
	NotNull bool
}

// ClassColumns returns a relation's columns in ordinal order.
func (c *Catalog) ClassColumns(classOID int64) ([]ClassColumn, error) {
	rows, err := c.q.ClassAttributes(context.Background(), classOID)
	if err != nil {
		return nil, fmt.Errorf("class columns %d: %w", classOID, err)
	}
	out := make([]ClassColumn, 0, len(rows))
	for _, r := range rows {
		out = append(out, ClassColumn{
			AttOID:  r.Oid,
			Name:    r.Name,
			TypeOID: r.TypeOid,
			NotNull: r.NotNull != 0,
		})
	}
	return out, nil
}

type CodegenColumn struct {
	Name     string
	TypeName string
	NotNull  bool
}

func (c *Catalog) ClassCodegenColumns(classOID int64) ([]CodegenColumn, error) {
	rows, err := c.q.ListClassColumns(context.Background(), classOID)
	if err != nil {
		return nil, fmt.Errorf("class codegen columns %d: %w", classOID, err)
	}
	out := make([]CodegenColumn, 0, len(rows))
	for _, r := range rows {
		out = append(out, CodegenColumn{
			Name:     r.ColumnName,
			TypeName: r.TypeName,
			NotNull:  r.NotNull != 0,
		})
	}
	return out, nil
}

type AttributeDetails struct {
	Schema        string
	Table         string
	Column        string
	Num           int
	DeclType      string
	TypeLength    int
	TypeScale     int
	AutoIncrement bool
	IsPrimaryKey  bool
	IsUnique      bool
	NotNull       bool
}

func (c *Catalog) LookupAttribute(attOID int64) (AttributeDetails, error) {
	r, err := c.q.LookupAttribute(context.Background(), attOID)
	if err != nil {
		return AttributeDetails{}, fmt.Errorf("lookup attribute %d: %w", attOID, err)
	}
	return AttributeDetails{
		Schema:        r.SchemaName,
		Table:         r.TableName,
		Column:        r.ColumnName,
		Num:           int(r.Num),
		DeclType:      r.DeclType,
		TypeLength:    int(r.TypeLength),
		TypeScale:     int(r.TypeScale),
		AutoIncrement: r.AutoIncrement != 0,
		IsPrimaryKey:  r.IsPrimaryKey != 0,
		IsUnique:      r.IsUnique != 0,
		NotNull:       r.NotNull != 0,
	}, nil
}

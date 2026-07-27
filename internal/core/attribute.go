package core

import (
	"context"
	"fmt"

	"github.com/sqlc-dev/sqlc/internal/core/catalogdb"
)

// AttributeSpec carries the data needed to register a column on a
// relation. It is the full-fidelity counterpart to CreateAttribute,
// which is kept as a thin wrapper for callers that don't need to set
// any of the "extra" metadata.
type AttributeSpec struct {
	ClassOID      int64
	Name          string
	TypeOID       int64
	Num           int // 1-based ordinal position
	NotNull       bool
	HasDefault    bool
	DeclType      string // original declared type, before normalization
	TypeLength    int    // varchar(N), numeric(p,s).p, etc.
	TypeScale     int    // numeric(p,s).s
	AutoIncrement bool
	IsPrimaryKey  bool
	IsUnique      bool
}

// CreateAttributeSpec inserts a column with the full set of attribute
// metadata.
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

// CreateAttribute inserts a column for the given relation. Wraps
// CreateAttributeSpec for callers that only set the basics.
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

// SetAttributePrimaryKey flips the is_primary_key (and not_null) flag
// on every named column of a relation. Used when a table-level PRIMARY
// KEY constraint is processed after the columns have been registered.
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

// SetAttributeUnique flips is_unique on every named column. Multi-column
// UNIQUE constraints don't make individual columns unique, so callers
// should pass length-1 column lists only.
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

// ColumnInfo describes a resolved column.
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

// ResolveColumn looks up a column by table name and column name.
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

// TableColumns returns all columns for a given table name, ordered by position.
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

// ClassColumn is a column of a relation identified by class OID, used by
// the analyzer to build its per-relation scope.
type ClassColumn struct {
	AttOID  int64
	Name    string
	TypeOID int64
	NotNull bool
}

// ClassColumns returns the columns of a relation by class OID, ordered by
// position.
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

// CodegenColumn is the minimal column shape codegen needs to emit a model
// field: the column name, its resolved type name, and nullability.
type CodegenColumn struct {
	Name     string
	TypeName string
	NotNull  bool
}

// ClassCodegenColumns returns the columns of a relation by class OID in the
// shape the codegen catalog projection consumes.
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

// AttributeDetails describes a column resolved by its OID. Populated
// from sql_attribute joined back to sql_class / sql_namespace, so the
// caller doesn't have to re-resolve names from raw OIDs.
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

// LookupAttribute returns the full source-side metadata for a column,
// keyed by attribute OID. Used by analyzers to populate Column.Source
// (and the per-column flag fields) once they've resolved an
// expression to a sql_attribute row.
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
